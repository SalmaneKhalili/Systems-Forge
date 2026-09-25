# M12-ex04 · Election

The term and vote rules from ex01–ex02 establish a candidate, but the cluster
still needs a majority before it can recognize a leader. This exercise applies
the strict-majority quorum that gates ex05's single-node and ex06's networked
elections.

## Shape

A node becomes leader by earning a strict majority of the cluster's votes,
not by declaring itself winner. The quorum for an `n`-node cluster is `n/2 + 1`
with integer division: a 2-node group needs both votes, a 3-node group needs 2
and a 4-node group needs 3. You write **`elect.go`** and implement:

```go
// Quorum returns the number of votes needed to win in a group of size n.
func Quorum(n int) int

// HasMajority reports whether yes votes are enough to win in a group of size n.
func HasMajority(n, yes int) bool
```

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized. A strict majority requires `yes >= n/2 + 1`; a tie at
`yes = n/2` is not a win. Never read the wall clock.

## Acceptance

The reference transcript is:

```text
n=1 quorum=1 has1=true
n=2 quorum=2 has1=false has2=true
n=3 quorum=2 has2=true
n=4 quorum=3 has2=false has3=true
n=5 quorum=3 has3=true
```

- `n=2 quorum=2` and `n=4 quorum=3` prove an even cluster needs one vote past
  half, so `has1=false` for n=2 and `has2=false` for n=4.
- A `floor(n/2)` quorum lets a 2-node cluster win on 1 and a 4-node cluster
  win on 2; treating those ties as wins fails the contract.

## Readings

- **Reading ladder** — election by majority is a quorum, and you already trained that in
  M10-ex04; here it is lifted onto term-and-log.
- Raft paper, §5.2 "Leader election" (quorum, vote granting): https://raft.github.io/raft.pdf
- Reread M10-ex04 (majority quorum).
