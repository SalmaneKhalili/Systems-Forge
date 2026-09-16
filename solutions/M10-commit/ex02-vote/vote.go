package main

// VoteFor returns a participant's vote: commit iff it can guarantee the
// write, else abort. A commit vote is a binding promise.
func VoteFor(canCommit bool) string {
	if canCommit {
		return "commit"
	}
	return "abort"
}