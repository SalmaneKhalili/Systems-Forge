package main

// LogInfo describes the last entry of a log.
type LogInfo struct {
	Term int
	Idx  int
}

// Voter is a follower deciding whether to grant a RequestVote.
type Voter struct {
	term  int
	voted bool
}

// NewVoter starts at term 0 with no vote cast.
func NewVoter() *Voter {
	return &Voter{term: 0, voted: false}
}

// upToDate reports whether cand is at least as up-to-date as mine: higher
// last-term wins; on equal term, higher-or-equal last-index wins.
func upToDate(cand, mine LogInfo) bool {
	if cand.Term != mine.Term {
		return cand.Term > mine.Term
	}
	return cand.Idx >= mine.Idx
}

// RequestVote handles a candidate RequestVote. It advances the term if the
// candidate's term is newer, then grants only if the term is not stale, we
// have not already voted, and the candidate's log is at least as up-to-date.
func (v *Voter) RequestVote(candTerm int, cand, mine LogInfo) bool {
	if candTerm > v.term {
		v.term = candTerm
		v.voted = false
	}
	if candTerm < v.term {
		return false
	}
	if v.voted {
		return false
	}
	if !upToDate(cand, mine) {
		return false
	}
	v.voted = true
	return true
}