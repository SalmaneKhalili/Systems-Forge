package main

import "fmt"

func main() {
	fmt.Printf("n=1 quorum=%d has1=%v\n", Quorum(1), HasMajority(1, 1))
	fmt.Printf("n=2 quorum=%d has1=%v has2=%v\n", Quorum(2), HasMajority(2, 1), HasMajority(2, 2))
	fmt.Printf("n=3 quorum=%d has2=%v\n", Quorum(3), HasMajority(3, 2))
	fmt.Printf("n=4 quorum=%d has2=%v has3=%v\n", Quorum(4), HasMajority(4, 2), HasMajority(4, 3))
	fmt.Printf("n=5 quorum=%d has3=%v\n", Quorum(5), HasMajority(5, 3))
}