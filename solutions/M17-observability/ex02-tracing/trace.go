package main

import (
	"fmt"
	"strings"
)

// Span is one unit of work in a request tree.
type Span struct {
	name     string
	durMs    int64
	children []*Span
}

// NewSpan creates a root span with the given name and duration.
func NewSpan(name string, durMs int64) *Span {
	return &Span{name: name, durMs: durMs}
}

// Child adds a child span under s and returns it.
func (s *Span) Child(name string, durMs int64) *Span {
	c := &Span{name: name, durMs: durMs}
	s.children = append(s.children, c)
	return c
}

// TotalMs returns the span's own duration, not including children.
func (s *Span) TotalMs() int64 {
	return s.durMs
}

// Dump renders the tree depth-first: indentation (2 spaces per depth),
// then name=durMsms.
func (s *Span) Dump() []string {
	return s.dump(0, nil)
}

func (s *Span) dump(depth int, acc []string) []string {
	line := strings.Repeat("  ", depth) + fmt.Sprintf("%s=%dms", s.name, s.durMs)
	acc = append(acc, line)
	for _, c := range s.children {
		acc = c.dump(depth+1, acc)
	}
	return acc
}