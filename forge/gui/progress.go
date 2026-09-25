package main

import (
	"time"

	"forge/internal/cur"
)

// Progress aggregates the dashboard numbers: overall pass counts, streak,
// activity heatmap (runs/day), focused minutes, and per-module bars.
func (a *App) Progress() *Progress {
	now := time.Now().UTC()

	var all []string
	for _, mod := range a.cur.Modules {
		all = append(all, exerciseIDs(mod)...)
	}
	passed, attempted, total, _ := a.eng.Store.OverallProgress(all)

	// runs per day, last 90 days
	start := now.AddDate(0, 0, -89)
	act, _ := a.eng.Store.ActivityByDay(start, now)
	var days []*DayActivity
	todayKey := now.Format("2006-01-02")
	todayCount := 0
	for d := start; !d.After(now); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		n := act[key]
		if key == todayKey {
			todayCount = n
		}
		days = append(days, &DayActivity{Day: key, Count: n})
	}

	// focused minutes, last 7 days + current week total
	weekStart := now.AddDate(0, 0, -6)
	minAct, _ := a.eng.Store.MinutesByDay(weekStart, now)
	weekMin, _ := a.eng.Store.MinutesBetween(weekStart, now, "")
	var minDays []*DayActivity
	for d := weekStart; !d.After(now); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		minDays = append(minDays, &DayActivity{Day: key, Count: minAct[key]})
	}

	var bars []*ModuleBar
	for _, mod := range a.cur.Modules {
		done, cnt, _ := a.eng.Store.ModuleProgress(exerciseIDs(mod))
		bars = append(bars, &ModuleBar{ID: mod.ID, Done: done, Total: cnt})
	}

	return &Progress{
		Passed:    passed,
		Attempted: attempted,
		Total:     total,
		Streak:    streakOf(act, now),
		TodayRuns: todayCount,
		Activity:  days,
		Minutes:   minDays,
		WeekMin:   weekMin,
		PerModule: bars,
	}
}

// streakOf computes the current active-day streak ending today or yesterday
// (so an untouched "today" doesn't kill the streak before the day is over).
func streakOf(act map[string]int, now time.Time) int {
	d := now
	if act[d.Format("2006-01-02")] == 0 {
		d = d.AddDate(0, 0, -1)
	}
	s := 0
	for act[d.Format("2006-01-02")] > 0 {
		s++
		d = d.AddDate(0, 0, -1)
	}
	return s
}

func exerciseIDs(mod *cur.Module) []string {
	out := make([]string, len(mod.Exercises))
	for i, ex := range mod.Exercises {
		out[i] = ex.ID
	}
	return out
}