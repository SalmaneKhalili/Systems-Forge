package store

import (
	"time"
)

// Session is one tracked study block (pomodoro or free focus session).
type Session struct {
	ID        int64
	Kind      string // "pomodoro" | "focus"
	Exercise  string // exercise id the session was attached to (may be "")
	StartedAt time.Time
	EndedAt   time.Time // zero while the session is still open
	Minutes   int
}

// StartSession opens a new session row.
func (s *Store) StartSession(kind, exercise string) (*Session, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO sessions (kind, exercise, started_at, ended_at, minutes)
		 VALUES (?, ?, ?, ?, 0)`,
		kind, exercise, now.Format(time.RFC3339Nano), time.Time{}.Format(time.RFC3339Nano))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Session{ID: id, Kind: kind, Exercise: exercise, StartedAt: now}, nil
}

// EndSession closes an open session row if one exists for the id.
func (s *Store) EndSession(id int64, minutes int) (*Session, error) {
	now := time.Now().UTC()
	_, err := s.db.Exec(
		`UPDATE sessions SET ended_at = ?, minutes = ? WHERE id = ?`,
		now.Format(time.RFC3339Nano), minutes, id)
	if err != nil {
		return nil, err
	}
	return s.Session(id)
}

// Session fetches one session by id.
func (s *Store) Session(id int64) (*Session, error) {
	row := s.db.QueryRow(
		`SELECT id, kind, exercise, started_at, ended_at, minutes FROM sessions WHERE id = ?`, id)
	return scanSession(row)
}

// RecentSessions returns the most recent sessions.
func (s *Store) RecentSessions(limit int) ([]*Session, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Query(
		`SELECT id, kind, exercise, started_at, ended_at, minutes
		   FROM sessions ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Session
	for rows.Next() {
		ss, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ss)
	}
	return out, rows.Err()
}

// MinutesBetween sums focused minutes in [start, end]. When exercise is
// non-empty only sessions attached to that exercise count.
func (s *Store) MinutesBetween(start, end time.Time, exercise string) (int, error) {
	q := `SELECT COALESCE(SUM(minutes), 0) FROM sessions
	       WHERE started_at >= ? AND started_at <= ?`
	args := []interface{}{start.UTC().Format(time.RFC3339Nano), end.UTC().Format(time.RFC3339Nano)}
	if exercise != "" {
		q += ` AND exercise = ?`
		args = append(args, exercise)
	}
	var total int
	err := s.db.QueryRow(q, args...).Scan(&total)
	return total, err
}

// MinutesByDay returns focused minutes per UTC day within [start, end].
func (s *Store) MinutesByDay(start, end time.Time) (map[string]int, error) {
	rows, err := s.db.Query(
		`SELECT started_at, minutes FROM sessions
		  WHERE started_at >= ? AND started_at <= ? AND minutes > 0`,
		start.UTC().Format(time.RFC3339Nano), end.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var at string
		var mins int
		if err := rows.Scan(&at, &mins); err != nil {
			return nil, err
		}
		t, perr := time.Parse(time.RFC3339Nano, at)
		if perr != nil {
			continue
		}
		out[t.UTC().Format("2006-01-02")] += mins
	}
	return out, rows.Err()
}

// MinutesByExercise returns focused minutes grouped by exercise.
func (s *Store) MinutesByExercise() (map[string]int, error) {
	rows, err := s.db.Query(
		`SELECT COALESCE(exercise, ''), SUM(minutes) FROM sessions
		  WHERE minutes > 0 AND exercise != '' GROUP BY exercise`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var ex string
		var total int
		if err := rows.Scan(&ex, &total); err != nil {
			return nil, err
		}
		out[ex] = total
	}
	return out, rows.Err()
}

func scanSession(row interface{ Scan(...any) error }) (*Session, error) {
	ss := new(Session)
	var started, ended string
	err := row.Scan(&ss.ID, &ss.Kind, &ss.Exercise, &started, &ended, &ss.Minutes)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	ss.StartedAt, _ = time.Parse(time.RFC3339Nano, started)
	ss.EndedAt, _ = time.Parse(time.RFC3339Nano, ended)
	return ss, nil
}