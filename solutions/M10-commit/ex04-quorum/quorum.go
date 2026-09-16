package main

// QuorumSize returns the smallest strict-majority quorum for n participants.
func QuorumSize(n int) int {
	return n/2 + 1
}

// HasQuorum reports whether yes votes reach a strict majority of n.
func HasQuorum(n, yes int) bool {
	return yes >= QuorumSize(n)
}