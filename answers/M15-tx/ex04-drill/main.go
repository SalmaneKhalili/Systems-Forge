package main

import "fmt"

func main() {
	d := NewDrill([]string{"a", "b", "c"})
	fmt.Println("drill a b c")
	d.Down("b")
	fmt.Println("down b")
	d.Tick()
	fmt.Println(d.Op("a") + " a")
	fmt.Println(d.Op("b") + " b")
	fmt.Println(d.Op("c") + " c")
	d.Tick()
	d.Up("b")
	d.Tick()
	fmt.Println("up b")
	fmt.Println(d.Op("b") + " b")
}
