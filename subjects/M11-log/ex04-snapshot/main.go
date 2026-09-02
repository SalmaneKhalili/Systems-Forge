package main

import "fmt"

func main() {
	l := &Log{}
	fmt.Printf("append a -> %d   append b -> %d   append c -> %d   append d -> %d\n",
		l.Append("a"), l.Append("b"), l.Append("c"), l.Append("d"))
	fmt.Printf("all=%s\n", join(l.All()))
	fmt.Printf("snapshot(1) -> removed %d\n", l.Snapshot(1))
	fmt.Printf("all=%s\n", join(l.All()))
	fmt.Printf("append e -> %d\n", l.Append("e"))
	fmt.Printf("all=%s\n", join(l.All()))
}

func join(s []string) string {
	out := ""
	for i, e := range s {
		if i > 0 {
			out += "|"
		}
		out += e
	}
	return out
}