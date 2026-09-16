package main

import "fmt"

func main() {
	s := NewSummary()
	for _, d := range []int64{10, 20, 30, 40, 50} {
		s.Add(d)
	}
	for _, kv := range s.SummaryValues() {
		fmt.Printf("%s=%s\n", kv.Key, kv.Val)
	}
}