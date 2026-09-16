package main

import "fmt"

func main() {
	l := &Log{}
	fmt.Printf("len=%d all=%s\n", l.Len(), join(l.All()))
	fmt.Printf("append set -> %d\n", l.Append("set"))
	fmt.Printf("append add -> %d\n", l.Append("add"))
	fmt.Printf("append del -> %d\n", l.Append("del"))
	fmt.Printf("len=%d all=%s\n", l.Len(), join(l.All()))
	fmt.Printf("at 1 = %s\n", l.All()[1])
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