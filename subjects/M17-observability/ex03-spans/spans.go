package main

import (
	"fmt"
	"sort"
	"strings"
)

// Span is a named unit of work with a duration and key=value attributes.
type Span struct {
	name     string
	durMs    int64
	attr     map[string]string
	children []*Span
}

// NewSpan creates a root span.
func NewSpan(name string, durMs int64) *Span {
	return &Span{name: name, durMs: durMs, attr: map[string]string{}}
}

// Child adds a child span under s and returns it.
func (s *Span) Child(name string, durMs int64) *Span {
	c := NewSpan(name, durMs)
	s.children = append(s.children, c)
	return c
}

// Attr sets an attribute key=value on the span.
func (s *Span) Attr(k, v string) {
	s.attr[k] = v
}

// Dump renders the tree depth-first as "name{attrs}=durMs".
func (s *Span) Dump() []string {
	var out []string
	s.dump(&out, 0)
	return out
}

func (s *Span) dump(out *[]string, depth int) {
	keys := make([]string, 0, len(s.attr))
	for k := range s.attr {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+s.attr[k])
	}
	attrs := strings.Join(parts, ",")
	indent := strings.Repeat("  ", depth)
	*out = append(*out, fmt.Sprintf("%s%s{%s}=%dms", indent, s.name, attrs, s.durMs))
	for _, c := range s.children {
		c.dump(out, depth+1)
	}
}

// FindAll returns every span (in any subtree) whose name matches, in
// depth-first order.
func (s *Span) FindAll(name string) []*Span {
	var out []*Span
	s.find(name, &out)
	return out
}

func (s *Span) find(name string, out *[]*Span) {
	if s.name == name {
		*out = append(*out, s)
	}
	for _, c := range s.children {
		c.find(name, out)
	}
}

// MaxDur returns the longest own duration among s and all descendants.
func (s *Span) MaxDur() int64 {
	best := s.durMs
	var walk func(*Span)
	walk = func(sp *Span) {
		for _, c := range sp.children {
			if c.durMs > best {
				best = c.durMs
			}
			walk(c)
		}
	}
	walk(s)
	return best
}