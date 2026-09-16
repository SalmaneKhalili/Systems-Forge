package main

import "fmt"

// Log keeps the uncommitted suffix plus a base offset so public indexes stay
// stable across snapshots.
type Log struct {
	entries []string
	base    int
}

// Append adds cmd at the end and returns its public index.
func (l *Log) Append(cmd string) int {
	idx := l.base + len(l.entries)
	l.entries = append(l.entries, cmd)
	return idx
}

// Snapshot drops everything with index <= upto and returns how many entries
// were removed.
func (l *Log) Snapshot(upto int) int {
	keep := upto - l.base + 1
	if keep > len(l.entries) {
		keep = len(l.entries)
	}
	if keep < 0 {
		keep = 0
	}
	removed := keep
	l.entries = l.entries[keep:]
	l.base += keep
	return removed
}

// All returns the available entries as "index:value" strings.
func (l *Log) All() []string {
	out := make([]string, len(l.entries))
	for i, e := range l.entries {
		out[i] = fmt.Sprintf("%d:%s", l.base+i, e)
	}
	return out
}