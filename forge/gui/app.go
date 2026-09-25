package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"forge/internal/check"
	"forge/internal/cur"
	"forge/internal/store"
)

// App is the Wails-bound backend. Every exported method is callable from the
// frontend via the generated bindings (window.go.main.App.*).
type App struct {
	ctx  context.Context
	root string
	eng  *check.Engine
	cur  *cur.Curriculum
}

// NewApp resolves the repository root (the directory containing the
// `subjects` tree, located by walking up from the working directory so the
// app works whether launched from the repo root or from forge/gui) and loads
// the curriculum + progress store.
func NewApp() (*App, error) {
	root, err := findRepoRoot()
	if err != nil {
		return nil, err
	}
	eng, err := check.New(root)
	if err != nil {
		return nil, err
	}
	curriculum, err := cur.Load(eng.Subjects)
	if err != nil {
		eng.Close()
		return nil, err
	}
	sortModulesNumeric(curriculum.Modules)
	return &App{root: root, eng: eng, cur: curriculum}, nil
}

// findRepoRoot walks up from the current directory looking for the repo root
// (a directory containing the subjects/ tree), mirroring `forge init`'s
// notion of the project root.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		st, serr := os.Stat(filepath.Join(dir, "subjects"))
		if serr == nil && st.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("subjects tree not found above %s (run the app from the repository root)", dir)
		}
		dir = parent
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	killAllTerms()
	if a.eng != nil {
		_ = a.eng.Close()
	}
}

// Log forwards a frontend log line (usually an error) to the backend's
// stderr so runtime issues are visible when debugging the webview.
func (a *App) Log(msg string) {
	fmt.Fprintf(os.Stderr, "[frontend] %s\n", msg)
}

// ---- curriculum -------------------------------------------------------

// Curriculum returns the full module tree with per-exercise status.
func (a *App) Curriculum() *Curriculum {
	fmt.Fprintf(os.Stderr, "[backend] curriculum requested\n")
	out := &Curriculum{}
	for _, mod := range a.cur.Modules {
		mi := &ModuleInfo{
			ID:    mod.ID,
			Title: mod.Title,
			Skill: skillFor(mod.ID),
		}
		for i, ex := range mod.Exercises {
			last, _ := a.eng.Store.LastRun(ex.ID)
			eid := mod.ID + "/" + ex.ID
			ei := &ExerciseInfo{
				ID:     ex.ID,
				Key:    eid,
				Title:  ex.Title,
				Lang:   ex.Lang,
				Status: "untried",
				IsGate: i == len(mod.Exercises)-1,
			}
			if last != nil {
				ei.Attempts = 1
				ei.LastAt = last.RanAt.UTC().Format("2006-01-02")
				if last.Status == store.StatusPass {
					ei.Status = "pass"
				} else if last.Status == store.StatusFail {
					ei.Status = "fail"
				} else {
					ei.Status = "error"
				}
			}
			if ei.Status == "pass" {
				mi.Done++
			}
			mi.Total++
			mi.Exercises = append(mi.Exercises, ei)
		}
		if mi.Total > 0 {
			// the module gate is its final exercise (the mini-capstone)
			if ex := mod.Exercises[len(mod.Exercises)-1]; ex != nil {
				mi.Gate = ex.ID
			}
		}
		out.Total += mi.Total
		out.Modules = append(out.Modules, mi)
	}
	return out
}

// GetExercise returns the full detail for one exercise. The id may be a
// plain slug ("ex03-...") or a globally-unique key ("module/exID").
func (a *App) GetExercise(exID string) (*ExerciseDetail, error) {
	mod, ex, err := a.resolveExercise(exID)
	if err != nil {
		return nil, err
	}
	d := &ExerciseDetail{
		ID:       ex.ID,
		Title:    ex.Title,
		Module:   mod.ID,
		Lang:     ex.Lang,
		Readings: []*Reading{},
		QA:       []*QA{},
		Methods:  methodSummaries(ex),
		Files:    []*FileInfo{},
		History:  []*RunInfo{},
	}
	if data, err := os.ReadFile(filepath.Join(ex.Dir, "subject.md")); err == nil {
		d.Spec = string(data)
	} else {
		// some exercises put the spec in a differently-named markdown file
		if entries, derr := os.ReadDir(ex.Dir); derr == nil {
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
					if data, rerr := os.ReadFile(filepath.Join(ex.Dir, e.Name())); rerr == nil {
						d.Spec = string(data)
					}
					break
				}
			}
		}
	}
	for _, r := range ex.Readings {
		d.Readings = append(d.Readings, &Reading{
			Source: r.Source, Author: r.Author, Chapter: r.Chapter,
			Section: r.Section, URL: r.URL,
		})
	}
	for _, qa := range ex.Answers {
		d.QA = append(d.QA, &QA{Q: qa.Q, A: qa.A})
	}
	_ = methodSummaries(ex) // result already set in the struct literal
	// files in the mirrored answers workspace
	if work, werr := a.eng.WorkDir(ex); werr == nil {
		_ = filepath.WalkDir(work, func(p string, de os.DirEntry, err error) error {
			if err != nil || de.IsDir() {
				return nil
			}
			rel, rerr := filepath.Rel(work, p)
			if rerr != nil {
				return nil
			}
			fi, ferr := de.Info()
			if ferr != nil {
				return nil
			}
			d.Files = append(d.Files, &FileInfo{
				Path: filepath.ToSlash(rel),
				Name: de.Name(),
				Size: fi.Size(),
			})
			return nil
		})
	}
	// status + history
	last, _ := a.eng.Store.LastRun(ex.ID)
	if last != nil {
		d.Status = last.Status
		d.LastRun = runInfo(last)
	}
	hist, _ := a.eng.Store.RunHistory(ex.ID, 10)
	for _, r := range hist {
		d.History = append(d.History, runInfo(r))
	}
	if d.Status == "" {
		d.Status = "untried"
	}
	return d, nil
}

