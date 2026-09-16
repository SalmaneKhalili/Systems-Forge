package main

// Voter returns a participant's vote: true = vote to commit.
type Voter func(id int) bool

// Coordinator drives two-phase commit over n participants. Run returns
// "commit" only if EVERY participant votes commit, else "abort".
type Coordinator struct {
	n int
}

// Run queries every participant once and returns the agreed decision.
func (c *Coordinator) Run(vote Voter) string {
	for i := 0; i < c.n; i++ {
		if !vote(i) {
			return "abort"
		}
	}
	return "commit"
}