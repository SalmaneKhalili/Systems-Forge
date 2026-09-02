// Package check orchestrates a single exercise check: locate the student
// workspace dir, apply every method runner, aggregate the result, and
// record it in the progress store.
package check

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"forge/internal/cur"
	"forge/internal/methods"
	"forge/internal/sandbox"
	"forge/internal/store"
)

// MethodReport is the outcome of one method within a check.
type MethodReport struct {
	Type   string
	Label  string
	Pass   bool
	Parts  []*PartReport `json:",omitempty"`
	Output string        `json:",omitempty"`
	DurMs  int64
}

// PartReport is one indivisible assertion within a method.
type PartReport struct {
	Name   string
	Pass   bool
	Detail string `json:",omitempty"`
}

// Report is the outcome of one exercise check.
type Report struct {
	Exercise string
	Pass     bool
	Methods  []*MethodReport
	Duration time.Duration
}

// Engine performs checks against the answers workspace.
type Engine struct {
	Root     string // repository root
	Subjects string // canonical subjects dir (under Root)
	Answers  string // personal workspace dir (under Root)
	Store    *store.Store
}

// New creates an engine rooted at the repository directory.
func New(root string) (*Engine, error) {
	st, err := store.Open(filepath.Join(root, "progress.db"))
	if err != nil {
		return nil, err
	}
	return &Engine{
		Root:     root,
		Subjects: filepath.Join(root, "subjects"),
		Answers:  filepath.Join(root, "answers"),
		Store:    st,
	}, nil
}

// Close releases engine resources.
func (e *Engine) Close() error { return e.Store.Close() }

// FindExercise locates an exercise by ID across all modules.
func (e *Engine) FindExercise(id string) (*cur.Module, *cur.Exercise, error) {
	curriculum, err := cur.Load(e.Subjects)
	if err != nil {
		return nil, nil, err
	}
	for _, mod := range curriculum.Modules {
		for _, ex := range mod.Exercises {
			if ex.ID == id {
				return mod, ex, nil
			}
		}
	}
	// Substring match on the slug part (ex03-...).
	for _, mod := range curriculum.Modules {
		for _, ex := range mod.Exercises {
			if strings.Contains(ex.ID, "-") && strings.SplitN(ex.ID, "-", 2)[1] == id {
				return mod, ex, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("exercise %q not found in curriculum", id)
}

// WorkDir returns the personal workspace dir for an exercise, creating it
// and mirroring subject scaffolding as needed.
func (e *Engine) WorkDir(ex *cur.Exercise) (string, error) {
	rel, err := filepath.Rel(e.Subjects, ex.Dir)
	if err != nil {
		return "", err
	}
	dst := filepath.Join(e.Answers, filepath.FromSlash(rel))
	if err := ensureMirror(ex.Dir, dst); err != nil {
		return "", err
	}
	return dst, nil
}

// Check runs every method for an exercise and records the outcome.
func (e *Engine) Check(ctx context.Context, ex *cur.Exercise) (*Report, error) {
	work, err := e.WorkDir(ex)
	if err != nil {
		return nil, err
	}

	// Pre-clean generated artifacts.
	for _, m := range ex.Methods {
		if m.Artifact != nil && m.Artifact.Wipe {
			_ = os.Remove(filepath.Join(work, m.Artifact.Path))
		}
	}

	c := &methods.Ctx{Ex: ex, Dir: work, Root: e.Root}
	report := &Report{Exercise: ex.ID}
	start := time.Now()
	for _, m := range ex.Methods {
		res, err := methods.Dispatch(ctx, c, m)
		dur := time.Since(start)
		mr := &MethodReport{Type: m.Type, Label: m.DisplayName(), Pass: false, DurMs: dur.Milliseconds()}
		if err != nil {
			mr.Output = "checker error: " + err.Error()
		} else {
			mr.Pass = res.Pass
			for _, p := range res.Parts {
				mr.Parts = append(mr.Parts, &PartReport{Name: p.Name, Pass: p.Pass, Detail: p.Detail})
			}
			mr.Output = strings.TrimSpace(res.Stderr)
			if mr.Output == "" {
				mr.Output = strings.TrimSpace(res.Stdout)
			}
		}
		report.Methods = append(report.Methods, mr)
		if !mr.Pass {
			break
		}
	}
	report.Duration = time.Since(start)
	report.Pass = true
	for _, m := range report.Methods {
		if !m.Pass {
			report.Pass = false
			break
		}
	}

	// Persist the run.
	r := &store.Run{
		Exercise:   ex.ID,
		Status:     store.StatusFail,
		Detail:     detailJSON(report),
		DurationMs: report.Duration.Milliseconds(),
		RanAt:      time.Now(),
	}
	if report.Pass {
		r.Status = store.StatusPass
	}
	for _, m := range report.Methods {
		if !m.Pass {
			r.Stdout = m.Output
			break
		}
	}
	_ = e.Store.RecordRun(r)
	return report, nil
}

func detailJSON(r *Report) string {
	b, err := json.Marshal(r.Methods)
	if err != nil {
		return ""
	}
	return string(b)
}

// ensureMirror copies subject scaffolding into the answers workspace so the
// exercise is runnable there unchanged; student files are never overwritten.
func ensureMirror(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, ent := range entries {
		sp := filepath.Join(src, ent.Name())
		dp := filepath.Join(dst, ent.Name())
		info, err := ent.Info()
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := ensureMirror(sp, dp); err != nil {
				return err
			}
			continue
		}
		spInfo, err := os.Stat(sp)
		if err != nil {
			return err
		}
		if dpInfo, err := os.Stat(dp); err == nil {
			// Never clobber student-written files.
			if !spInfo.ModTime().After(dpInfo.ModTime()) {
				continue
			}
		}
		data, err := sandbox.ReadFile(sp)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dp, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
