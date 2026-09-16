package main

// Candidate is a node volunteering in an election round.
type Candidate struct {
	ID    int
	Stamp int // logical election timestamp
}

// Elect returns the winning candidate, or nil when none ran. Newest stamp
// wins; ties broken by higher id. Order-independent.
func Elect(cs []Candidate) *Candidate {
	var best *Candidate
	for i := range cs {
		c := &cs[i]
		if best == nil ||
			c.Stamp > best.Stamp ||
			(c.Stamp == best.Stamp && c.ID > best.ID) {
			best = c
		}
	}
	return best
}