package main

import (
	"os"
	"path/filepath"
	"strings"

	"forge/internal/store"
)

// HintsResult is the staged-hints payload for one exercise. Hints are
// derived from the exercise's own metadata (readings, Q&A, grader methods)
// — never from solutions/ — so they guide without spoiling.
type HintsResult struct {
	Hints []string `json:"hints"`
	Fails int      `json:"fails"`
}

// Hints returns three staged study hints. Hint i (0-based) is considered
// unlocked once the learner has accrued i failed runs on the exercise.
func (a *App) Hints(exID string) (*HintsResult, error) {
	ex, err := a.findExercise(exID)
	if err != nil {
		return nil, err
	}
	hist, _ := a.eng.Store.RunHistory(ex.ID, 1000)
	fails := 0
	for _, r := range hist {
		if r.Status == store.StatusFail {
			fails++
		}
	}

	out := &HintsResult{Fails: fails}
	hints := make([]string, 0, 3)

	// Tier 1 — ground yourself in the readings.
	if len(ex.Readings) > 0 {
		var b strings.Builder
		for i, r := range ex.Readings {
			if i > 0 {
				b.WriteString("\n")
			}
			b.WriteString(r.Source)
			if r.Author != "" {
				b.WriteString(" — " + r.Author)
			}
			if r.Chapter != "" {
				b.WriteString(", " + r.Chapter)
			}
			if r.Section != "" {
				b.WriteString(" §" + r.Section)
			}
			if r.URL != "" {
				b.WriteString("\n  " + r.URL)
			}
		}
		hints = append(hints, "Ground yourself in the readings:\n"+b.String())
	}

	// Tier 2 — the comprehension checkpoint this exercise opens with.
	if len(ex.Answers) > 0 {
		qa := ex.Answers[0]
		hints = append(hints,
			"Before writing code, be able to answer this checkpoint from the spec:\n"+qa.Q)
	}

	// Tier 3 — what the grader will check.
	if len(ex.Methods) > 0 {
		var b strings.Builder
		for _, m := range ex.Methods {
			b.WriteString("• " + m.DisplayName() + "\n")
		}
		hints = append(hints,
			"The grader will check these behaviors:\n"+strings.TrimRight(b.String(), "\n"))
	}

	out.Hints = hints
	return out, nil
}

// ReadNotes returns the exercise's personal notes ("" when none exist yet).
func (a *App) ReadNotes(exID string) (string, error) {
	ex, err := a.findExercise(exID)
	if err != nil {
		return "", err
	}
	work, err := a.eng.WorkDir(ex)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(work, "notes.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

// WriteNotes saves the exercise's personal notes into its answers workspace
// (never into subjects/).
func (a *App) WriteNotes(exID, content string) error {
	ex, err := a.findExercise(exID)
	if err != nil {
		return err
	}
	work, err := a.eng.WorkDir(ex)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(work, "notes.md"), []byte(content), 0o644)
}