package main

import "fmt"

func main() {
	l := NewLog()
	l.Append("set a=1")
	l.Append("set b=2")
	l.Append("set a=3")
	fmt.Println("entries ->", len(l.Entries()))
	st := l.Replay()
	fmt.Println("a ->", st["a"])
	fmt.Println("b ->", st["b"])
}
