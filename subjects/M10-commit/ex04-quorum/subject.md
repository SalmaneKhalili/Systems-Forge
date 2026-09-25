# M10-ex04 · Quorum

A coordinator does not need every unreliable participant to agree; it needs a
**quorum**. M10 reduces that requirement to strict-majority arithmetic, where
any two majorities overlap and opposing decisions cannot both reach a quorum.

## Shape

The safest general rule is a strict majority, `floor(n/2) + 1`. Any two
majority quorums **overlap**, so two opposing decisions cannot both reach a
majority; that overlap keeps the group safe against split-brain. You write
**`quorum.go`** and implement:

```go
// QuorumSize returns the smallest strict-majority quorum for n participants.
func QuorumSize(n int) int

// HasQuorum reports whether yes votes reach a strict majority of n.
func HasQuorum(n, yes int) bool
```

The exercise ships `main.go`, which prints the quorum size and several
yes-vote checks for a few group sizes. `make all` must build `test`, and
`./test` must print the reference transcript exactly. Use Go and the standard
library only. The reference transcript is `expected.txt`, with whitespace
normalized; never read the wall clock.

## Acceptance

The reference transcript is:

```text
n=1 quorum 1  has 1? true
n=2 quorum 2  has 1? false  has 2? true
n=3 quorum 2  has 1? false  has 2? true
n=4 quorum 3  has 2? false  has 3? true
n=5 quorum 3  has 2? false  has 3? true
```

- A 2-node group needs **both** nodes, so quorum is 2; a 1-1 tie is not a
  quorum.
- A 4-node group needs 3, not 2; 2-2 is a tie, not a majority.
- Using `floor(n/2)` (or `n/2`) for quorum is the classic bug because it lets
  a tie count as a decision.

## Readings

- **Reading ladder** — a quorum is a fact about *counts*; majority arithmetic does the work.
- *Designing Data-Intensive Applications*, Chapter 9 — quorums and majority reads/writes.
- Wikipedia, "Quorum (distributed computing)":
  https://en.wikipedia.org/wiki/Quorum_(distributed_computing)
