package main

import (
	"time"

	"forge/internal/store"
)

// StartFocus opens a tracked study block (kind: "pomodoro" | "focus"), tied
// to an exercise when one is active.
func (a *App) StartFocus(kind, exID string) (*StudySession, error) {
	s, err := a.eng.Store.StartSession(kind, exID)
	if err != nil {
		return nil, err
	}
	return toStudySession(s), nil
}

// EndFocus closes an open study block and records its focused minutes.
func (a *App) EndFocus(id int64, minutes int) (*StudySession, error) {
	if minutes < 0 {
		minutes = 0
	}
	s, err := a.eng.Store.EndSession(id, minutes)
	if err != nil {
		return nil, err
	}
	return toStudySession(s), nil
}

// Stats returns total focused minutes, per-exercise totals, and recent
// sessions for the profile/stats views.
func (a *App) Stats() *Stats {
	st := &Stats{}
	now := time.Now().UTC()
	st.TotalMinutes, _ = a.eng.Store.MinutesBetween(time.Time{}, now, "")
	if st.TotalMinutes < 0 {
		st.TotalMinutes = 0
	}
	st.PerExercise, _ = a.eng.Store.MinutesByExercise()
	if st.PerExercise == nil {
		st.PerExercise = map[string]int{}
	}
	recent, _ := a.eng.Store.RecentSessions(20)
	st.Recent = make([]*StudySession, 0, len(recent))
	for _, s := range recent {
		st.Recent = append(st.Recent, toStudySession(s))
	}
	return st
}

// GetSettings returns the whole settings KV map.
func (a *App) GetSettings() map[string]string {
	m, _ := a.eng.Store.AllSettings()
	if m == nil {
		m = map[string]string{}
	}
	return m
}

// SaveSettings upserts several settings at once.
func (a *App) SaveSettings(kv map[string]string) error {
	return a.eng.Store.SetSettingsBulk(kv)
}

func toStudySession(s *store.Session) *StudySession {
	out := &StudySession{
		ID:       s.ID,
		Kind:     s.Kind,
		Exercise: s.Exercise,
		Minutes:  s.Minutes,
		StartedAt: s.StartedAt.UTC().Format(time.RFC3339),
	}
	out.Open = s.EndedAt.IsZero()
	return out
}