# M11-ex02 · Commit index

The append-only log from ex01 can contain tentative entries as well as durable
ones. This exercise adds the commit index so the state machine applies only
the safe prefix and can still roll back everything above it after a leader
change.

## Shape

Only entries at or below the **commit index** are durable and safe to apply;
entries above it may still be rolled back. You write **`log.go`** and track the
boundary beside the log:

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

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized. `Committed` returns only the durable prefix; never read
the wall clock.

## Acceptance

The reference transcript is:

```text
append a -> 0
append b -> 1
append c -> 2
committed=[] commit=-1
commit 1 -> 1
committed=a|b
commit 2 -> 2
committed=a|b|c
```

- Initially `c` is beyond the commit index, so the committed prefix is empty.
- After `commit 1`, only `a` and `b` are durable; `c` becomes committed only
  when the commit index reaches 2.

## Readings

- **Reading ladder** — "committed" is a *boundary on the log*, not a state of the node; the
  Raft safety rules say when an entry is safe to apply.
- Raft paper, §5.3 (log replication) and §5.4.2 (commitment / election safety):
  https://raft.github.io/raft.pdf
- Jay Kreps, *The Log* — durability of committed offsets.
