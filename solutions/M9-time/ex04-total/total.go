package main

import "sort"

// Event is one timestamped event: lamport stamp and process id.
type Event struct {
	Lamport int
	Pid     int
}

// TotalOrder returns the indices of es sorted by (lamport, pid) so the
// partial happens-before order is extended to a total order.
func TotalOrder(es []Event) []int {
	idx := make([]int, len(es))
	for i := range es {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool {
		A, B := es[idx[a]], es[idx[b]]
		if A.Lamport != B.Lamport {
			return A.Lamport < B.Lamport
		}
		return A.Pid < B.Pid
	})
	return idx
}