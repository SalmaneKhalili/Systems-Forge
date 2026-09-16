package main

import "fmt"

func main() {
	root := NewSpan("request", 120)
	db := root.Child("db.query", 40)
	db.Child("db.scan", 25)
	root.Child("cache.get", 15)
	root.Child("api.out", 65)

	for _, line := range root.Dump() {
		fmt.Println(line)
	}
	fmt.Printf("root.dur=%d\n", root.TotalMs())
}