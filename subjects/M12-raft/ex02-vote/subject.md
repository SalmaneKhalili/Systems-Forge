# M12-ex02 · Vote

The monotonic term from ex01 orders ballots, but freshness decides which
candidate may receive one. This exercise makes a follower grant a vote only to
a candidate whose last log entry is at least as up-to-date as its own.

## Shape

A stale node must not overtake a fresh one. Compare logs by the candidate's
last term, falling back to the last index when the terms tie. A node that has
already voted in the term refuses every later candidate. You write
**`vote.go`** and implement:

```go
type LogInfo struct {
	Term int // term of the last entry
	Idx  int // index of the last entry
}

type Voter struct {
	term  int
	voted bool
}

// NewVoter starts at term 0 with no vote cast.
func NewVoter() *Voter

// RequestVote handles a candidate's RequestVote at term candTerm whose log
// ends at cand LogInfo. It advances our term if candTerm is newer, then
// grants the vote if ALL hold: the candidate's term is not stale, we have not
// voted, and the candidate's log is at least as up-to-date as ours. Returns
// true if the vote is granted.
func (v *Voter) RequestVote(candTerm int, cand LogInfo, mine LogInfo) bool
```

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized. For the up-to-date comparison, a candidate last-term
above the node's is fresher; when terms are equal, the candidate's last index
must be greater than or equal to the node's. Never read the wall clock.

## Acceptance

The reference transcript is:

```text
case a: term1 idx1 vs term1 idx0  -> true
case b: term1 idx0 vs term1 idx1  -> false
case c: term2 idx0 vs term1 idx9  -> true   (higher term wins regardless of index)
case d: term1 idx5 vs term0 idx8  -> true   (equal? no: ours term0, cand term1 -> fresher)
case e: already voted             -> false
```

- A candidate with a younger last-term loses even if it has more entries, as
  case d shows with the terms swapped.
- A candidate trailing in index at the same term loses, as case b shows.
- Granting on length alone, or granting twice, fails the election restriction.

## Readings

- **Reading ladder** — the vote rule is the election restriction that protects a fresh log;
  the guide shows it, the paper states it.
- Raft paper, §5.2 "Leader election" (the up-to-date-log check): https://raft.github.io/raft.pdf
- Raft visual guide: https://raft.github.io/
