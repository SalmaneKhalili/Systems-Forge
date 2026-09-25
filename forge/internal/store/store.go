// Package store records grading history and per-exercise progress in a
// local SQLite database (cgo-free driver).
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Status values recorded for a graded run.
const (
	StatusPass = "pass"
	StatusFail = "fail"
	StatusErr  = "error"
)

// Run is one grading attempt.
type Run struct {
	ID         int64
	Exercise   string
	Status     string
	Stdout     string
	Stderr     string
	Detail     string // JSON summary of method-level results
	DurationMs int64
	RanAt      time.Time
}

// Store is the progress database.
type Store struct {
	db *sql.DB
}

// Open opens (creating if needed) the progress database.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS runs (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		exercise   TEXT    NOT NULL,
		status     TEXT    NOT NULL,
		stdout     TEXT    NOT NULL DEFAULT '',
		stderr     TEXT    NOT NULL DEFAULT '',
		detail     TEXT    NOT NULL DEFAULT '',
		duration_ms INTEGER NOT NULL DEFAULT 0,
		ran_at     TEXT    NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_runs_exercise ON runs(exercise, id);

	CREATE TABLE IF NOT EXISTS review_state (
		key          TEXT PRIMARY KEY,
		ease         REAL    NOT NULL DEFAULT 2.5,
		interval_days REAL   NOT NULL DEFAULT 0,
		due          TEXT    NOT NULL,
		last_review  TEXT    NOT NULL,
		reps         INTEGER NOT NULL DEFAULT 0,
		lapses       INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		kind       TEXT    NOT NULL,
		exercise   TEXT    NOT NULL DEFAULT '',
		started_at TEXT    NOT NULL,
		ended_at   TEXT    NOT NULL,
		minutes    INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_sessions_started ON sessions(started_at);

	CREATE TABLE IF NOT EXISTS settings (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS user_cards (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		q            TEXT    NOT NULL,
		a            TEXT    NOT NULL,
		created_at   TEXT    NOT NULL
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

// RecordRun inserts a grading attempt.
func (s *Store) RecordRun(r *Run) error {
	_, err := s.db.Exec(
		`INSERT INTO runs (exercise, status, stdout, stderr, detail, duration_ms, ran_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		r.Exercise, r.Status, r.Stdout, r.Stderr, r.Detail, r.DurationMs,
		r.RanAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

// LastRun fetches the most recent run for an exercise (nil when none).
func (s *Store) LastRun(exercise string) (*Run, error) {
	row := s.db.QueryRow(
		`SELECT id, exercise, status, stdout, stderr, detail, duration_ms, ran_at
		   FROM runs WHERE exercise = ?
		  ORDER BY id DESC LIMIT 1`, exercise)
	return scanRun(row)
}

// RunHistory returns runs for an exercise, most recent first.
func (s *Store) RunHistory(exercise string, limit int) ([]*Run, error) {
	rows, err := s.db.Query(
		`SELECT id, exercise, status, stdout, stderr, detail, duration_ms, ran_at
		   FROM runs WHERE exercise = ? ORDER BY id DESC LIMIT ?`, exercise, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Run
	for rows.Next() {
		r, err := scanRun(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ModuleProgress aggregates pass/fail across a list of exercise IDs.
func (s *Store) ModuleProgress(exercises []string) (passed, total int, err error) {
	total = len(exercises)
	if total == 0 {
		return 0, 0, nil
	}
	for _, ex := range exercises {
		last, err := s.LastRun(ex)
		if err != nil {
			return 0, 0, err
		}
		if last != nil && last.Status == StatusPass {
			passed++
		}
	}
	return passed, total, nil
}

// OverallProgress computes stats over all known exercises.
func (s *Store) OverallProgress(all []string) (passed, attempted, total int, err error) {
	total = len(all)
	for _, ex := range all {
		last, lerr := s.LastRun(ex)
		if lerr != nil {
			return 0, 0, 0, lerr
		}
		if last != nil {
			attempted++
			if last.Status == StatusPass {
				passed++
			}
		}
	}
	return passed, attempted, total, nil
}

// ActivityDay is the number of graded runs on a UTC calendar day.
type ActivityDay struct {
	Day   time.Time // UTC midnight
	Count int
}

// ActivityByDay returns the total number of runs per UTC day within the
// [start, end] window (inclusive), for the contribution heatmap. Days with no
// runs are omitted; callers fill gaps from the window.
func (s *Store) ActivityByDay(start, end time.Time) (map[string]int, error) {
	rows, err := s.db.Query(
		`SELECT ran_at FROM runs WHERE ran_at >= ? AND ran_at <= ?`,
		start.UTC().Format(time.RFC3339Nano), end.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dayKey := func(t time.Time) string { return t.UTC().Format("2006-01-02") }
	out := map[string]int{}
	for rows.Next() {
		var ranAt string
		if err := rows.Scan(&ranAt); err != nil {
			return nil, err
		}
		t, perr := time.Parse(time.RFC3339Nano, ranAt)
		if perr != nil {
			continue
		}
		out[dayKey(t)]++
	}
	return out, rows.Err()
}

// PassesOnDay counts grading runs recorded as "pass" within the UTC calendar
// day containing day.
func (s *Store) PassesOnDay(day time.Time) (int, error) {
	start := day.UTC()
	end := start.AddDate(0, 0, 1)
	var n int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM runs WHERE status = ? AND ran_at >= ? AND ran_at < ?`,
		StatusPass, start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano)).Scan(&n)
	return n, err
}

// ReviewedOnDay counts review_state rows last reviewed within the UTC
// calendar day containing day. Rows with a zero last_review (never reviewed)
// fall outside any real day range and are naturally excluded.
func (s *Store) ReviewedOnDay(day time.Time) (int, error) {
	start := day.UTC()
	end := start.AddDate(0, 0, 1)
	var n int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM review_state WHERE last_review >= ? AND last_review < ?`,
		start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano)).Scan(&n)
	return n, err
}

func scanRun(row interface{ Scan(...any) error }) (*Run, error) {
	r := new(Run)
	var ranAt string
	err := row.Scan(&r.ID, &r.Exercise, &r.Status, &r.Stdout, &r.Stderr,
		&r.Detail, &r.DurationMs, &ranAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	r.RanAt, _ = time.Parse(time.RFC3339Nano, ranAt)
	return r, nil
}

// Close closes the database.
func (s *Store) Close() error {
	return s.db.Close()
}
