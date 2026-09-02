package main

import "fmt"

func main() {
	// run 0 is older, run 1 is newer. Tombstone Val == "".
	runs := [][]KV{
		{KV{Key: "a", Val: "old"}, KV{Key: "c", Val: "1"}, KV{Key: "e", Val: "3"}, KV{Key: "g", Val: "1"}},
		{KV{Key: "a", Val: "stale"}, KV{Key: "b", Val: "new"}, KV{Key: "c", Val: ""}, KV{Key: "d", Val: "3"}, KV{Key: "e", Val: ""}, KV{Key: "f", Val: "x"}},
	}
	fmt.Println("after compaction:")
	for _, kv := range Compact(runs[0], runs[1]) {
		fmt.Printf("%s=%s\n", kv.Key, kv.Val)
	}
}