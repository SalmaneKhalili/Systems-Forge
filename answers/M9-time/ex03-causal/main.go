package main

import "fmt"

// Three events from a fixed script: e1 depends on e0 (via a message), e2 is
// on a third process with no messages, e3 repeats e0's exact stamp.
var events = [4][3]int{
	{3, 1, 0}, // e0
	{3, 2, 1}, // e1 — depends on e0
	{1, 0, 5}, // e2 — concurrent
	{3, 1, 0}, // e3 — identical to e0
}

func main() {
	fmt.Printf("e0 -> e1: %v\n", CausallyBefore(events[0], events[1]))
	fmt.Printf("e0 -> e2: %v\n", CausallyBefore(events[0], events[2]))
	fmt.Printf("e1 -> e2: %v\n", CausallyBefore(events[1], events[2]))
	fmt.Printf("e2 -> e1: %v\n", CausallyBefore(events[2], events[1]))
	fmt.Printf("e0 -> e3: %v\n", CausallyBefore(events[0], events[3]))
	fmt.Printf("e0 || e1: %v\n", Concurrent(events[0], events[1]))
	fmt.Printf("e0 || e2: %v\n", Concurrent(events[0], events[2]))
	fmt.Printf("e1 || e2: %v\n", Concurrent(events[1], events[2]))
	fmt.Printf("e0 || e3: %v\n", Concurrent(events[0], events[3]))
}