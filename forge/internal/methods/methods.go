// Package methods implements every grader method type the curriculum uses.
//
// Each runner receives a Ctx describing the exercise and the sandboxed
// working directory the student solution lives in (under answers/), and
// returns a fine-grained Result the orchestrator turns into a run record.
package methods

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"forge/internal/cur"
	"forge/internal/sandbox"
)

// Ctx is the shared context for one exercise check.
type Ctx struct {
	Ex *cur.Exercise
	// Dir is the working directory (answers/<module>/<ex>).
	Dir string
	// Root is the repository root, for resolving bundled tools.
	Root string
	// BaseEnv is injected into every command (e.g. tool paths).
	BaseEnv []string
}

// Opts returns sandbox defaults derived from the exercise config, running
// in the exercise working directory by default.
func (c *Ctx) Opts() sandbox.Options {
	o := sandbox.Options{
		Dir:     c.Dir,
		Timeout: time.Duration(c.Ex.TimeoutS) * time.Second,
		Env:     c.BaseEnv,
	}
	if u := c.Ex.Ulimits; u != nil {
		o.AS, o.FSize, o.NProc, o.Nofile = u.AS, u.FSize, u.NProc, u.Nofile
	}
	return o
}

// Abs resolves a path relative to the exercise working directory.
func (c *Ctx) Abs(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(c.Dir, filepath.Clean(p))
}

// ReadExerciseFile reads a file from the exercise dir.
func (c *Ctx) ReadExerciseFile(p string) ([]byte, error) {
	return sandbox.ReadFile(c.Abs(p))
}

// Tool returns the absolute path of a bundled tool under tools/.
func (c *Ctx) Tool(name string) string {
	return filepath.Join(c.Root, "tools", name)
}

// Part is one indivisible check within a method.
type Part struct {
	Name   string
	Pass   bool
	Detail string
}

// Result is the outcome of one method.
type Result struct {
	Pass   bool
	Parts  []*Part
	Detail string
	Stdout string
	Stderr string
	DurMs  int64
}

// Finish marks the overall method result from its parts.
func (r *Result) Finish() {
	r.Pass = true
	for _, p := range r.Parts {
		if !p.Pass {
			r.Pass = false
			return
		}
	}
}

// DetailIfEmpty fills the console detail from a failing part.
func (r *Result) DetailIfEmpty() {
	if r.Detail != "" {
		return
	}
	for _, p := range r.Parts {
		if !p.Pass && p.Detail != "" {
			r.Detail = p.Detail
			return
		}
	}
}

// Runner executes one grading method.
type Runner func(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error)

var runners = map[string]Runner{
	"build":    runBuild,
	"stdout":   runStdout,
	"artifact": runArtifact,
	"process":  runProcess,
	"net":      runNet,
	"quiz":     runQuiz,
	"report":   runReport,
	"scenario": runScenario,
	"fault":    runScenario,
	"lincheck": runLincheck,
}

// Dispatch runs any declared method type by name.
func Dispatch(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	rn, ok := runners[m.Type]
	if !ok {
		return nil, fmt.Errorf("no runner for method type %q", m.Type)
	}
	res, err := rn(ctx, c, m)
	if err != nil {
		return nil, err
	}
	res.DetailIfEmpty()
	res.Finish()
	return res, nil
}

// ResultText renders a single-line summary for TUI/CLI display.
func ResultText(res *Result) string {
	if res.Pass {
		return "PASS"
	}
	if res.Detail != "" {
		return "FAIL :: " + strings.ReplaceAll(res.Detail, "\n", " ")
	}
	return "FAIL"
}

// decodeBytes turns a configured raw string into actual bytes, honoring
// Go quoted-string escapes when the value starts with a double quote.
func decodeBytes(s string) []byte {
	if s == "" {
		return nil
	}
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		if out, err := strconv.Unquote(s); err == nil {
			return []byte(out)
		}
	}
	return []byte(s)
}

// ContainsBytes reports whether haystack contains needle (needle non-empty).
func ContainsBytes(hay, needle []byte) bool {
	return len(needle) > 0 && bytes.Contains(hay, needle)
}

// Normalize prepares candidate text for comparison.
func normalize(s, mode string) string {
	switch mode {
	case "strip": // collapse all whitespace
		var b strings.Builder
		for _, f := range strings.Fields(s) {
			b.WriteString(f)
		}
		return b.String()
	case "blank": // trim trailing ws per line, collapse empty runs
		lines := strings.Split(s, "\n")
		var out []string
		prevBlank := false
		for _, ln := range lines {
			ln = strings.TrimRight(ln, " \t\r")
			blank := ln == ""
			if blank && prevBlank {
				continue
			}
			prevBlank = blank
			out = append(out, ln)
		}
		for len(out) > 0 && out[len(out)-1] == "" {
			out = out[:len(out)-1]
		}
		return strings.Join(out, "\n")
	case "", "trailing": // trim trailing ws per line, drop final blanks
		lines := strings.Split(s, "\n")
		for i := range lines {
			lines[i] = strings.TrimRight(lines[i], " \t\r")
		}
		for len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		return strings.Join(lines, "\n")
	default:
		return s
	}
}

// diffSummary explains the first mismatch between expected and actual.
func diffSummary(expected, actual string) string {
	el := strings.Split(expected, "\n")
	al := strings.Split(actual, "\n")
	n := len(el)
	if len(al) < n {
		n = len(al)
	}
	for i := 0; i < n; i++ {
		if el[i] != al[i] {
			return fmt.Sprintf("line %d:\n  expected: %q\n  got:      %q", i+1, el[i], al[i])
		}
	}
	if len(el) != len(al) {
		return fmt.Sprintf("line count: expected %d, got %d", len(el), len(al))
	}
	return "no difference"
}
