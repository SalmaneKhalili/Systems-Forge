package main

import "fmt"

func main() {
	root := NewSpan("request", 30)
	db1 := root.Child("db.query", 40)
	db1.Attr("db", "mysql")
	db1.Child("db.scan", 25)
	cache := root.Child("cache.get", 15)
	cache.Attr("hit", "true")
	db2 := root.Child("db.query", 50)
	db2.Attr("db", "postgres")

	for _, line := range root.Dump() {
		fmt.Println(line)
	}
	fmt.Printf("count db.query=%d\n", len(root.FindAll("db.query")))
	fmt.Printf("slowest=%d\n", root.MaxDur())
}