package main

// Node is a raft follower tracking its term and vote state.
type Node struct {
	term  int
	voted bool
}

// NewNode starts a follower at term 0 with no vote cast.
func NewNode() *Node {
	return &Node{term: 0, voted: false}
}

// Term returns the current term.
func (n *Node) Term() int {
	return n.term
}

// Voted reports whether this node has already cast a vote this term.
func (n *Node) Voted() bool {
	return n.voted
}

// SeeTerm observes a term from the wire. A higher term steps the node down to
// follower (term rises, vote resets). A lower or equal term is ignored.
func (n *Node) SeeTerm(t int) int {
	if t > n.term {
		n.term = t
		n.voted = false
	}
	return n.term
}

// StartElection begins a new term and votes for self.
func (n *Node) StartElection() int {
	n.term++
	n.voted = true
	return n.term
}