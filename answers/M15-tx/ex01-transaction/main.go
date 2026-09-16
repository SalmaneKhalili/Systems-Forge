package main

import "fmt"

func main() {
	s := NewStore()
	s.Set("a", "1")
	fmt.Println("store.get a ->", s.Get("a"))

	t := s.Begin()
	t.Set("a", "2")
	fmt.Println("tx.get a ->", t.Get("a"))
	fmt.Println("store.get a ->", s.Get("a"))
	t.Rollback()
	fmt.Println("rollback done")
	fmt.Println("store.get a ->", s.Get("a"))

	t2 := s.Begin()
	t2.Set("a", "5")
	t2.Set("b", "6")
	t2.Commit()
	fmt.Println("commit done")
	fmt.Println("store.get a ->", s.Get("a"))
	fmt.Println("store.get b ->", s.Get("b"))
}
