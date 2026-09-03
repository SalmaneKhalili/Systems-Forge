// Package cur models the curriculum: modules, exercises, and the JSON
// grading configuration that ships inside each exercise directory.
//
// The canonical tree lives under subjects/. Every exercise directory is
// exNN-<slug> and must contain:
//
//	exercise.json    the machine-readable grading config (this package)
//	subject.md       the human-readable specification (goal, constraints,
//	                 acceptance criteria, readings pointers)
//
// A module directory is <NN>-<slug> and may contain module.md.
package cur

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// MethodName is the allowed set of grader method types.
var MethodNames = map[string]bool{
	"build":    true,
	"stdout":   true,
	"artifact": true,
	"process":  true,
	"net":      true,
	"quiz":     true,
	"report":   true,
	"scenario": true,
	"fault":    true,
	"lincheck": true,
}

// Exercise is the top-level grading config found in exercise.json.
//
// Dir holds the on-disk exercise directory (populated by the loader, not
// part of the JSON config).
//
// expect is used by the selftest harness (tools/fixtures) to declare the
// intended outcome of a fixture exercise.
type Exercise struct {
	ID         string     `json:"id"`
	Module     string     `json:"module"`
	Title      string     `json:"title"`
	Lang       string     `json:"lang"`
	TimeoutS   int        `json:"timeout_s,omitempty"`
	Ulimits    *Ulimits   `json:"ulimits,omitempty"`
	Sanitize   bool       `json:"sanitize,omitempty"`
	Sanitizers []string   `json:"sanitizers,omitempty"`
	Methods    []*Method  `json:"methods"`
	Readings   []*Reading `json:"readings,omitempty"`
	Answers    []*QA      `json:"answers,omitempty"`
	Expect     string     `json:"expect,omitempty"` // selftest marker: pass|fail

	Dir string `json:"-"`
}

// Ulimits applied to every command run for this exercise.
type Ulimits struct {
	AS     int64 `json:"as,omitempty"`     // address space, bytes
	FSize  int64 `json:"fsize,omitempty"`  // file size, bytes
	NProc  int64 `json:"nproc,omitempty"`  // processes/threads
	Nofile int64 `json:"nofile,omitempty"` // open files
}

// Method configures one grading method for an exercise.
type Method struct {
	Type string `json:"type"`

	// exactly one of the following is expected per method
	Build    *BuildSpec    `json:"build,omitempty"`
	Stdout   *StdoutSpec   `json:"stdout,omitempty"`
	Artifact *ArtifactSpec `json:"artifact,omitempty"`
	Process  *ProcessSpec  `json:"process,omitempty"`
	Net      *NetSpec      `json:"net,omitempty"`
	Quiz     *QuizSpec     `json:"quiz,omitempty"`
	Report   *ReportSpec   `json:"report,omitempty"`
	Scenario *ScenarioSpec `json:"scenario,omitempty"`
	Fault    *FaultSpec    `json:"fault,omitempty"`
	Lincheck *LincheckSpec `json:"lincheck,omitempty"`

	Label string `json:"label,omitempty"` // display name
}

// DisplayName returns the method label or its type when unlabeled.
func (m *Method) DisplayName() string {
	if m.Label != "" {
		return m.Label
	}
	return m.Type
}

// BuildSpec compiles the exercise (typically with a student-authored
// Makefile) under strict flags and runs verification binaries.
type BuildSpec struct {
	Make  string    `json:"make,omitempty"`  // target, default "all"
	Clean string    `json:"clean,omitempty"` // target to run before build, default "fclean"
	Sub   string    `json:"sub,omitempty"`   // subdirectory to build in, relative to exercise dir
	Runs  []*RunCmd `json:"runs,omitempty"`
}

// RunCmd is a single program invocation with expected output.
type RunCmd struct {
	Name        string   `json:"name"`
	Cmd         []string `json:"cmd"`
	Input       string   `json:"input,omitempty"`
	Expect      string   `json:"expect,omitempty"` // relative path to expected-output file
	Normalize   string   `json:"normalize,omitempty"`
	MustSucceed *bool    `json:"must_succeed,omitempty"` // default true
	ExpectExit  int      `json:"expect_exit,omitempty"`  // default 0
	StderrOK    bool     `json:"stderr_ok,omitempty"`    // default: stderr must be empty
}

