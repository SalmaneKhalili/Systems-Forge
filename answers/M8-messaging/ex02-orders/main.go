package main

import (
	"fmt"
	"strings"
)

type entry struct{ partition, offset int }

// script fixes the arrival order and is deliberately tangled: partition 0's
// offset 0 lands after its offset 1, so only a gap-holding queue can produce
// the single allowed delivery order.
var script = []entry{
	{0, 1}, // a1 — held: a0 has not arrived yet
	{1, 0}, // b0 — deliverable
	{0, 0}, // a0 — now a1 may go too
	{1, 1}, //
	{0, 2}, //
	{1, 2}, //
}

func names(es []entry) string {
	parts := make([]string, len(es))
	for i, e := range es {
		parts[i] = fmt.Sprintf("%s%d", string(rune('a'+e.partition)), e.offset)
	}
	return strings.Join(parts, "|")
}

func main() {
	q := NewQueue()
	for _, e := range script {
		q.Put(e.partition, e.offset)
	}
	out, err := q.Drain()
	if err != nil {
		fmt.Printf("drain: %v\n", err)
		return
	}
	fmt.Printf("delivered: %s\n", names(out))
}