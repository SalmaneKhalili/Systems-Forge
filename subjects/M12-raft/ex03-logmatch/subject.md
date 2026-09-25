# M12-ex03 · Log-match

The election safety rules from ex02 still need follower logs to stay aligned
with the leader. This exercise implements Raft's **matching-prefix check**:
AppendEntries may add an entry only after a predecessor with the expected
term.

## Shape

The leader sends the entry before `prevIndex`. The follower accepts only if
its entry at `prevIndex` has the same term as `prevTerm`; otherwise it rejects
and the leader backs up to a lower index. You write **`match.go`** and
implement:

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

The follower's log holds `Entry{term, cmd}` at indexes 0,1,… . If `prevIndex`
is a valid index and `log[prevIndex].Term == prevTerm`, the new entry is
appended after it.

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized. A match at `prevIndex` with the right term appends;
anything else rejects. Never read the wall clock.

## Acceptance

The reference transcript is:

```text
log   =b2 c3 d4
match prev(0,1) -> b2 c3 d4 e5
match prev(0,1) -> b2 c3 d4 e5
mismatch prev(0,2) -> REJECT
```

If the follower's entry at the previous index does not carry the leader's
`prevTerm`, nothing is appended: the log remains untouched and `"REJECT"` is
returned so the leader can back up. Appending anyway breaks log matching.

## Readings

- **Reading ladder** — the Log Matching Property is the single rule that keeps replicas'
  logs in lockstep.
- Raft paper, §5.3 "Log replication" (Log Matching Property): https://raft.github.io/raft.pdf
