package main

import "fmt"

func main() {
	for _, k := range []string{"apple", "banana", "cherry", "date", "elder"} {
		fmt.Printf("%s @4 -> %d\n", k, SlotOf(k, 4))
	}
}