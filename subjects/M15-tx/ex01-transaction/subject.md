# M15-ex01 · Transaction

## Goal

A **transaction** groups several writes so they either all take effect
(commit) or none do (rollback). This exercise models a tiny in-memory key-value
store where writes inside a transaction are **staged** and only become visible
to other readers on commit. `"6"` writes nothing dangling.

Implement `store.go`:

```go
type Store struct{ data map[string]string }

func NewStore() *Store
func (s *Store) Set(k, v string)          // apply immediately (outside any txn)
func (s *Store) Get(k string) string      // committed value ("" if absent)

type Tx struct{ ... }

func (s *Store) Begin() *Tx               // start a transaction
func (t *Tx) Set(k, v string)             // stage a write (own txn only)
func (t *Tx) Get(k string) string         // own staged write if present, else committed
func (t *Tx) Commit()                     // apply all staged writes atomically
func (t *Tx) Rollback()                   // discard all staged writes
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `store.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Staged writes must NOT be visible to `Store.Get` until `Commit`.
- `Rollback` must discard staged writes without touching the store.
- No wall clock; a transaction is a pure data model driven by the harness.

## Acceptance

Reference transcript:

```
store.get a -> 1
tx.get a -> 2
store.get a -> 1
rollback done
store.get a -> 1
commit done
store.get a -> 5
store.get b -> 6
```

`m.Snapshot` is used by the grader runner only; a transaction that loses its
staged writes on `Commit`, or that applies rollback's writes anyway, is the bug.