// StdoutSpec runs a program once and diffs its stdout.
type StdoutSpec struct {
	Cmd       []string `json:"cmd"`
	Input     string   `json:"input,omitempty"`
	Expect    string   `json:"expect"` // relative path to expected-output file
	Normalize string   `json:"normalize,omitempty"`
	StderrOK  bool     `json:"stderr_ok,omitempty"`
}

// ArtifactSpec asserts properties of produced files.
type ArtifactSpec struct {
	Path     string   `json:"path"`
	Wipe     bool     `json:"wipe,omitempty"`     // remove before grading
	Exists   *bool    `json:"exists,omitempty"`   // default true
	IsA      string   `json:"is_a,omitempty"`     // prefix match of `file` output
	MinSize  int64    `json:"min_size,omitempty"` // bytes
	MaxSize  int64    `json:"max_size,omitempty"` // bytes
	Symbols  []string `json:"symbols,omitempty"`  // required `nm` symbols
	Contains []string `json:"contains,omitempty"` // required substrings in text payload
}

// ProcessSpec spawns a program, probes it while running, then stops it.
type ProcessSpec struct {
	Start     []string          `json:"start"`
	WaitMs    int               `json:"wait_ms,omitempty"` // default 1000
	TargetEnv map[string]string `json:"target_env,omitempty"`
	Probes    []*Probe          `json:"probes,omitempty"`
	Kill      string            `json:"kill,omitempty"` // "term" (default) | "kill"
}

// Probe runs a command while the target process is alive.
type Probe struct {
	Exec       []string          `json:"exec"`
	ExpectExit int               `json:"expect_exit,omitempty"` // default 0
	ExpectOut  string            `json:"expect_out,omitempty"`  // substring of probe stdout
	Env        map[string]string `json:"env,omitempty"`
}

// NetSpec drives a scripted protocol against a spawned server.
type NetSpec struct {
	Start       []string          `json:"start"`
	WaitMs      int               `json:"wait_ms,omitempty"`
	Host        string            `json:"host,omitempty"` // default 127.0.0.1
	Port        int               `json:"port"`
	TimeoutMs   int               `json:"timeout_ms,omitempty"`
	StartEnv    map[string]string `json:"start_env,omitempty"`
	Steps       []*NetStep        `json:"steps"`
	Connections []*NetSession     `json:"connections,omitempty"`
}

// NetSession is one logical connection in a multi-connection net run. Each
// session dials a fresh TCP connection against the same server process.
type NetSession struct {
	Steps []*NetStep `json:"steps"`
}

// NetStep is one send/receive interaction.
type NetStep struct {
	Send      string `json:"send,omitempty"`     // raw bytes, Go-style escapes decoded
	SendHex   string `json:"send_hex,omitempty"` // hex encode, overrides Send
	Expect    string `json:"expect,omitempty"`   // full-match bytes, Go-style escapes decoded
	ExpectHex string `json:"expect_hex,omitempty"`
	Contains  string `json:"contains,omitempty"` // substring match on received bytes
	SleepMs   int    `json:"sleep_ms,omitempty"`

	// Signal sends an OS signal to the server process group.
	// Accepted names: TERM, KILL, USR1, USR2, HUP, INT (SIG prefix optional).
	Signal string `json:"signal,omitempty"`

	// WaitExit waits for the server process to exit and checks the exit code.
	// Nil means "just wait, don't check the code". Non-nil means check the code matches.
	WaitExit *int `json:"wait_exit,omitempty"`

	// Restart performs a supervisor restart cycle: kill the running server
	// (using Signal, default TERM), wait for it to exit (checking WaitExit if
	// set), respawn the same start command, and re-dial. Subsequent steps in
	// this connection run against the restarted server. This grades that a
	// service comes back up and serves again after being killed (a bounded
	// restart policy, cf. M7-ex05 watchdog / Kubernetes Restart Policies).
	Restart bool `json:"restart,omitempty"`
}

