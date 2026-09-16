package main

import "sort"

// Conflict tells whether two transaction write-sets collide on any key: both
// writing the same key is a potential lost update.
func Conflict(writesA, writesB map[string]string) bool {
	for k := range writesA {
		if _, ok := writesB[k]; ok {
			return true
		}
	}
	return false
}

// Resolve merges a peer's writes into yours, but when both wrote the same key
// YOUR write wins (the opponent's conflicting write is dropped as a lost
// update). Returns the merged result as a copy.
func Resolve(mine, theirs map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range mine {
		out[k] = v
	}
	for k, v := range theirs {
		if _, ok := mine[k]; !ok {
			out[k] = v
		}
	}
	return out
}

// SortedKeys returns the map's keys in sorted order (helper for tests).
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
