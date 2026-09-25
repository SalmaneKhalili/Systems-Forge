# M12-ex01 · Term

Raft begins every leadership round with a **term**, a monotonically increasing
integer that orders elections. This exercise makes a node's term unidirectional:
a higher wire term raises it and clears the vote, while a lower or equal term
changes nothing.

## Shape

A node never moves its term backwards. A node that discovers a higher term
immediately steps down to follower. You write **`node.go`** and implement:

```go
type Node struct {
	term int
	vote int // who this node voted for in the current term, -1 = none
}

// NewNode starts a follower at term 0 with no vote cast.
func NewNode() *Node

// Term returns the current term.
func (n *Node) Term() int

// SeeTerm observes a term from the wire. If it is higher than our own,
// step down to follower (term rises, vote resets to -1). Returns our term
// after the observation.
func (n *Node) SeeTerm(t int) int

// StartElection begins a new term and votes for self.
func (n *Node) StartElection() int
```

`StartElection` advances the term by 1 and casts the node's own vote; the term
never goes backwards.

The exercise ships `main.go`, which prints the transcript. `make all` must
build `test`, and `./test` must print the reference transcript exactly. Use Go
and the standard library only. The reference transcript is `expected.txt`, with
whitespace normalized. The term must be strictly monotonic: a `SeeTerm` with a
lower or equal term must neither change the term nor reset the vote. Never
read the wall clock.

## Acceptance

The reference transcript is:

```text
t0 term=0 vote=-1
see 2 -> term=2 vote=-1
see 5 -> term=5 vote=-1
elect -> term=6 vote=self
see 3 -> term=6 vote=self
t0 term=6 vote=self
```

At term 6, stale `see 3` is ignored, so both the term and vote remain. Moving
the term downward, or resetting its vote on a lower or equal term, fails the
contract.

## Readings

- **Reading ladder** — the visual guide gives the intuition, then the paper's §5.1 pins the
  term rules.
- Raft visual guide (interactive): https://raft.github.io/
- Raft paper, §5.1 "Raft basics" — terms: https://raft.github.io/raft.pdf
