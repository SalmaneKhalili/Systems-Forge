package main

import (
	"fmt"
	"strings"
	"time"

	"forge/internal/store"
)

// UserCardInfo is the personal-card row shown in the review manager, merged
// with its SM-2 scheduling state.
type UserCardInfo struct {
	ID       int64   `json:"id"`
	Q        string  `json:"q"`
	A        string  `json:"a"`
	Reps     int     `json:"reps"`
	Interval float64 `json:"interval"`
	Due      string  `json:"due"`
}

// AddCard creates a personal flashcard and registers its SM-2 state.
func (a *App) AddCard(q, ans string) (*Card, error) {
	q = strings.TrimSpace(q)
	ans = strings.TrimSpace(ans)
	if q == "" {
		return nil, fmt.Errorf("a card needs a question")
	}
	id, err := a.eng.Store.AddUserCard(q, ans)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	if err := a.eng.Store.StartReview(key); err != nil {
		return nil, err
	}
	return &Card{Key: key, Module: "personal", ExID: "personal", Title: "Personal card", Q: q, A: ans, New: true, Reviewable: true}, nil
}

// ListUserCards returns every personal card with its current scheduling.
func (a *App) ListUserCards() ([]*UserCardInfo, error) {
	cards, err := a.eng.Store.ListUserCards()
	if err != nil {
		return nil, err
	}
	states, _ := a.eng.Store.AllReviews()
	out := make([]*UserCardInfo, 0, len(cards))
	for _, uc := range cards {
		info := &UserCardInfo{ID: uc.ID, Q: uc.Q, A: uc.A}
		if rs := states[fmt.Sprintf("user:%d", uc.ID)]; rs != nil {
			info.Reps = rs.Reps
			info.Interval = rs.Interval
			info.Due = rs.Due.UTC().Format("2006-01-02")
		}
		out = append(out, info)
	}
	return out, nil
}

// UpdateCard edits a personal card's contents.
func (a *App) UpdateCard(id int64, q, ans string) error {
	q = strings.TrimSpace(q)
	ans = strings.TrimSpace(ans)
	if q == "" {
		return fmt.Errorf("a card needs a question")
	}
	return a.eng.Store.UpdateUserCard(id, q, ans)
}

// DeleteCard removes a personal card and its review state.
func (a *App) DeleteCard(id int64) error {
	return a.eng.Store.DeleteUserCard(id)
}

// answerUserCard applies an SM-2 rating to a personal card.
func (a *App) answerUserCard(key string, rating int) (*Card, error) {
	var id int64
	if _, err := fmt.Sscanf(key, "user:%d", &id); err != nil {
		return nil, fmt.Errorf("bad card key %q", key)
	}
	uc, err := a.eng.Store.UserCard(id)
	if err != nil {
		return nil, err
	}
	if uc == nil {
		return nil, fmt.Errorf("card not found")
	}
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
		Key: key, ExID: "personal", Module: "personal", Title: "Personal card",
		Q: uc.Q, A: uc.A,
		Ease: rs.Ease, Interval: rs.Interval, Due: rs.Due.UTC().Format("2006-01-02"),
		Reps: rs.Reps, Lapses: rs.Lapses,
	}, nil
}