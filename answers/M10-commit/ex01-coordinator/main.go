package main

import "fmt"

func main() {
	c := &Coordinator{n: 3}
	votes := []bool{true, true, false}
	// Phase one: prepare each participant and print its vote.
	voice := make([]bool, len(votes))
	for i := 0; i < c.n; i++ {
		voice[i] = votes[i]
		if voice[i] {
			fmt.Printf("prepare %d -> yes\n", i)
		} else {
			fmt.Printf("prepare %d -> no\n", i)
		}
	}
	// Phase two: the coordinator decides from the collected votes.
	out := c.Run(func(id int) bool { return votes[id] })
	fmt.Printf("DECISION %s\n", out)
	for i := 0; i < c.n; i++ {
		if out == "commit" {
			fmt.Printf("commit %d\n", i)
		} else {
			fmt.Printf("abort %d\n", i)
		}
	}
}