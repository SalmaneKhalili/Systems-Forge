package main

// Quorum returns the number of votes needed to win in a group of size n.
func Quorum(n int) int {
	return n/2 + 1
}

// HasMajority reports whether yes votes are enough to win in a group of size n.
func HasMajority(n, yes int) bool {
	return yes >= Quorum(n)
}