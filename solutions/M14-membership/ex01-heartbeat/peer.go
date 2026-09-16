package main

// Peer tracks one monitored peer's liveness by heartbeat countdown.
type Peer struct {
	timeout int
	tick    int
	failed  bool
}

// NewPeer starts a peer alive with a full timeout of grace.
func NewPeer(timeout int) *Peer {
	return &Peer{timeout: timeout, tick: timeout, failed: false}
}

// Left returns the countdown remaining.
func (p *Peer) Left() int { return p.tick }

// Failed reports whether the peer is currently failed.
func (p *Peer) Failed() bool { return p.failed }

// Tick advances one unit of dead time. At zero the peer is marked failed.
func (p *Peer) Tick() string {
	if p.tick > 0 {
		p.tick--
	}
	if p.tick == 0 {
		p.failed = true
	}
	if p.failed {
		return "failed"
	}
	return "alive"
}

// Beat resets the countdown to a full timeout and revives the peer.
func (p *Peer) Beat() {
	p.tick = p.timeout
	p.failed = false
}