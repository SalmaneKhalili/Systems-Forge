package main

// Log is an append-only sequence with a commit index marking durability.
type Log struct {
	entries []string
	commit  int // highest committed index, or -1 when none
}

// NewLog returns an empty log with no committed entries.
func NewLog() *Log {
	return &Log{commit: -1}
}

// Append adds cmd at the end and returns its index.
func (l *Log) Append(cmd string) int {
	l.entries = append(l.entries, cmd)
	return len(l.entries) - 1
}

// Commit raises the commit index to idx and returns it.
func (l *Log) Commit(idx int) int {
	if idx > l.commit {
		l.commit = idx
	}
	return l.commit
}

// Committed returns the entries at or below the commit index.
func (l *Log) Committed() []string {
	if l.commit < 0 {
		return []string{}
	}
	if l.commit+1 > len(l.entries) {
		l.commit = len(l.entries) - 1
	}
	out := make([]string, l.commit+1)
	copy(out, l.entries)
	return out
}

// CommitIndex returns the current commit index.
func (l *Log) CommitIndex() int {
	return l.commit
}