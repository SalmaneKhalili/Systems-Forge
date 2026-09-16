package main

// NodeFor returns the index into nodes of the node owning keyHash. nodes is a
// sorted ascending list of node hash positions; the owner is the first node at
// or after keyHash, wrapping to node 0.
func NodeFor(keyHash int, nodes []int) int {
	for i, n := range nodes {
		if keyHash <= n {
			return i
		}
	}
	return 0
}