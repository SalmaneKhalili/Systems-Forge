package main

import "fmt"

func main() {
	nodes := []int{100, 300, 500}
	fmt.Printf("nodes")
	for _, n := range nodes {
		fmt.Printf(" %d", n)
	}
	fmt.Println()
	for _, k := range []int{50, 100, 150, 300, 350, 600} {
		fmt.Printf("key %-3d -> node %d\n", k, NodeFor(k, nodes))
	}
}