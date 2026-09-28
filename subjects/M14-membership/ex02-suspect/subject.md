# M14-ex02 · Suspect

## Goal

Losing a peer instantly would be too eager — a slow reply is not a dead node.
Membership uses a **suspect** state: after the first threshold a peer is
`suspect`, and only after a second, longer threshold does it become `failed`.
A heartbeat that arrives while `suspect` clears the suspicion entirely.

Implement `lifecycle.go`:

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

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `lifecycle.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `failAfter > suspectAfter`; a tick at exactly `suspectAfter` yields
  `suspect`, at exactly `failAfter` yields `failed`.
- Never read the wall clock — ticks are the only time model.

## Acceptance

Reference transcript:

```
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

With `suspectAfter=2` and `failAfter=5`, the peer turns suspect at tick 2 and
failed at tick 5. A `beat` before failure clears suspicion and restarts the
count. A lifecycle that skips `suspect` (goes straight to `failed`) or that
fails before `failAfter` ticks is the bug.

## Readings

- **Reading ladder** — suspicion is a deliberate third lifecycle stage between alive and
  failed; the SWIM paper motivates why losing a peer immediately is too eager.
- SWIM paper, suspicion mechanism: https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Cassandra's practical lifecycle (alive → suspect → dead).
