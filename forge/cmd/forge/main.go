// forge is the systems-forge platform: a TUI and CLI over a self-graded
// systems-engineering curriculum.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "forge:", err)
		os.Exit(1)
	}
}
