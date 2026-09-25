package store

import (
	"path/filepath"
	"testing"
)

// TestUserCardCRUD exercises the personal-card table end-to-end on a throw
// away database, including the SM-2 state cleanup on delete.
func TestUserCardCRUD(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "progress.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	id, err := s.AddUserCard("What is EINTR?", "Retry the syscall.")
	if err != nil {
		t.Fatalf("AddUserCard: %v", err)
	}
	if id <= 0 {
		t.Fatalf("AddUserCard: bad id %d", id)
	}

	// The card key used across the app is "user:<id>"; its SM-2 state row
	// must exist and be mergeable with the deck.
	if err := s.StartReview("user:" + itoa(id)); err != nil {
		t.Fatalf("StartReview: %v", err)
	}
	rs, err := s.GetReview("user:" + itoa(id))
	if err != nil || rs == nil {
		t.Fatalf("GetReview: rs=%v err=%v", rs, err)
	}

	uc, err := s.UserCard(id)
	if err != nil || uc == nil {
		t.Fatalf("UserCard: uc=%v err=%v", uc, err)
	}
	if uc.Q != "What is EINTR?" {
		t.Fatalf("UserCard: wrong q %q", uc.Q)
	}

	if err := s.UpdateUserCard(id, "What is EPIPE?", "Broken pipe — handle SIGPIPE or ignore it."); err != nil {
		t.Fatalf("UpdateUserCard: %v", err)
	}
	uc, _ = s.UserCard(id)
	if uc.Q != "What is EPIPE?" {
		t.Fatalf("UpdateUserCard: q not updated %q", uc.Q)
	}

	all, err := s.ListUserCards()
	if err != nil || len(all) != 1 {
		t.Fatalf("ListUserCards: %d cards err=%v", len(all), err)
	}

	// Delete must remove both the card and its review_state so no orphan rows
	// linger in the deck.
	if err := s.DeleteUserCard(id); err != nil {
		t.Fatalf("DeleteUserCard: %v", err)
	}
	uc, _ = s.UserCard(id)
	if uc != nil {
		t.Fatal("DeleteUserCard: card still present")
	}
	rs, err = s.GetReview("user:" + itoa(id))
	if err != nil || rs != nil {
		t.Fatalf("DeleteUserCard: review state not removed rs=%v err=%v", rs, err)
	}
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}