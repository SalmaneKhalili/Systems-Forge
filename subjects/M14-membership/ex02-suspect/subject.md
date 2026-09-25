# M14-ex02 · Suspect

A single missed heartbeat is too eager to declare failure, so the heartbeat
model gains a second threshold and a reversible suspicion state. This exercise
makes `alive → suspect → failed` explicit under harness-driven ticks; a beat
during suspicion returns the peer to alive and restarts its clock. You deliver
`lifecycle.go` with that state machine.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`lifecycle.go`** using only the Go standard library and implement:

```go
// Lifecycle tracks a peer through alive -> suspect -> failed.
type Lifecycle struct {
	suspectAfter int // ticks until suspect
	failAfter    int // ticks until failed
	tick         int
	state        string // "alive" | "suspect" | "failed"
}

// NewLifecycle starts alive. suspectAfter and failAfter are tick counts
// (failAfter > suspectAfter).
func NewLifecycle(suspectAfter, failAfter int) *Lifecycle

// Tick advances one unit. Crosses suspectAfter -> "suspect", failAfter ->
// "failed". Returns the state.
func (lc *Lifecycle) Tick() string

// Beat clears suspicion: the peer returns to alive, tick resets to 0.
func (lc *Lifecycle) Beat()

// State returns the current state string.
func (lc *Lifecycle) State() string
```

`failAfter` is greater than `suspectAfter`. A tick at exactly `suspectAfter`
yields `suspect`, and one at exactly `failAfter` yields `failed`. Never read the
wall clock; ticks are the only time model. The reference transcript is
`expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
suspectAfter=2 failAfter=5
tick 1: alive
tick 2: suspect
tick 3: suspect
beat: alive
tick 1: alive
tick 2: suspect
tick 3: suspect
tick 4: suspect
tick 5: failed
```

With `suspectAfter=2` and `failAfter=5`, the peer becomes suspect at tick 2 and
failed at tick 5. A beat before failure clears suspicion and restarts the count.
Skipping `suspect`, or declaring failure before `failAfter` ticks, breaks the
lifecycle contract.

## Readings

- **Reading ladder** — suspicion is a deliberate third lifecycle stage between alive and
  failed; the SWIM paper motivates why losing a peer immediately is too eager.
- SWIM paper, suspicion mechanism: https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Cassandra's practical lifecycle (alive → suspect → dead).
