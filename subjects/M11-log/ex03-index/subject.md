# M11-ex03 · Index

With append and commit in place, M11 isolates the index relationships that
keep them aligned. Pure arithmetic maps the commit index and current length to
the next write slot and the uncommitted count, preventing skipped or duplicated
public indexes.

## Shape

Replicas constantly move among the commit index, log length, next write slot
and uncommitted count. Mixing up `len-1` for `len` makes a leader silently skip
or duplicate an index. You write **`index.go`** with pure arithmetic over an
(index, length) pair:

```go
// State holds the commit index and log length.
type State struct {
	Commit int // highest committed index, or -1
	Len    int // number of entries
}

// NextIndex returns the index the next Append writes at.
func NextIndex(s State) int

// Uncommitted returns how many entries are beyond the commit index.
func Uncommitted(s State) int
```

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized; never read the wall clock.

## Acceptance

The reference transcript is:

```text
next(c=-1,l=0) = 0      uncommitted = 0
next(c=0,l=3)  = 3      uncommitted = 2
next(c=2,l=3)  = 3      uncommitted = 0
next(c=1,l=5)  = 5      uncommitted = 3
```

- An empty log has next index 0 and no uncommitted entries.
- With three entries and commit at 0, entries `1` and `2` are uncommitted.
- The next write is always the *current length*, never `len-1`.

## Readings

- **Reading ladder** — two integers, seven relations; the Raft index rules supply the
  arithmetic the exercise checks.
- Raft paper, §5.3 "Log replication" (nextIndex / commitIndex):
  https://raft.github.io/raft.pdf
