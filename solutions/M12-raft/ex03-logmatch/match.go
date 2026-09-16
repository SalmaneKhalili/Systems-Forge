package main

// Entry is one log record.
type Entry struct {
	Term int
	Cmd  string
}

// AppendEntries appends ent after the follower entry that matches prevTerm at
// prevIndex. On a match it returns the extended follower log as a copy; on a
// mismatch it returns nil (the leader must back up to a lower index).
func AppendEntries(log []Entry, prevIndex, prevTerm int, ent Entry) []Entry {
	if prevIndex < 0 || prevIndex >= len(log) || log[prevIndex].Term != prevTerm {
		return nil
	}
	out := append([]Entry{}, log...)
	return append(out, ent)
}