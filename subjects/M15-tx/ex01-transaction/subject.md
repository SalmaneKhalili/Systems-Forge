# M15-ex01 · Transaction

The module begins by turning several key-value writes into one atomic unit. A
transaction stages its changes, reads its own pending values, and publishes them
to the shared store only on commit, so rollback leaves nothing dangling. You
deliver `store.go` with the in-memory store and transaction state needed by
replay, conflict handling, and the later gateway.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`store.go`** using only the Go standard library and implement:

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

Staged writes must not be visible to `Store.Get` until `Commit`, while `Tx.Get`
must prefer its own staged value. `Rollback` discards staged writes without
touching the store. Never read the wall clock; the harness drives a pure data
model. The reference transcript is `expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
store.get a -> 1
tx.get a -> 2
store.get a -> 1
rollback done
store.get a -> 1
commit done
store.get a -> 5
store.get b -> 6
```

`m.Snapshot` is used by the grader runner only. A transaction that loses staged
writes on `Commit` or applies rollback's writes anyway fails the atomicity
contract and leaves state dangling.

## Readings

- **Reading ladder** — stage, commit, or roll back; DDIA Chapter 7 is the definitive *why*.
- *Designing Data-Intensive Applications*, Chapter 7 "Transactions" — atomicity and commit.
- Wikipedia, "ACID": https://en.wikipedia.org/wiki/ACID
