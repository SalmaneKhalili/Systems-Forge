# M11-ex04 · Snapshot

## Goal

A log that only ever grows is unbounded. Replicas compact it with a
**snapshot**: the state machine state at some index is captured, and everything
at or below that index is dropped from the in-memory log. The committed prefix
is replaced by the snapshot; only the uncommitted suffix stays as raw entries.

Implement `log.go`:

```go
type Log struct {
	// entries holds ONLY the uncommitted suffix after compaction.
	entries []string
	base    int // index the first entry is at (0 when never compacted)
}

// Append adds cmd at the end; returns Log's public index.
func (l *Log) Append(cmd string) int

// Snapshot drops everything with index <= upto; returns the number of
// entries removed (on-disk snapshot money is out of scope).
func (l *Log) Snapshot(upto int) int

// All returns entries with their public indexes preserved.
func (l *Log) All() []string
```

After a snapshot the log no longer knows the compacted entries, but the
*indexes* of the remaining ones stay stable. Use a `base` offset so `All`
keeps reporting original indexes.

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `log.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `All` must preserve each entry's pre-snapshot index.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
append a -> 0   append b -> 1   append c -> 2   append d -> 3
all=0:a|1:b|2:c|3:d
snapshot(1) -> removed 2
all=2:c|3:d
append e -> 4
all=2:c|3:d|4:e
```

After `snapshot(1)`, entries 0 and 1 are gone but `c` is still reported at
index 2 and `d` at 3; the next append lands at 4. A log that renumbers the
tail after compaction (or that fails to drop the compacted prefix) is the bug.