package store

import (
	"time"
)

// ReviewState is the SM-2 spaced-repetition state for one review card.
// Card keys have the shape "<exerciseID>:<index>" (see gui app).
type ReviewState struct {
	Key         string    // "<exID>:<n>"
	Ease        float64   // SM-2 ease factor (2.5 base)
	Interval    float64   // interval in days
	Due         time.Time // next review time
	LastReview  time.Time // when it was last reviewed (zero when never)
	Reps        int       // number of successful reviews
	Lapses      int       // number of times the card lapsed ("again")
}

// StartReview inserts an initial review-state row for a newly seen card.
func (s *Store) StartReview(key string) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO review_state (key, ease, interval_days, due, last_review, reps, lapses)
		 VALUES (?, 2.5, 0, ?, ?, 0, 0)`,
		key, time.Now().UTC().Format(time.RFC3339Nano), time.Time{}.Format(time.RFC3339Nano))
	return err
}

// GetReview fetches the review state for a card (nil when it has no state yet).
func (s *Store) GetReview(key string) (*ReviewState, error) {
	row := s.db.QueryRow(
		`SELECT key, ease, interval_days, due, last_review, reps, lapses
		   FROM review_state WHERE key = ?`, key)
	rs := new(ReviewState)
	var due, last string
	err := row.Scan(&rs.Key, &rs.Ease, &rs.Interval, &due, &last, &rs.Reps, &rs.Lapses)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	rs.Due, _ = time.Parse(time.RFC3339Nano, due)
	rs.LastReview, _ = time.Parse(time.RFC3339Nano, last)
	return rs, nil
}

// SaveReview writes the SM-2 state for a card.
func (s *Store) SaveReview(rs *ReviewState) error {
	_, err := s.db.Exec(
		`INSERT INTO review_state (key, ease, interval_days, due, last_review, reps, lapses)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		   ease = excluded.ease,
		   interval_days = excluded.interval_days,
		   due = excluded.due,
		   last_review = excluded.last_review,
		   reps = excluded.reps,
		   lapses = excluded.lapses`,
		rs.Key, rs.Ease, rs.Interval, rs.Due.UTC().Format(time.RFC3339Nano),
		rs.LastReview.UTC().Format(time.RFC3339Nano), rs.Reps, rs.Lapses)
	return err
}

// AllReviews returns every stored review state, keyed by card key.
func (s *Store) AllReviews() (map[string]*ReviewState, error) {
	rows, err := s.db.Query(
		`SELECT key, ease, interval_days, due, last_review, reps, lapses FROM review_state`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]*ReviewState{}
	for rows.Next() {
		rs := new(ReviewState)
		var due, last string
		if err := rows.Scan(&rs.Key, &rs.Ease, &rs.Interval, &due, &last, &rs.Reps, &rs.Lapses); err != nil {
			return nil, err
		}
		rs.Due, _ = time.Parse(time.RFC3339Nano, due)
		rs.LastReview, _ = time.Parse(time.RFC3339Nano, last)
		out[rs.Key] = rs
	}
	return out, rows.Err()
}

func isNoRows(err error) bool {
	return err != nil && err.Error() == "sql: no rows in result set"
}