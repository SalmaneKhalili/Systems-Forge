# M15-ex02 · Replay

Atomic commit defines what a transaction publishes; the write-ahead log defines
how that result survives a crash. This exercise records intended changes before
they become state, then reconstructs that state by replaying entries in order.
You deliver `log.go` with the deterministic last-writer-wins model the durable
store will persist.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`log.go`** using only the Go standard library and implement:

```go
type Log struct{ entries []string }

func NewLog() *Log
func (l *Log) Append(op string)          // append "set k=v" to the log
func (l *Log) Entries() []string         // a copy of the entries
func (l *Log) Replay() map[string]string // replay ops in order -> key-value state
```

Each op is `set k=v`. `Entries` returns a copy, and `Replay` applies entries
`0..n-1` in order so later writes win. Never read the wall clock; the log is a
pure in-memory model. The reference transcript is `expected.txt` with whitespace
normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
entries -> 3
a -> 3
b -> 2
```

After appending `set a=1`, `set b=2`, and `set a=3`, replay yields `a=3` because
the last write wins, plus `b=2`; the entry count is 3. Applying entries out of
order, or allowing the earlier `a=1` to win, breaks deterministic recovery.

## Readings

- **Reading ladder** — write the intent first, replay on crash; the WAL docs and DDIA agree
  on the order.
- PostgreSQL, "Write-Ahead Logging (WAL)": https://www.postgresql.org/docs/current/wal-intro.html
- Wikipedia, "Write-ahead logging": https://en.wikipedia.org/wiki/Write-ahead_logging
