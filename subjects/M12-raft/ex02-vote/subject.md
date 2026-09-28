# M12-ex02 · Vote

## Goal

A follower grants a vote only to a candidate whose log is **at least as
up-to-date** as its own — the rule that keeps a stale node from overtaking a
fresh one. A log is compared by its last term, falling back to its last index
when terms tie. A node that already voted this term refuses everyone.

Implement `vote.go`:

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

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `vote.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Up-to-date comparison: candidate's last-term > ours → fresher; equal term →
  candidate's last-index >= ours → at least as up-to-date.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
case a: term1 idx1 vs term1 idx0  -> true
case b: term1 idx0 vs term1 idx1  -> false
case c: term2 idx0 vs term1 idx9  -> true   (higher term wins regardless of index)
case d: term1 idx5 vs term0 idx8  -> true   (equal? no: ours term0, cand term1 -> fresher)
case e: already voted             -> false
```

The tell: a candidate with a younger last-term loses even if it has more
entries (case d with swapped terms), and a candidate trailing in index at the
same term loses (case b). Granting on length alone, or granting twice, is the
bug.

## Readings

- **Reading ladder** — the vote rule is the election restriction that protects a fresh log;
  the guide shows it, the paper states it.
- Raft paper, §5.2 "Leader election" (the up-to-date-log check): https://raft.github.io/raft.pdf
- Raft visual guide: https://raft.github.io/
