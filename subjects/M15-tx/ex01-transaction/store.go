package main

import "sort"

// Store is an in-memory key-value store. Writes outside a transaction apply
// immediately; writes inside a transaction are staged until Commit.
type Store struct {
	data map[string]string
}

// Tx is an in-flight transaction: writes are staged and only become visible
// to other readers when the transaction commits.
type Tx struct {
	store  *Store
	staged map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Set(k, v string) {
	s.data[k] = v
}

func (s *Store) Get(k string) string {
	return s.data[k]
}

// Snapshot returns a sorted copy of the committed keys (used by tests).
func (s *Store) Snapshot() []string {
	var out []string
	for k, v := range s.data {
		if v != "" {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func (s *Store) Begin() *Tx {
	return &Tx{store: s, staged: make(map[string]string)}
}

// Set stages a write visible to this transaction (and committed on Commit).
func (t *Tx) Set(k, v string) {
	t.staged[k] = v
}

// Get reads this transaction's own staged write if present, else committed.
func (t *Tx) Get(k string) string {
	if v, ok := t.staged[k]; ok {
		return v
	}
	return t.store.data[k]
}

// Commit applies all staged writes to the store atomically.
func (t *Tx) Commit() {
	for k, v := range t.staged {
		t.store.data[k] = v
	}
	t.staged = make(map[string]string)
}

// Rollback discards all staged writes without touching the store.
func (t *Tx) Rollback() {
	t.staged = make(map[string]string)
}
