package main

import "fmt"

func render(v []Member) string {
	s := ""
	for i, m := range v {
		if i > 0 {
			s += " "
		}
		s += fmt.Sprintf("%s%d", m.Node, m.Version)
	}
	return s
}

func main() {
	a := []Member{{"x", 1}, {"y", 1}}
	b := []Member{{"x", 2}, {"y", 1}, {"z", 1}}
	fmt.Printf("a: %s\n", render(a))
	fmt.Printf("b: %s\n", render(b))
	fmt.Printf("merge(a,b): %s\n", render(Merge(a, b)))
	fmt.Printf("merge(b,a): %s\n", render(Merge(b, a)))
}