package main

import "sort"

// Member is one entry in the membership view.
type Member struct {
	Node string
	Seq  int
}

// EvictStale returns the members still live at logical time now: those with
// now-Seq <= staleAfter, sorted by node name.
func EvictStale(members []Member, now, staleAfter int) []Member {
	out := []Member{}
	for _, m := range members {
		if now-m.Seq <= staleAfter {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Node < out[j].Node })
	return out
}