// QuizSpec grades a checkpoint answers file against the exercise's Answers.
type QuizSpec struct {
	File string `json:"file,omitempty"` // default "quiz.txt"
}

// ReportSpec compares structured output against a reference run.
type ReportSpec struct {
	Cmd       []string           `json:"cmd"`
	Reference string             `json:"reference"`          // relative path to reference k=v file
	Tolerate  map[string]float64 `json:"tolerate,omitempty"` // key -> relative tolerance
}

// ScenarioSpec / FaultSpec run a driver that boots nodes and injects faults.
// The driver exits 0 on success; a line "FORGE_RESULT: pass|fail" overrides
// its exit code.
type ScenarioSpec struct {
	Driver   []string          `json:"driver"`
	Env      map[string]string `json:"env,omitempty"`
	TimeoutS int               `json:"timeout_s,omitempty"` // default from exercise
}

// FaultSpec is structurally identical to ScenarioSpec. Kept distinct so a
// fault grader can diverge later (e.g., assert properties, not outputs).
type FaultSpec struct {
	Driver   []string          `json:"driver"`
	Env      map[string]string `json:"env,omitempty"`
	TimeoutS int               `json:"timeout_s,omitempty"`
}

// LincheckSpec checks recorded operation histories for linearizability.
type LincheckSpec struct {
	History string `json:"history"`       // glob/path to history.json produced by a driver
	Err     string `json:"err,omitempty"` // substring expected in checker stderr (optional)
}

// Reading points the learner at a specific chapter/section/article.
type Reading struct {
	Source  string `json:"source"`
	Author  string `json:"author,omitempty"`
	Chapter string `json:"chapter,omitempty"`
	Section string `json:"section,omitempty"`
	URL     string `json:"url,omitempty"`
}

// QA is one reading-comprehension checkpoint item.
type QA struct {
	Q string `json:"q"`
	A string `json:"a"`
}

// Free-standing concrete types so callers don't need the Exercise wrapper.

// Module is a curriculum module and its exercises.
type Module struct {
	ID        string // e.g. "M0-tools"
	Dir       string // absolute path of the module directory
	Title     string
	Exercises []*Exercise
}

// Exercise returns the exercise with the given ID inside this module.
func (m *Module) Exercise(id string) *Exercise {
	for _, ex := range m.Exercises {
		if ex.ID == id {
			return ex
		}
	}
	return nil
}

// Curriculum is the full loaded subject tree.
type Curriculum struct {
	Root    string // canonical subjects root
	Modules []*Module
}

// Load parses the curriculum tree rooted at subjectsDir.
func Load(subjectsDir string) (*Curriculum, error) {
	abs, err := filepath.Abs(subjectsDir)
	if err != nil {
		return nil, err
	}
	cur := &Curriculum{Root: abs}
	entries, err := os.ReadDir(abs)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("subjects tree not found at %s (run `forge init`?)", abs)
		}
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		mod, err := loadModule(filepath.Join(abs, e.Name()))
		if err != nil {
			return nil, err
		}
		cur.Modules = append(cur.Modules, mod)
	}
	sort.Slice(cur.Modules, func(i, j int) bool { return cur.Modules[i].Dir < cur.Modules[j].Dir })
	return cur, nil
}

var exDirRe = regexp.MustCompile(`^ex\d{2}-.+$`)

func loadModule(dir string) (*Module, error) {
	mod := &Module{ID: filepath.Base(dir), Dir: dir}
	title, err := moduleTitle(dir)
	if err == nil {
		mod.Title = title
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() || !exDirRe.MatchString(e.Name()) {
			continue
		}
		exDir := filepath.Join(dir, e.Name())
		ex, err := LoadExercise(exDir)
		if err != nil {
			return nil, err
		}
		mod.Exercises = append(mod.Exercises, ex)
	}
	sort.Slice(mod.Exercises, func(i, j int) bool { return mod.Exercises[i].ID < mod.Exercises[j].ID })
	return mod, nil
}

const defaultTimeoutS = 20
const defaultWaitMs = 1000

