package main

import "fmt"

func main() {
	m := NewMemTable()
	m.Put("cherry", "3")
	m.Put("banana", "2")
	m.Put("apple", "1")
	m.Put("banana", "4") // overwrite

	fmt.Println("keys in order:")
	for _, kv := range m.All() {
		fmt.Printf("%s=%s\n", kv.Key, kv.Val)
	}

	val, ok := m.Get("banana")
	fmt.Printf("get banana: %s (ok=%t)\n", val, ok)
	_, ok2 := m.Get("dragon")
	fmt.Printf("get dragon: (ok=%t)\n", ok2)
}
