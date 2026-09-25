# M11-ex04 · Snapshot

A log that only grows is unbounded. This exercise compacts it with a
**snapshot**, replacing the committed prefix with state captured at an index
while retaining the uncommitted suffix and every surviving entry's original
public index.

## Shape

The state machine state at some index is captured, and everything at or below
that index is dropped from the in-memory log. The committed prefix is replaced
by the snapshot; only the uncommitted suffix remains as raw entries. You write
**`log.go`**:

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

After compaction the log no longer knows the removed entries, but the
*indexes* of the remaining ones stay stable. Use a `base` offset so `All`
continues to report original indexes.

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized. `All` must preserve every entry's pre-snapshot index;
never read the wall clock.

## Acceptance

The reference transcript is:

```text
append a -> 0   append b -> 1   append c -> 2   append d -> 3
all=0:a|1:b|2:c|3:d
snapshot(1) -> removed 2
all=2:c|3:d
append e -> 4
all=2:c|3:d|4:e
```

After `snapshot(1)`, entries 0 and 1 are gone, while `c` remains at public
index 2 and `d` at 3; the next append lands at 4. Renumbering the tail after
compaction, or failing to drop the compacted prefix, fails the contract.

## Readings

- **Reading ladder** — a snapshot is a checkpoint of applied state; the Raft compaction note
  and The Log each frame the bounded-memory trade-off.
- Raft paper (log compaction section): https://raft.github.io/raft.pdf
- Jay Kreps, *The Log* — log compaction and retention.
- *Designing Data-Intensive Applications*, Chapter 11 "Stream Processing" — log compaction as
  a design pattern.
