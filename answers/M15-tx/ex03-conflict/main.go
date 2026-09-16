package main

import "fmt"

func main() {
	a := map[string]string{"x": "1"}
	b := map[string]string{"x": "2"}
	fmt.Println("conflicts ->", Conflict(a, b))

	mine := map[string]string{"x": "1"}
	theirs := map[string]string{"x": "b", "y": "c"}
	merged := Resolve(mine, theirs)
	for _, k := range []string{"x", "y"} {
		fmt.Printf("%s -> %s\n", k, merged[k])
	}
}
