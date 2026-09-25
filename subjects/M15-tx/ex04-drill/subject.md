# M15-ex04 · Drill

The in-memory transaction and conflict rules now need a failure boundary. This
exercise drives node up/down state through a logical tick and makes every
operation's outcome a pure function of the target's state at that tick. You
deliver `drill.go`, the deterministic availability model the gateway exercises
before the durable store joins the stack.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`drill.go`** using only the Go standard library and implement:

```go
type Drill struct{ up map[string]bool; tick int }

func NewDrill(up []string) *Drill
func (d *Drill) Tick()                // advance the logical clock one step
func (d *Drill) Up(node string)       // mark a node up (recovery)
func (d *Drill) Down(node string)     // mark a node down (failure)
func (d *Drill) Op(node string) string // "ok" if node up at current tick, else "fail"
```

`Tick` is the only time model; never read the wall clock. An operation against a
down node must return `fail`, never `ok`, and therefore must not count as
applied. The reference transcript is `expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
drill a b c
down b
ok a
fail b
ok c
up b
ok b
```

With `b` down at the current tick, `Op(b)` is `fail` while `a` and `c` are
`ok`; after recovery, `b` returns `ok` again. Returning `ok` for an operation on
a down node breaks the availability rule.

## Readings

- **Reading ladder** — a chaos drill makes node-churn *exercised* rather than theoretical.
- Principles of Chaos Engineering: https://principlesofchaos.org/
- *Designing Data-Intensive Applications*, Chapter 8 "The Trouble with Distributed Systems".
