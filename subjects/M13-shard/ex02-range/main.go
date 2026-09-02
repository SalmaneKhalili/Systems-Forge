package main

import "fmt"

func main() {
	bounds := []string{"k", "p", "u"}
	fmt.Printf("bounds")
	for _, b := range bounds {
		fmt.Printf(" %s", b)
	}
	fmt.Println()
	for _, k := range []string{"apple", "kiwi", "k", "mango", "pear", "u", "zebra"} {
		n := RangeShard(k, bounds)
		fmt.Printf("%-6s -> %d\n", k, n)
	}
}