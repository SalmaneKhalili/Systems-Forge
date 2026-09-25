# M14-ex01 · Heartbeat

Failure detection begins with the signal the later lifecycle, gossip, and eviction
exercises consume: a peer stays alive while heartbeats arrive. Here the harness,
not the wall clock, advances time by `Tick`, and a late heartbeat can revive a
failed peer. You deliver `peer.go` with the countdown state and its transition.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`peer.go`** using only the Go standard library and implement:

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

A `Beat` always resets the countdown to the full timeout, even if the peer just
failed: a late heartbeat revives it. Never read the wall clock; `Tick` and `Beat`
are the only time model. The reference transcript is `expected.txt` with
whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
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

With a timeout of 3, each `Tick` counts the peer down and a `Beat` resets it to
3. Three consecutive unbeat ticks fail the peer; a late `beat` then revives it.
Failing one tick early breaks the inclusive countdown, and refusing to revive on
a beat fails the final transition.

## Readings

- **Reading ladder** — SWIM is the canonical membership protocol behind heartbeats; Cassandra's
  failure detector shows the practical thresholding.
- "SWIM: Scalable Weakly-consistent Infection-style process group membership" (Das, Gupta,
  Motivala): https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Datastax, "Failure Detector":
  https://docs.datastax.com/en/archived/cassandra/3.0/cassandra/operations/opsFailureDetector.html
