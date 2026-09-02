package main

import "fmt"

func stamp(v [2]int) string { return fmt.Sprintf("(%d,%d)", v[0], v[1]) }

func main() {
	p0 := &VC{}
	p1 := &VC{}

	// p0: local event
	p0.Local(0)
	fmt.Printf("e0 %s\n", stamp(p0.Now()))

	// p0: send m0 (a send is a local event of p0)
	s0 := p0.Local(0)
	fmt.Printf("p0 send m0 %s\n", stamp(s0))

	// p1: receive m0 — fold the vector in, then the receive ticks p1
	p1.Merge(s0)
	pr := p1.Local(1)
	fmt.Printf("p1 recv m0 %s\n", stamp(pr))

	// p1: local event after the receive
	p1.Local(1)
	fmt.Printf("e1 %s\n", stamp(p1.Now()))

	// p1: send m1 carrying p0's stale knowledge
	s1 := p1.Local(1)
	fmt.Printf("p1 send m1 %s\n", stamp(s1))

	// p0: receive m1 — must fold p1's component, then tick p0
	p0.Merge(s1)
	p0.Local(0)
	fmt.Printf("p0 recv m1 %s\n", stamp(p0.Now()))

	// p0: local event afterwards
	p0.Local(0)
	fmt.Printf("e2 %s\n", stamp(p0.Now()))
}