package main

import (
	"fmt"
	"strings"
)

func ownersLine(nodes, keys []int) string {
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%d:%d", k, OwnerHash(k, nodes))
	}
	return strings.Join(parts, " ")
}

func main() {
	base := []int{100, 300, 500}
	keys := []int{50, 150, 250, 350, 450, 550}
	fmt.Printf("owners [100 300 500]: %s\n", ownersLine(base, keys))
	fmt.Printf("add 200 -> %d moved\n", Moved(base, []int{100, 200, 300, 500}, keys))
	fmt.Printf("add 700 -> %d moved\n", Moved(base, []int{100, 300, 500, 700}, keys))
	fmt.Printf("remove 300 -> %d moved\n", Moved(base, []int{100, 500}, keys))
}