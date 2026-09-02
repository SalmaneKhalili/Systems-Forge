# M10-ex04 · Quorum

## Goal

When lots of participants — some unreliable — must agree, a coordinator
doesn't need all of them, only a **quorum**. The safest general rule is a
strict majority: `floor(n/2) + 1`. The key property is that any two majority
quorums **overlap** — so two opposing decisions can never both reach a
majority. That is what keeps a group safe against split-brain.

Implement `quorum.go`:

```go
// QuorumSize returns the smallest strict-majority quorum for n participants.
func QuorumSize(n int) int

// HasQuorum reports whether yes votes reach a strict majority of n.
func HasQuorum(n, yes int) bool
```

The provided `main.go` prints, for a few group sizes, the quorum size and
whether a given count of yes-votes is a quorum. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `quorum.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Never read the wall clock.

## Acceptance

Reference transcript:

```
n=1 quorum 1  has 1? true
n=2 quorum 2  has 1? false  has 2? true
n=3 quorum 2  has 1? false  has 2? true
n=4 quorum 3  has 2? false  has 3? true
n=5 quorum 3  has 2? false  has 3? true
```

The tell: a 2-node group needs **both** nodes (quorum 2) — a bare tie of 1-1
is not a quorum; and a 4-node group needs 3, not 2 (2-2 is a tie, not a
majority). Getting quorum as `floor(n/2)` (or `n/2`) is the classic bug that
lets a tie be mistaken for a decision.