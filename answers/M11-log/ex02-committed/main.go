package main

import "fmt"

func main() {
	l := NewLog()
	l.Append("a")
	l.Append("b")
	l.Append("c")
	fmt.Printf("committed=%s commit=%d\n", join(l.Committed()), l.CommitIndex())
	fmt.Printf("commit 1 -> %d\n", l.Commit(1))
	fmt.Printf("committed=%s\n", join(l.Committed()))
	fmt.Printf("commit 2 -> %d\n", l.Commit(2))
	fmt.Printf("committed=%s\n", join(l.Committed()))
}

func join(s []string) string {
	if len(s) == 0 {
		return "[]"
	}
	out := ""
	for i, e := range s {
		if i > 0 {
			out += "|"
		}
		out += e
	}
	return out
}