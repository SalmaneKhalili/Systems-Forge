package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"forge/internal/store"
)

// QuizDeck returns every review card built from the exercises' Answers[]
// metadata, merged with persisted SM-2 state. An empty moduleFilter means all
// modules; otherwise only cards from that module are returned.
func (a *App) QuizDeck(moduleFilter string) []*Card {
	states, _ := a.eng.Store.AllReviews()
	now := time.Now().UTC()
	cards := make([]*Card, 0)
	for _, mod := range a.cur.Modules {
		if moduleFilter != "" && mod.ID != moduleFilter {
			continue
		}
		for _, ex := range mod.Exercises {
			for i, qa := range ex.Answers {
				exid := mod.ID + "/" + ex.ID
				key := fmt.Sprintf("%s:%d", exid, i)
				c := &Card{
					Key: key, Module: mod.ID, ExID: exid, Title: ex.Title,
					Q: qa.Q, A: qa.A, New: true, Reviewable: true,
				}
				if rs := states[key]; rs != nil {
					c.New = false
					c.Ease = rs.Ease
					c.Interval = rs.Interval
					c.Due = rs.Due.UTC().Format("2006-01-02")
					c.Reps = rs.Reps
					c.Lapses = rs.Lapses
					c.Reviewable = !now.Before(rs.Due)
				}
				cards = append(cards, c)
			}
		}
	}
	// Personal cards (SM-2 state merged the same way). Only included when the
	// filter is empty or "personal".
	if moduleFilter == "" || moduleFilter == "personal" {
		users, _ := a.eng.Store.ListUserCards()
		for _, uc := range users {
			key := fmt.Sprintf("user:%d", uc.ID)
			c := &Card{
				Key: key, Module: "personal", ExID: "personal", Title: "Personal card",
				Q: uc.Q, A: uc.A, New: true, Reviewable: true,
			}
			if rs := states[key]; rs != nil {
				c.New = false
				c.Ease = rs.Ease
				c.Interval = rs.Interval
				c.Due = rs.Due.UTC().Format("2006-01-02")
				c.Reps = rs.Reps
				c.Lapses = rs.Lapses
				c.Reviewable = !now.Before(rs.Due)
			}
			cards = append(cards, c)
		}
	}
	return cards
}

// AnswerCard applies an SM-2 rating (0 again, 1 hard, 2 good, 3 easy) to the
// card and returns its updated state.
func (a *App) AnswerCard(key string, rating int) (*Card, error) {
	if strings.HasPrefix(key, "user:") {
		return a.answerUserCard(key, rating)
	}
	exID, idx, err := parseCardKey(key)
	if err != nil {
		return nil, err
	}
	ex, err := a.findExercise(exID)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(ex.Answers) {
		return nil, fmt.Errorf("bad card index %d", idx)
	}
	qa := ex.Answers[idx]

	now := time.Now().UTC()
	var rs *store.ReviewState
	if existing, _ := a.eng.Store.GetReview(key); existing != nil {
		rs = existing
	} else {
		rs = &store.ReviewState{Key: key, Ease: 2.5}
		_ = a.eng.Store.StartReview(key)
	}
	applySM2(rs, rating, now)
	_ = a.eng.Store.SaveReview(rs)

	return &Card{
		Key: key, ExID: exID, Module: ex.Module, Title: ex.Title,
		Q: qa.Q, A: qa.A,
		Ease: rs.Ease, Interval: rs.Interval, Due: rs.Due.UTC().Format("2006-01-02"),
		Reps: rs.Reps, Lapses: rs.Lapses,
	}, nil
}

// applySM2 is a small SM-2 variant: ratings are 0 again / 1 hard / 2 good /
// 3 easy. Failures shrink the ease factor and send the card back to a
// relearn step; successes grow the interval toward ~ease×interval.
func applySM2(rs *store.ReviewState, rating int, now time.Time) {
	const minEase = 1.3
	switch rating {
	case 0: // again — relearn
		rs.Ease = math.Max(minEase, rs.Ease-0.20)
		rs.Interval = 0
		rs.Reps = 0
		rs.Lapses++
		rs.Due = now.Add(10 * time.Minute)
	case 1: // hard — small step up, slight ease penalty
		if rs.Interval < 1 {
			rs.Interval = 1
		}
		rs.Interval = math.Max(rs.Interval*1.2, rs.Interval+1)
		rs.Ease = math.Max(minEase, rs.Ease-0.15)
		rs.Reps++
		rs.Due = addDays(now, rs.Interval)
	case 2: // good — classic SM-2 steps
		rs.Reps++
		rs.Interval = sm2Interval(rs.Reps, rs.Interval, rs.Ease)
		rs.Due = addDays(now, rs.Interval)
	case 3: // easy — bonus interval + ease boost (cap 3.0)
		rs.Reps++
		rs.Interval = sm2Interval(rs.Reps, rs.Interval, rs.Ease) * 1.3
		rs.Ease = math.Min(3.0, rs.Ease+0.15)
		rs.Due = addDays(now, rs.Interval)
	}
	rs.LastReview = now
}

func sm2Interval(reps int, prev, ease float64) float64 {
	switch reps {
	case 1:
		return 1
	case 2:
		return 6
	default:
		if prev < 1 {
			return 1
		}
		return prev * ease
	}
}

func addDays(t time.Time, days float64) time.Time {
	if days >= 1 {
		return t.AddDate(0, 0, int(days))
	}
	return t.Add(time.Duration(days * float64(24*time.Hour)))
}

func parseCardKey(key string) (string, int, error) {
	i := strings.LastIndexByte(key, ':')
	if i <= 0 {
		return "", 0, fmt.Errorf("bad card key %q", key)
	}
	idx, err := strconv.Atoi(key[i+1:])
	if err != nil {
		return "", 0, err
	}
	return key[:i], idx, nil
}