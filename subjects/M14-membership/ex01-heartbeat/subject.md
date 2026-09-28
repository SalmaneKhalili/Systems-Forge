# M14-ex01 · Heartbeat

## Goal

Failure detection starts with heartbeats. A peer is **alive** while heartbeats
keep arriving; if nothing arrives for a configured timeout, the peer moves to
**failed**. In this deterministic model, time is advanced by the harness as a
`Tick` count — never by reading the wall clock.

Implement `peer.go`:

```go
// Peer tracks one monitored peer's liveness by heartbeat countdown.
type Peer struct {
	timeout int // heartbeats of grace before the peer is marked failed
	tick    int // countdown remaining
	failed  bool
}

// NewPeer starts a peer alive with a full timeout of grace.
func NewPeer(timeout int) *Peer

// Beat resets the countdown to a full timeout (a new heartbeat arrived).
func (p *Peer) Beat()

// Tick advances one unit of dead time. When the countdown reaches zero the
// peer is marked failed. Returns "alive" or "failed".
func (p *Peer) Tick() string

// Failed reports whether the peer is currently failed.
func (p *Peer) Failed() bool
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `peer.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- A `Beat` always resets the countdown to the full timeout, even if the peer
  just failed (a late heartbeat revives it).
- Never read the wall clock — tick/beat are the only time model.

## Acceptance

Reference transcript:

```
timeout=3
tick: alive (2 left)
tick: alive (1 left)
beat
tick: alive (2 left)
tick: alive (1 left)
tick: failed (0 left)
tick: failed (0 left)
beat: alive
```

With a timeout of 3, each `Tick` counts the peer down; a `Beat` resets to 3.
After three consecutive unbeat ticks the peer fails, and a late `beat` revives
it. A peer that fails one tick early (off-by-one on the countdown), or that
fails and never revives on a beat, is the bug.

## Readings

- **Reading ladder** — SWIM is the canonical membership protocol behind heartbeats; Cassandra's
  failure detector shows the practical thresholding.
- "SWIM: Scalable Weakly-consistent Infection-style process group membership" (Das, Gupta,
  Motivala): https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Datastax, "Failure Detector":
  https://docs.datastax.com/en/archived/cassandra/3.0/cassandra/operations/opsFailureDetector.html
