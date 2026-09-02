package main

import "fmt"

func show(name string, v *Voter, ct int, cand, mine LogInfo) {
	fmt.Printf("%s: cand(%d,%d) vs mine(%d,%d) -> %v\n",
		name, cand.Term, cand.Idx, mine.Term, mine.Idx,
		v.RequestVote(ct, cand, mine))
}

func main() {
	show("a", NewVoter(), 1, LogInfo{1, 1}, LogInfo{1, 0})
	show("b", NewVoter(), 1, LogInfo{1, 0}, LogInfo{1, 1})
	show("c", NewVoter(), 2, LogInfo{2, 0}, LogInfo{1, 9})
	show("d", NewVoter(), 1, LogInfo{1, 5}, LogInfo{0, 8})
	v := NewVoter()
	show("e", v, 1, LogInfo{1, 0}, LogInfo{1, 0})
	show("f", v, 1, LogInfo{1, 0}, LogInfo{1, 0})
}