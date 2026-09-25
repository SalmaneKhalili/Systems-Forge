package store

import (
	"fmt"
	"time"
)

// UserCard is a learner-authored flashcard. SM-2 scheduling state lives in
// review_state under the synthetic key "user:<id>" so existing review logic
// and persistence apply unchanged.
type UserCard struct {
	ID        int64
	Q         string
	A         string
	CreatedAt time.Time
}

// AddUserCard inserts a personal card and returns its id.
func (s *Store) AddUserCard(q, a string) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO user_cards (q, a, created_at) VALUES (?, ?, ?)`,
		q, a, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UserCard fetches one personal card (nil when missing).
func (s *Store) UserCard(id int64) (*UserCard, error) {
	row := s.db.QueryRow(
		`SELECT id, q, a, created_at FROM user_cards WHERE id = ?`, id)
	uc := new(UserCard)
	var created string
	err := row.Scan(&uc.ID, &uc.Q, &uc.A, &created)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	uc.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return uc, nil
}

// ListUserCards returns every personal card, newest first.
func (s *Store) ListUserCards() ([]*UserCard, error) {
	rows, err := s.db.Query(
		`SELECT id, q, a, created_at FROM user_cards ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*UserCard
	for rows.Next() {
		uc := new(UserCard)
		var created string
		if err := rows.Scan(&uc.ID, &uc.Q, &uc.A, &created); err != nil {
			return nil, err
		}
		uc.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		out = append(out, uc)
	}
	return out, rows.Err()
}

// UpdateUserCard rewrites a personal card's contents.
func (s *Store) UpdateUserCard(id int64, q, a string) error {
	_, err := s.db.Exec(
		`UPDATE user_cards SET q = ?, a = ? WHERE id = ?`, q, a, id)
	return err
}

// DeleteUserCard removes the card and its SM-2 state.
func (s *Store) DeleteUserCard(id int64) error {
	if _, err := s.db.Exec(
		`DELETE FROM review_state WHERE key = ?`, fmt.Sprintf("user:%d", id)); err != nil {
		return err
	}
	_, err := s.db.Exec(`DELETE FROM user_cards WHERE id = ?`, id)
	return err
}