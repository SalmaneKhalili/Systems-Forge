# M12-ex01 · Term

## Goal

Raft's heartbeat is a **term**, a monotonically increasing integer that
orders rounds of leadership. A node never moves its term backwards, and a node
that discovers a higher term immediately steps down to follower. Implement
`node.go`:

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

`StartElection` advances the term by 1 and casts the node's own vote; the
term never goes backwards.

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `node.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- The term must be strictly monotonic; a `SeeTerm` with a lower/equal term
  must NOT change the term and must NOT reset the vote.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
t0 term=0 vote=-1
see 2 -> term=2 vote=-1
see 5 -> term=5 vote=-1
elect -> term=6 vote=self
see 3 -> term=6 vote=self
t0 term=6 vote=self
```

Once the node is at term 6, a stale `see 3` is ignored (still term 6, vote
kept). A node that resets its term downward, or that resets its vote on a
lower/equal term, is the bug.