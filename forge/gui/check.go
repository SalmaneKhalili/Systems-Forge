package main

import (
	"context"
	"strings"

	"forge/internal/check"
)

// RunCheck runs the autograder for one exercise and returns the full report.
// The workdir mirror + SQLite record happen inside check.Engine.Check, so the
// UI stays in sync with the CLI/TUI progress store automatically.
func (a *App) RunCheck(exID string) (*CheckResult, error) {
	ex, err := a.findExercise(exID)
	if err != nil {
		return nil, err
	}
	report, err := a.eng.Check(context.Background(), ex)
	if err != nil {
		return nil, err
	}
	return toCheckResult(report), nil
}

func toCheckResult(r *check.Report) *CheckResult {
	out := &CheckResult{
		ID:         r.Exercise,
		Pass:       r.Pass,
		DurationMs: r.Duration.Milliseconds(),
		Methods:    make([]*MethodResult, 0, len(r.Methods)),
	}
	for _, m := range r.Methods {
		out.Methods = append(out.Methods, toMethodResult(m))
		if !m.Pass {
			out.Output = strings.TrimSpace(m.Output)
		}
	}
	if out.Output == "" && !r.Pass {
		// fall back to the first failing part detail
		for _, m := range r.Methods {
			for _, p := range m.Parts {
				if !p.Pass && p.Detail != "" {
					out.Output = p.Detail
					break
				}
			}
			if out.Output != "" {
				break
			}
		}
	}
	return out
}

func toMethodResult(m *check.MethodReport) *MethodResult {
	mr := &MethodResult{
		Type:   m.Type,
		Label:  m.Label,
		Pass:   m.Pass,
		DurMs:  m.DurMs,
		Output: m.Output,
		Parts:  make([]*PartResult, 0, len(m.Parts)),
	}
	for _, p := range m.Parts {
		mr.Parts = append(mr.Parts, &PartResult{
			Name: p.Name, Pass: p.Pass, Detail: p.Detail,
		})
	}
	return mr
}

// RunCheckMethod re-runs a single grader method (no store record) for a fast
// debug loop on one failing check.
func (a *App) RunCheckMethod(exID string, idx int) (*CheckResult, error) {
	ex, err := a.findExercise(exID)
	if err != nil {
		return nil, err
	}
	mr, err := a.eng.CheckMethod(context.Background(), ex, idx)
	if err != nil {
		return nil, err
	}
	return &CheckResult{
		ID:         ex.ID,
		Pass:       mr.Pass,
		DurationMs: mr.DurMs,
		Methods:    []*MethodResult{toMethodResult(mr)},
	}, nil
}