package main

import "fmt"

func main() {
	lc := NewLifecycle(2, 5)
	fmt.Println("suspectAfter=2 failAfter=5")
	lc.Tick()
	fmt.Printf("tick 1: %s\n", lc.State())
	lc.Tick()
	fmt.Printf("tick 2: %s\n", lc.State())
	lc.Tick()
	fmt.Printf("tick 3: %s\n", lc.State())
	lc.Beat()
	fmt.Println("beat: alive")
	lc.Tick()
	fmt.Printf("tick 1: %s\n", lc.State())
	lc.Tick()
	fmt.Printf("tick 2: %s\n", lc.State())
	lc.Tick()
	fmt.Printf("tick 3: %s\n", lc.State())
	lc.Tick()
	fmt.Printf("tick 4: %s\n", lc.State())
	lc.Tick()
	fmt.Printf("tick 5: %s\n", lc.State())
}