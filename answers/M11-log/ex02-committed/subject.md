# M11-ex02 · Commit index

## Goal

Not every log entry can be applied to the replicated state machine. Only the
entries at or below the **commit index** are durable and safe to apply; the
ones above it are tentative and may still be rolled back by a leader change.
Track the commit index alongside the log.

Implement `log.go`:

```go
type Log struct {
	entries []string
	commit  int // highest committed index, or -1 when none
}

// Append adds cmd at the end and returns its index.
func (l *Log) Append(cmd string) int

// Commit raises the commit index to idx and returns it.
func (l *Log) Commit(idx int) int

// Committed returns the entries at or below the commit index.
func (l *Log) Committed() []string
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `log.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `Committed` returns only the durable prefix.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
append a -> 0
append b -> 1
append c -> 2
committed=[] commit=-1
commit 1 -> 1
committed=a|b
commit 2 -> 2
committed=a|b|c
```

`c` sits beyond the commit index at first, so `committed` is empty; after
`commit 1`, only `a` and `b` are durable; `c` becomes committed only after
the commit index reaches 2.