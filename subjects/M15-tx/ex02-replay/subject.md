# M15-ex02 · Replay

## Goal

A **write-ahead log** is the backbone of durability. Every intended change is
appended to a log first; a crash is recovered by replaying the log. Replay must
be deterministic and *last-writer-wins*: a later `set` on the same key
overwrites the earlier one.

Implement `log.go`:

```go
type Log struct{ entries []string }

func NewLog() *Log
func (l *Log) Append(op string)          // append "set k=v" to the log
func (l *Log) Entries() []string         // a copy of the entries
func (l *Log) Replay() map[string]string // replay ops in order -> key-value state
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `log.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Each op is `set k=v`; replay applies op 0..n-1 in order, later wins.
- No wall clock; the log is a pure in-memory model.

## Acceptance

Reference transcript:

```
entries -> 3
a -> 3
b -> 2
```

After appending `set a=1`, `set b=2`, `set a=3`, replay must give `a=3`
(last write wins) and `b=2`; entries length is 3. A replay that applies
out-of-order or lets an earlier write to `a` win is the bug.

## Readings

- **Reading ladder** — write the intent first, replay on crash; the WAL docs and DDIA agree
  on the order.
- PostgreSQL, "Write-Ahead Logging (WAL)": https://www.postgresql.org/docs/current/wal-intro.html
- Wikipedia, "Write-ahead logging": https://en.wikipedia.org/wiki/Write-ahead_logging
