package main

// State holds the commit index and log length.
type State struct {
	Commit int // highest committed index, or -1
	Len    int // number of entries
}

// NextIndex returns the index the next Append writes at.
func NextIndex(s State) int {
	return s.Len
}

// Uncommitted returns how many entries are beyond the commit index.
func Uncommitted(s State) int {
	return s.Len - s.Commit - 1
}