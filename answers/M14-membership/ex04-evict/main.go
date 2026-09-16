package main

import "fmt"

func render(v []Member) string {
	s := ""
	for i, m := range v {
		if i > 0 {
			s += " "
		}
		s += fmt.Sprintf("%s%d", m.Node, m.Seq)
	}
	return s
}

func main() {
	members := []Member{{"a", 7}, {"b", 8}, {"c", 9}, {"d", 4}}
	fmt.Println("now=10 staleAfter=3")
	fmt.Printf("members %s\n", render(members))
	live := EvictStale(members, 10, 3)
	liveSet := map[string]bool{}
	for _, m := range live {
		liveSet[m.Node] = true
	}
	evicted := []Member{}
	for _, m := range members {
		if !liveSet[m.Node] {
			evicted = append(evicted, m)
		}
	}
	fmt.Printf("live: %s\n", render(live))
	fmt.Printf("evicted: %s\n", render(evicted))
}