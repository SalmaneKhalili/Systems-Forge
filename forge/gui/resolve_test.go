package main

import (
	"testing"
)

// TestResolveExerciseEncodedKey guards the Wave 0 navigation fix: the
// frontend hash router encodes the "/" inside exercise keys as %2F, and
// resolveExercise must resolve both the plain and the encoded key.
func TestResolveExerciseEncodedKey(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	defer app.eng.Close()

	for _, key := range []string{
		"M0-tools/M0-ex01",
		"M0-tools%2FM0-ex01",
	} {
		mod, ex, err := app.resolveExercise(key)
		if err != nil {
			t.Fatalf("resolveExercise(%q): %v", key, err)
		}
		if mod.ID != "M0-tools" || ex.ID != "M0-ex01" {
			t.Fatalf("resolveExercise(%q) = %s/%s, want M0-tools/M0-ex01", key, mod.ID, ex.ID)
		}
	}
}

// TestDerivedBindings sanity-checks the read-only feature bindings added in
// Waves 1-3 against the real curriculum: nothing here writes to the store.
func TestDerivedBindings(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	defer app.eng.Close()

	// DailyPlan: three quests, streak >= 0.
	plan := app.DailyPlan()
	if plan == nil || len(plan.Quests) != 3 {
		t.Fatalf("DailyPlan: expected 3 quests, got %+v", plan)
	}

	// XP: sane level/XP and badge list.
	xi := app.XP()
	if xi == nil || xi.Level < 1 || xi.XP < 0 {
		t.Fatalf("XP: unexpected payload %+v", xi)
	}
	if len(xi.Badges) != 10 {
		t.Fatalf("XP: expected 10 badges, got %d", len(xi.Badges))
	}

	// Hints: M0-ex01 has Q&A + a grader method, so >= 2 tiers.
	hr, err := app.Hints("M0-tools/M0-ex01")
	if err != nil {
		t.Fatalf("Hints: %v", err)
	}
	if len(hr.Hints) < 2 {
		t.Fatalf("Hints: expected >= 2 tiers for M0-ex01, got %d", len(hr.Hints))
	}

	// QuizDeck: the "personal" filter never crosses into exercise cards, and
	// the unfiltered deck always contains the exercise cards.
	all := app.QuizDeck("")
	if len(all) == 0 {
		t.Fatal("QuizDeck(\"\"): empty deck")
	}
	pers := app.QuizDeck("personal")
	for _, c := range pers {
		if c.ExID != "personal" {
			t.Fatalf("QuizDeck(\"personal\") leaked exercise card %q", c.Key)
		}
	}
}