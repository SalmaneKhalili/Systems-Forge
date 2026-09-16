package main

import "fmt"

func status(p *Peer) string {
	if p.Failed() {
		return "failed"
	}
	return "alive"
}

func main() {
	p := NewPeer(3)
	fmt.Println("timeout=3")
	p.Tick()
	fmt.Printf("tick: %s (%d left)\n", status(p), p.Left())
	p.Tick()
	fmt.Printf("tick: %s (%d left)\n", status(p), p.Left())
	p.Beat()
	fmt.Println("beat")
	p.Tick()
	fmt.Printf("tick: %s (%d left)\n", status(p), p.Left())
	p.Tick()
	fmt.Printf("tick: %s (%d left)\n", status(p), p.Left())
	p.Tick()
	fmt.Printf("tick: %s (%d left)\n", status(p), p.Left())
	p.Tick()
	fmt.Printf("tick: %s (%d left)\n", status(p), p.Left())
	p.Beat()
	fmt.Printf("beat: %s\n", status(p))
}