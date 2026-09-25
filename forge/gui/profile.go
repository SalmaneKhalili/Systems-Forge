package main

import (
	"math"
	"time"

	"forge/internal/store"
)

// BadgeInfo is one achievement shown in the profile view.
type BadgeInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Icon       string `json:"icon"`
	Earned     bool   `json:"earned"`
	UnlockedAt string `json:"unlockedAt,omitempty"`
}

// XPInfo is the XP/level/badge payload for the profile view.
type XPInfo struct {
	XP     int          `json:"xp"`
	Level  int          `json:"level"`
	XPIn   int          `json:"xpIn"`   // xp earned inside the current level
	XPNext int          `json:"xpNext"` // xp needed to complete the current level
	Badges []*BadgeInfo `json:"badges"`
	Passed int          `json:"passed"`
	Total  int          `json:"total"`
}

// XP computes the learner's XP, level, and badge set. Badge unlock times are
// persisted in settings (keys "badges.<id>") on first earning.
func (a *App) XP() *XPInfo {
	all := make([]string, 0)
	for _, mod := range a.cur.Modules {
		all = append(all, exerciseIDs(mod)...)
	}
	passed, attempted, total, _ := a.eng.Store.OverallProgress(all)
	states, _ := a.eng.Store.AllReviews()
	minutes, _ := a.eng.Store.MinutesBetween(time.Time{}, time.Now().UTC(), "")

	xp := passed*100 + reviewReps(states)*5 + minutes*2 + attempted*10
	level := levelFor(xp)
	in := xp - levelBase(level)
	next := levelNext(level)

	return &XPInfo{
		XP: xp, Level: level, XPIn: in, XPNext: next,
		Badges: a.computeBadges(passed, total, minutes, reviewReps(states)),
		Passed: passed, Total: total,
	}
}

func reviewReps(states map[string]*store.ReviewState) int {
	n := 0
	for _, rs := range states {
		n += rs.Reps
	}
	return n
}

// levelFor: level 1 at 0 XP, and each level costs (2·level−1)·100 XP (i.e.
// cumulative 100·level²). square-root ladder keeps early levels fast.
func levelFor(xp int) int {
	if xp < 100 {
		return 1
	}
	return int(math.Sqrt(float64(xp)/100)) + 1
}

// levelBase is the XP required to have completed level l-1.
func levelBase(level int) int {
	return 100 * (level - 1) * (level - 1)
}

// levelNext is the XP needed to finish the current level.
func levelNext(level int) int {
	return 100 * (2*level - 1)
}

type badgeDef struct {
	id, name, desc, icon string
	earned               bool
}

func (a *App) computeBadges(passed, total, minutes, repCount int) []*BadgeInfo {
	now := time.Now().UTC()
	settings, _ := a.eng.Store.AllSettings()
	streak := a.currentStreak()
	gates := a.gatesPassed()

	defs := []badgeDef{
		{"first_pass", "First Steps", "Pass your first exercise.", "🌱", passed >= 1},
		{"pass_10", "Rising Forge", "10 exercises passed.", "🔨", passed >= 10},
		{"pass_25", "Quarter Way", "25 exercises passed.", "🔥", passed >= 25},
		{"pass_50", "Half Forged", "50 exercises passed.", "⚔️", passed >= 50},
		{"pass_all", "Complete", "Every exercise passed.", "🏆", total > 0 && passed >= total},
		{"gate_slayer", "Gate Crasher", "Pass your first module gate.", "🚪", gates >= 1},
		{"reviews_50", "Card Shark", "50 flashcard reviews.", "🃏", repCount >= 50},
		{"streak_7", "On a Roll", "7-day study streak.", "📈", streak >= 7},
		{"streak_30", "Unstoppable", "30-day study streak.", "👑", streak >= 30},
		{"focus_10h", "Deep Focus", "10 hours of focused study.", "🧘", minutes >= 600},
	}

	newly := map[string]string{}
	for _, d := range defs {
		if d.earned {
			if _, ok := settings["badges."+d.id]; !ok {
				newly["badges."+d.id] = now.Format(time.RFC3339)
			}
		}
	}
	if len(newly) > 0 {
		_ = a.eng.Store.SetSettingsBulk(newly)
	}

	out := make([]*BadgeInfo, 0, len(defs))
	for _, d := range defs {
		b := &BadgeInfo{ID: d.id, Name: d.name, Desc: d.desc, Icon: d.icon, Earned: d.earned}
		if at, ok := settings["badges."+d.id]; ok && d.earned {
			b.UnlockedAt = at
		}
		out = append(out, b)
	}
	return out
}

// currentStreak is the active-day grading streak (shared with the dashboard).
func (a *App) currentStreak() int {
	now := time.Now().UTC()
	act, _ := a.eng.Store.ActivityByDay(now.AddDate(0, 0, -89), now)
	return streakOf(act, now)
}

// gatesPassed counts module gates (final exercise of each module) with a
// passing latest run.
func (a *App) gatesPassed() int {
	n := 0
	for _, mod := range a.cur.Modules {
		if len(mod.Exercises) == 0 {
			continue
		}
		gate := mod.Exercises[len(mod.Exercises)-1]
		if last, _ := a.eng.Store.LastRun(gate.ID); last != nil && last.Status == store.StatusPass {
			n++
		}
	}
	return n
}