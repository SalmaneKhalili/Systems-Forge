package main

import "fmt"

// Four elections over fixed candidate sets.
func main() {
	e0 := []Candidate{{3, 5}, {1, 2}} // node 3 newer stamp
	e1 := []Candidate{{3, 5}, {5, 4}} // id5 at OLDER stamp
	e2 := []Candidate{{4, 3}, {2, 3}} // tie on stamp -> higher id
	e3 := []Candidate{{2, 9}, {1, 9}, {0, 1}}
	fmt.Printf("e0: leader %d\n", Elect(e0).ID)
	fmt.Printf("e1: leader %d\n", Elect(e1).ID)
	fmt.Printf("e2: leader %d\n", Elect(e2).ID)
	fmt.Printf("e3: leader %d\n", Elect(e3).ID)
}