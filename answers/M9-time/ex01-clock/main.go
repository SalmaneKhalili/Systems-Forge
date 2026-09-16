package main

import "fmt"

func main() {
	p0 := &Lamport{}
	p1 := &Lamport{}

	// p0: local event
	p0.Tick()
	fmt.Printf("e0 t%d\n", p0.Now())

	// p0: send m0 — a send is a local event, stamped before transmission
	s0 := p0.Tick()
	fmt.Printf("p0 send m0 t%d\n", s0)

	// p1: receive m0 — merge past the sender's stamp, then count the receive
	p1.Add(s0)
	fmt.Printf("p1 recv m0 t%d\n", p1.Now())

	// p1: local event after the receive
	p1.Tick()
	fmt.Printf("e1 t%d\n", p1.Now())

	// p0: local, unrelated to p1's progress
	p0.Tick()
	fmt.Printf("e2 t%d\n", p0.Now())

	// p1: send m1 carrying a stamp p0 has not seen
	s1 := p1.Tick()
	fmt.Printf("p1 send m1 t%d\n", s1)

	// p0: receive m1 — must jump past p1's stamp
	p0.Add(s1)
	fmt.Printf("p0 recv m1 t%d\n", p0.Now())

	// p0: local event afterwards
	p0.Tick()
	fmt.Printf("e3 t%d\n", p0.Now())
}