# M12-ex03 · Log-match

## Goal

Raft keeps every replica's log in lockstep with a single rule: the leader
sends the entry before `prevIndex`; the follower accepts only if its entry at
`prevIndex` has the same term as `prevTerm` — the **matching-prefix check**.
If it does not match, the follower rejects and the leader backs up to a lower
index.

Implement `match.go`:

```go
type Entry struct {
	Term int
	Cmd  string
}

// AppendEntries attempts to append ent after a follower entry that matches
// prevTerm at prevIndex. Returns the resulting log as term strings (or
// "REJECT" if no match).
func AppendEntries(log []Entry, prevIndex int, prevTerm int, ent Entry) []string
```

The follower's log holds `Entry{term, cmd}` at indexes 0,1,…  The check is:
if `prevIndex` is a valid index and `log[prevIndex].Term == prevTerm`, then
the new entry is appended after it.

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `match.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- A match at `prevIndex` with the right term appends; anything else rejects.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
log   =b2 c3 d4
match prev(0,1) -> b2 c3 d4 e5
match prev(0,1) -> b2 c3 d4 e5
mismatch prev(0,2) -> REJECT
```

The tell: when the follower's entry at the previous index does not carry the
leader's `prevTerm`, nothing is appended (the log is left untouched and
"REJECT" is returned) — the leader must back up. Appending anyway is the bug.

## Readings

- **Reading ladder** — the Log Matching Property is the single rule that keeps replicas'
  logs in lockstep.
- Raft paper, §5.3 "Log replication" (Log Matching Property): https://raft.github.io/raft.pdf
