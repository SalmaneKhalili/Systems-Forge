# M15-ex04 · Drill

## Goal

The **chaos drill** proves a cluster stays correct while nodes go down. Each op
targets one node; it succeeds only if that node is **up at the current tick**.
A down node returns `fail`, and the operation must not be counted as applied.
Availability is a pure, deterministic function of the node's state — never the
wall clock.

Implement `drill.go`:

```go
type Drill struct{ up map[string]bool; tick int }

func NewDrill(up []string) *Drill
func (d *Drill) Tick()                // advance the logical clock one step
func (d *Drill) Up(node string)       // mark a node up (recovery)
func (d *Drill) Down(node string)     // mark a node down (failure)
func (d *Drill) Op(node string) string // "ok" if node up at current tick, else "fail"
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `drill.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- No wall clock — tick is the only time model.
- An op against a down node must return `fail` (never `ok`).

## Acceptance

Reference transcript:

```
drill a b c
down b
ok a
fail b
ok c
up b
ok b
```

With `b` down at the current tick, `Op(b)` is `fail` while `a`/`c` are `ok`; a
recovered `b` returns `ok` again. An op on a down node returning `ok` is the bug.

## Readings

- **Reading ladder** — a chaos drill makes node-churn *exercised* rather than theoretical.
- Principles of Chaos Engineering: https://principlesofchaos.org/
- *Designing Data-Intensive Applications*, Chapter 8 "The Trouble with Distributed Systems".
