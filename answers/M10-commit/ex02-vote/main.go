package main

import "fmt"

func main() {
	environ := []bool{true, false, true, false} // can each transaction commit?
	for i, ok := range environ {
		fmt.Printf("t%d: %s\n", i, VoteFor(ok))
	}
}