// LoadExercise parses exercise.json from dir and applies defaults.
func LoadExercise(dir string) (*Exercise, error) {
	cfgPath := filepath.Join(dir, "exercise.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("exercise %s: %w", dir, err)
	}
	var ex Exercise
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&ex); err != nil {
		return nil, fmt.Errorf("exercise %s: malformed exercise.json: %w", dir, err)
	}

	// Validate identity from the directory name.
	base := filepath.Base(dir)
	if !strings.HasPrefix(base, "ex") {
		return nil, fmt.Errorf("exercise dir must be named exNN-<slug>, got %q", base)
	}
	if ex.ID == "" {
		ex.ID = base
	}
	modName := filepath.Base(filepath.Dir(dir))
	if ex.Module == "" {
		ex.Module = modName
	}
	ex.Dir = dir
	if ex.TimeoutS <= 0 {
		ex.TimeoutS = defaultTimeoutS
	}
	if err := ex.normalize(); err != nil {
		return nil, fmt.Errorf("exercise %s: %w", dir, err)
	}
	return &ex, nil
}

func (ex *Exercise) normalize() error {
	if len(ex.Methods) == 0 {
		return errors.New("exercise has no grading methods")
	}
	seen := map[string]bool{}
	labels := map[string]bool{}
	for i, m := range ex.Methods {
		if !MethodNames[m.Type] {
			return fmt.Errorf("method %d: unknown method type %q", i, m.Type)
		}
		if seen[m.Type] {
			return fmt.Errorf("duplicate method type %q", m.Type)
		}
		seen[m.Type] = true
		if _, ok := labels[m.DisplayName()]; ok {
			return fmt.Errorf("duplicate method label %q", m.DisplayName())
		}
		labels[m.DisplayName()] = true
		switch m.Type {
		case "build":
			if m.Build == nil {
				return errors.New("build method requires build spec")
			}
		case "stdout":
			if m.Stdout == nil {
				return errors.New("stdout method requires stdout spec")
			}
			if m.Stdout.Expect == "" {
				return errors.New("stdout method requires expect file")
			}
		case "artifact":
			if m.Artifact == nil {
				return errors.New("artifact method requires artifact spec")
			}
			if m.Artifact.Path == "" {
				return errors.New("artifact method requires path")
			}
		case "process":
			if m.Process == nil {
				return errors.New("process method requires process spec")
			}
			if len(m.Process.Start) == 0 {
				return errors.New("process method requires start command")
			}
			if m.Process.WaitMs <= 0 {
				m.Process.WaitMs = defaultWaitMs
			}
		case "net":
			if m.Net == nil {
				return errors.New("net method requires net spec")
			}
			if len(m.Net.Start) == 0 {
				return errors.New("net method requires start command")
			}
			if m.Net.Port <= 0 {
				return errors.New("net method requires port")
			}
			if m.Net.WaitMs <= 0 {
				m.Net.WaitMs = defaultWaitMs
			}
		case "quiz":
			if m.Quiz == nil {
				return errors.New("quiz method requires quiz spec")
			}
			if m.Quiz.File == "" {
				m.Quiz.File = "quiz.txt"
			}
			if len(ex.Answers) == 0 {
				return errors.New("quiz method requires exercise answers")
			}
		case "report":
			if m.Report == nil {
				return errors.New("report method requires report spec")
			}
			if len(m.Report.Cmd) == 0 || m.Report.Reference == "" {
				return errors.New("report method requires cmd and reference")
			}
		case "scenario":
			if m.Scenario == nil {
				return errors.New("scenario method requires scenario spec")
			}
			if len(m.Scenario.Driver) == 0 {
				return errors.New("scenario method requires driver")
			}
		case "fault":
			if m.Fault == nil {
				return errors.New("fault method requires fault spec")
			}
			if len(m.Fault.Driver) == 0 {
				return errors.New("fault method requires driver")
			}
		case "lincheck":
			if m.Lincheck == nil {
				return errors.New("lincheck method requires lincheck spec")
			}
			if m.Lincheck.History == "" {
				return errors.New("lincheck method requires history")
			}
		}
	}
	return nil
}

func moduleTitle(dir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dir, "module.md"))
	if err != nil {
		return filepath.Base(dir), nil
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# ")), nil
		}
	}
	return filepath.Base(dir), nil
}
