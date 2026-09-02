package main

import "fmt"

// Six events: e1/e4 share Lamport 1 and e0/e3 share Lamport 2 — the ties are
// broken by process id so the order is fully deterministic.
var script = []Event{
	{2, 0}, // e0
	{1, 0}, // e1
	{3, 0}, // e2
	{2, 1}, // e3 — concurrent with e0 at the same lamport
	{1, 1}, // e4 — concurrent with e1 at the same lamport
	{4, 1}, // e5
}

func main() {
	order := TotalOrder(script)
	for i, idx := range order {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("e%d", idx)
	}
	fmt.Println()
}