// ---- helpers ----------------------------------------------------------

func (a *App) findExercise(exID string) (*cur.Exercise, error) {
	_, ex, err := a.resolveExercise(exID)
	return ex, err
}

// resolveExercise accepts a plain slug, an exact basename, or the fully
// qualified "module/exID" key and returns the matching exercise.
func (a *App) resolveExercise(key string) (*cur.Module, *cur.Exercise, error) {
	// Guard against URL-encoded keys coming from the frontend hash router
	// (e.g. "M0-tools%2Fex01-makefile").
	if dec, err := url.PathUnescape(key); err == nil && dec != key {
		key = dec
	}
	if i := strings.IndexByte(key, '/'); i > 0 {
		modID, exID := key[:i], key[i+1:]
		for _, mod := range a.cur.Modules {
			if mod.ID == modID {
				for _, ex := range mod.Exercises {
					if ex.ID == exID {
						return mod, ex, nil
					}
				}
				return nil, nil, fmt.Errorf("exercise %q not found in module %s", exID, modID)
			}
		}
		return nil, nil, fmt.Errorf("module %q not found", modID)
	}
	return a.eng.FindExercise(key)
}

func methodSummaries(ex *cur.Exercise) []*MethodInfo {
	out := make([]*MethodInfo, 0, len(ex.Methods))
	for _, m := range ex.Methods {
		mi := &MethodInfo{Type: m.Type, Label: m.DisplayName(), Parts: 1}
		switch {
		case m.Build != nil:
			mi.Parts = len(m.Build.Runs)
		case m.Net != nil:
			if len(m.Net.Connections) > 0 {
				mi.Parts = len(m.Net.Connections)
			} else {
				mi.Parts = len(m.Net.Steps)
			}
		case m.Process != nil:
			mi.Parts = len(m.Process.Probes)
		}
		out = append(out, mi)
	}
	return out
}

func runInfo(r *store.Run) *RunInfo {
	out := &RunInfo{
		Status:     r.Status,
		At:         r.RanAt.UTC().Format(time.RFC3339),
		DurationMs: r.DurationMs,
		Output:     r.Stdout,
	}
	if out.Output == "" {
		out.Output = r.Stderr
	}
	return out
}

func sortModulesNumeric(mods []*cur.Module) {
	sort.SliceStable(mods, func(i, j int) bool {
		return moduleNum(mods[i].ID) < moduleNum(mods[j].ID)
	})
}

func moduleNum(id string) int {
	id = strings.TrimPrefix(id, "M")
	if i := strings.IndexByte(id, '-'); i >= 0 {
		id = id[:i]
	}
	n, err := strconv.Atoi(id)
	if err != nil {
		return 0
	}
	return n
}

func safeJoin(base, rel string) (string, error) {
	rel = filepath.FromSlash(rel)
	if rel == "" || filepath.IsAbs(rel) {
		return "", errors.New("invalid path")
	}
	p := filepath.Clean(filepath.Join(base, rel))
	if p == base {
		return "", errors.New("invalid path")
	}
	r, err := filepath.Rel(base, p)
	if err != nil {
		return "", err
	}
	if r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes the exercise workspace: %s", rel)
	}
	return p, nil
}

// moduleSkills mirrors the skill labels used by the TUI profile view.
var moduleSkills = map[string]string{
	"M0-tools":          "Tooling & Shell Scripting",
	"M1-clib":           "C Foundations",
	"M2-procs":          "Processes & Signals",
	"M3-memory":         "Memory & Allocation",
	"M4-concurrency":    "Concurrency",
	"M5-filesio":        "Files & I/O",
	"M6-networking":     "Networking & Sockets",
	"M7-resilience":     "Resilience & Supervision",
	"M8-messaging":      "Messaging & Fault Injection",
	"M9-time":           "Time & Ordering",
	"M10-commit":        "Commit & Consensus",
	"M11-log":           "Replicated Logs",
	"M12-raft":          "Raft Consensus",
	"M13-shard":         "Sharding",
	"M14-membership":    "Membership & Failure Detection",
	"M15-tx":            "Transactions",
	"M16-storage":       "Storage Engines",
	"M17-observability": "Observability",
}

func skillFor(modID string) string {
	if s, ok := moduleSkills[modID]; ok {
		return s
	}
	return modID
}