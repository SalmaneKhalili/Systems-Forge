package main

import "strings"

// Log is an append-only write-ahead log. Replaying it recovers the state,
// with later writes to the same key winning (last-writer-wins).
type Log struct {
	entries []string
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Append(op string) {
	l.entries = append(l.entries, op)
}

func (l *Log) Entries() []string {
	out := make([]string, len(l.entries))
	copy(out, l.entries)
	return out
}

// Replay applies every logged op in order and returns the recovered key-value
// state. Each op is "set k=v"; a later op on the same key overwrites.
func (l *Log) Replay() map[string]string {
	state := make(map[string]string)
	for _, op := range l.entries {
		if !strings.HasPrefix(op, "set ") {
			continue
		}
		body := strings.TrimPrefix(op, "set ")
		i := strings.Index(body, "=")
		if i < 0 {
			continue
		}
		state[body[:i]] = body[i+1:]
	}
	return state
}
