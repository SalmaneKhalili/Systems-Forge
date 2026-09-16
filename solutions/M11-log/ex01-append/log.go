package main

// Log is an append-only sequence of commands, indexed from 0.
type Log struct {
	entries []string
}

// Append adds cmd at the end and returns its index.
func (l *Log) Append(cmd string) int {
	l.entries = append(l.entries, cmd)
	return len(l.entries) - 1
}

// Len returns how many entries are in the log.
func (l *Log) Len() int {
	return len(l.entries)
}

// All returns a copy of the entries.
func (l *Log) All() []string {
	out := make([]string, len(l.entries))
	copy(out, l.entries)
	return out
}