# M12-ex04 · Election

## Goal

A node does not become leader by hoping — it wins an **election** by earning
a strict majority of the cluster's votes. The quorum of an `n`-node cluster is
`n/2 + 1` (integer division): a 2-node group needs both votes, a 3-node group
needs 2, a 4-node group needs 3.

Implement `elect.go`:

```go
// Quorum returns the number of votes needed to win in a group of size n.
func Quorum(n int) int

// HasMajority reports whether yes votes are enough to win in a group of size n.
func HasMajority(n, yes int) bool
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `elect.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Strict majority: `yes >= n/2 + 1`; a tie (`yes = n/2`) is NOT a win.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
n=1 quorum=1 has1=true
n=2 quorum=2 has1=false has2=true
n=3 quorum=2 has2=true
n=4 quorum=3 has2=false has3=true
n=5 quorum=3 has3=true
```

The tell: `n=2 quorum=2`, `n=4 quorum=3` — an even cluster needs one past the
half, so `has1=false` for n=2 and `has2=false` for n=4. A `floor(n/2)` quorum
(2-node winning on 1, 4-node winning on 2) is the bug.

## Readings

- **Reading ladder** — election by majority is a quorum, and you already trained that in
  M10-ex04; here it is lifted onto term-and-log.
- Raft paper, §5.2 "Leader election" (quorum, vote granting): https://raft.github.io/raft.pdf
- Reread M10-ex04 (majority quorum).
