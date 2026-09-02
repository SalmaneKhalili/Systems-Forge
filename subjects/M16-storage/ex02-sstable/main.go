package main

import (
	"fmt"
	"os"
)

func main() {
	path := os.TempDir() + "/sstable_demo.sst"
	items := []KV{
		{Key: "alpha", Val: "v0"},
		{Key: "amber", Val: "v9"},
		{Key: "beta", Val: "v2"},
		{Key: "delta", Val: "v5"},
		{Key: "kappa", Val: "v1"},
		{Key: "zeta", Val: "v6"},
		{Key: "zed", Val: "v0"},
	}
	if err := WriteSSTable(path, items); err != nil {
		fmt.Println("write error:", err)
		return
	}
	fmt.Println("sstable write ok")

	s, err := OpenSSTable(path)
	if err != nil {
		fmt.Println("open error:", err)
		return
	}
	for _, k := range []string{"zed", "kappa"} {
		if v, ok := s.Get(k); ok {
			fmt.Printf("get %s -> (%s)\n", k, v)
		} else {
			fmt.Printf("get %s -> missing\n", k)
		}
	}
	_, ok := s.Get("alpha")
	fmt.Printf("get alpha -> ok=%t\n", ok)
	_, ok = s.Get("amber")
	fmt.Printf("get amber -> ok=%t\n", ok)
	_, ok = s.Get("missing")
	fmt.Printf("get missing -> ok=%t\n", ok)
	fmt.Printf("blocks: %d\n", s.Blocks())

	os.Remove(path)
}
