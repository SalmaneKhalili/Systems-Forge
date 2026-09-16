package main

// Lamport is a logical clock that orders events causally.
type Lamport struct {
	c int
}

// Tick counts a local event (including a send) and returns the new value.
func (l *Lamport) Tick() int {
	l.c++
	return l.c
}

// Now returns the current value without advancing.
func (l *Lamport) Now() int {
	return l.c
}

// Add reflects a received event carrying the stamp other: the clock jumps
// past it, then counts the receive itself.
func (l *Lamport) Add(other int) int {
	if other > l.c {
		l.c = other
	}
	l.c++
	return l.c
}