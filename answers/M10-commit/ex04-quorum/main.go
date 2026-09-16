package main

import "fmt"

func main() {
	fmt.Printf("n=1 quorum %d  has 1? %v\n", QuorumSize(1), HasQuorum(1, 1))
	fmt.Printf("n=2 quorum %d  has 1? %v  has 2? %v\n", QuorumSize(2), HasQuorum(2, 1), HasQuorum(2, 2))
	fmt.Printf("n=3 quorum %d  has 1? %v  has 2? %v\n", QuorumSize(3), HasQuorum(3, 1), HasQuorum(3, 2))
	fmt.Printf("n=4 quorum %d  has 2? %v  has 3? %v\n", QuorumSize(4), HasQuorum(4, 2), HasQuorum(4, 3))
	fmt.Printf("n=5 quorum %d  has 2? %v  has 3? %v\n", QuorumSize(5), HasQuorum(5, 2), HasQuorum(5, 3))
}