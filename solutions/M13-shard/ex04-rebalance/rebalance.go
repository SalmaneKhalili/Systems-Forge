package main

// OwnerHash returns the node hash position owning key on a ring built from
// nodes (sorted ascending). The owner is the first node at or after key,
// wrapping to the first node.
func OwnerHash(key int, nodes []int) int {
	for _, n := range nodes {
		if key <= n {
			return n
		}
	}
	return nodes[0]
}

// Moved returns the number of keys whose owning node hash differs between the
// old and new node sets.
func Moved(oldNodes, newNodes []int, keys []int) int {
	n := 0
	for _, k := range keys {
		if OwnerHash(k, oldNodes) != OwnerHash(k, newNodes) {
			n++
		}
	}
	return n
}