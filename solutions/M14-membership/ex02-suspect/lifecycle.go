package main

// Lifecycle tracks a peer through alive -> suspect -> failed by tick count.
type Lifecycle struct {
	suspectAfter int
	failAfter    int
	tick         int
	state        string
}

// NewLifecycle starts alive. suspectAfter and failAfter are tick counts.
func NewLifecycle(suspectAfter, failAfter int) *Lifecycle {
	return &Lifecycle{suspectAfter: suspectAfter, failAfter: failAfter, state: "alive"}
}

// State returns the current state string.
func (lc *Lifecycle) State() string { return lc.state }

// Tick advances one unit, crossing suspect then failed thresholds.
func (lc *Lifecycle) Tick() string {
	if lc.state == "failed" {
		return lc.state
	}
	lc.tick++
	if lc.tick >= lc.failAfter {
		lc.state = "failed"
	} else if lc.tick >= lc.suspectAfter {
		lc.state = "suspect"
	}
	return lc.state
}

// Beat clears suspicion: peer returns to alive and tick resets to 0.
func (lc *Lifecycle) Beat() {
	lc.tick = 0
	lc.state = "alive"
}