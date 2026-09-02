package main

import "fmt"

func main() {
	n := NewNode()
	fmt.Printf("t0 term=%d voted=%v\n", n.Term(), n.Voted())
	fmt.Printf("see 2 -> term=%d voted=%v\n", n.SeeTerm(2), n.Voted())
	fmt.Printf("see 5 -> term=%d voted=%v\n", n.SeeTerm(5), n.Voted())
	fmt.Printf("elect -> term=%d voted=%v\n", n.StartElection(), n.Voted())
	fmt.Printf("see 3 -> term=%d voted=%v\n", n.SeeTerm(3), n.Voted())
	fmt.Printf("t0 term=%d voted=%v\n", n.Term(), n.Voted())
}