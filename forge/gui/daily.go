package main

import (
	"time"
)

// Quest is one goal in the daily plan. Progress is derived from the store,
// never persisted, so the plan re-computes each day.
type Quest struct {
	Key    string `json:"key"` // "solve" | "review" | "focus"
	Label  string `json:"label"`
	Target int    `json:"target"`
	Done   int    `json:"done"`
	Unit   string `json:"unit"` // "exercise" | "cards" | "minutes"
}

// DailyPlan is the dashboard "today" summary: three quests, the number of
// review cards due, and the current active-day streak.
type DailyPlan struct {
	Date     string   `json:"date"` // UTC "2006-01-02"
	Quests   []*Quest `json:"quests"`
	DueToday int      `json:"dueToday"`
	Streak   int      `json:"streak"`
	DoneAll  bool     `json:"doneAll"`
}

// DailyPlan computes the current day's plan from existing progress data.
func (a *App) DailyPlan() *DailyPlan {
	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	plan := &DailyPlan{
		Date:   now.Format("2006-01-02"),
		Quests: make([]*Quest, 0, 3),
	}

	// 1. Solve one exercise (a grading run recorded as pass today).
	passed, _ := a.eng.Store.PassesOnDay(dayStart)
	plan.Quests = append(plan.Quests, &Quest{
		Key: "solve", Label: "Solve an exercise",
		Target: 1, Done: passed, Unit: "exercise",
	})

	// 2. Review cards: up to 10 today (target shrinks if the deck is small).
	deck := a.QuizDeck("")
	reviewed, _ := a.eng.Store.ReviewedOnDay(dayStart)
	target := len(deck)
	if target > 10 {
		target = 10
	}
	plan.Quests = append(plan.Quests, &Quest{
		Key: "review", Label: "Review flashcards",
		Target: target, Done: reviewed, Unit: "cards",
	})

	// 3. Focus 25 minutes (pomodoro/focus sessions ending today).
	minutes, _ := a.eng.Store.MinutesBetween(dayStart, now, "")
	plan.Quests = append(plan.Quests, &Quest{
		Key: "focus", Label: "Focus",
		Target: 25, Done: minutes, Unit: "minutes",
	})

	// How many review cards are waiting right now (new cards are "due" too).
	for _, c := range deck {
		if c.Reviewable {
			plan.DueToday++
		}
	}

	// Active-day streak from grading runs (same source as the heatmap).
	act, _ := a.eng.Store.ActivityByDay(now.AddDate(0, 0, -89), now)
	plan.Streak = streakOf(act, now)

	plan.DoneAll = questsDone(plan.Quests)
	return plan
}

func questsDone(qs []*Quest) bool {
	for _, q := range qs {
		if q.Target > 0 && q.Done < q.Target {
			return false
		}
	}
	return true
}