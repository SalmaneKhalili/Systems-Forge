package main

import "sort"

// Member is one entry in a membership view.
type Member struct {
	Node    string
	Version int
}

// Merge folds b into a: for every member seen in either view, keep the higher
// version. Result is deterministic, one entry per member, sorted by node name.
func Merge(a, b []Member) []Member {
	best := map[string]int{}
	order := []string{}
	for _, m := range append(append([]Member{}, a...), b...) {
		if _, seen := best[m.Node]; !seen {
			order = append(order, m.Node)
		}
		if m.Version > best[m.Node] {
			best[m.Node] = m.Version
		}
	}
	sort.Strings(order)
	out := make([]Member, 0, len(order))
	for _, n := range order {
		out = append(out, Member{n, best[n]})
	}
	return out
}