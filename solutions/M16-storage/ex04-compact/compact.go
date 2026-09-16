package main

import "sort"

// KV is a key-value pair, ordered by Key.
type KV struct {
	Key, Val string
}

// Compact merges several sorted runs into a single sorted run. Each key
// appears once; among runs containing the key, the highest run index wins.
// A result value of "" is a tombstone and drops the key entirely.
func Compact(runs ...[]KV) []KV {
	// Collect the winning value per key: iterate runs newest-first and keep
	// the first (i.e. newest) occurrence.
	best := map[string]string{}
	seen := map[string]bool{}
	for r := len(runs) - 1; r >= 0; r-- {
		for _, kv := range runs[r] {
			if !seen[kv.Key] {
				seen[kv.Key] = true
				best[kv.Key] = kv.Val
			}
		}
	}
	keys := make([]string, 0, len(best))
	for k := range best {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]KV, 0, len(keys))
	for _, k := range keys {
		v := best[k]
		if v == "" { // tombstone: drop the key
			continue
		}
		out = append(out, KV{Key: k, Val: v})
	}
	return out
}