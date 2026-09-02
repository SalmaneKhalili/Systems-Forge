package main

import "fmt"

func main() {
	fmt.Printf("next(c=-1,l=0) = %d      uncommitted = %d\n", NextIndex(State{-1, 0}), Uncommitted(State{-1, 0}))
	fmt.Printf("next(c=0,l=3)  = %d      uncommitted = %d\n", NextIndex(State{0, 3}), Uncommitted(State{0, 3}))
	fmt.Printf("next(c=2,l=3)  = %d      uncommitted = %d\n", NextIndex(State{2, 3}), Uncommitted(State{2, 3}))
	fmt.Printf("next(c=1,l=5)  = %d      uncommitted = %d\n", NextIndex(State{1, 5}), Uncommitted(State{1, 5}))
}