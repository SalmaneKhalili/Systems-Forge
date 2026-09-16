# M11-ex03 · Index

## Goal

Replicas shuffle the same numbers constantly: commit index, log length, next
write slot, and how many entries are still uncommitted. Get them right and
the append/commit dance is trivial; mix up `len-1` for `len` and a leader
silently skips or duplicates an index.

Implement `index.go` as pure arithmetic over an (index, length) pair:

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

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `index.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Never read the wall clock.

## Acceptance

Reference transcript:

```
next(c=-1,l=0) = 0      uncommitted = 0
next(c=0,l=3)  = 3      uncommitted = 2
next(c=2,l=3)  = 3      uncommitted = 0
next(c=1,l=5)  = 5      uncommitted = 3
```

Empty log: next write is 0, nothing uncommitted. With 3 entries and commit
at 0, two entries (`1`,`2`) are uncommitted. The next write is always the
*current length*, never `len-1`.