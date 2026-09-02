package main

import (
	"fmt"
	"strings"
)

func render(log []Entry) string {
	if log == nil {
		return "REJECT"
	}
	parts := make([]string, len(log))
	for i, e := range log {
		parts[i] = fmt.Sprintf("%s%d", e.Cmd, e.Term)
	}
	return strings.Join(parts, " ")
}

func main() {
	log := []Entry{{2, "b"}, {3, "c"}, {4, "d"}}
	fmt.Printf("log   =%s\n", render(log))
	fmt.Printf("match prev(0,2) -> %s\n", render(AppendEntries(log, 0, 2, Entry{5, "e"})))
	log2 := []Entry{{2, "b"}, {3, "c"}, {4, "d"}}
	fmt.Printf("mismatch prev(0,3) -> %s\n", render(AppendEntries(log2, 0, 3, Entry{5, "e"})))
}