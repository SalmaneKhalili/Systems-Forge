# M9-ex04 · Total order

Causality from ex03 is a **partial** order: concurrent events remain
incomparable. A log, ledger or replicated state machine needs one deterministic
sequence that still respects that order, so this exercise extends every stamp
with its sender's process id and breaks Lamport ties.

## Shape

The exercise ships `main.go` with six events, including ties at Lamport 1 and
2. You write **`total.go`** in the same package and implement:

```go
type Event struct {
	Lamport int
	Pid     int
}

// TotalOrder returns the indices of es sorted by (lamport, pid) so the
// partial happens-before order is extended to a total order.
func TotalOrder(es []Event) []int
```

Compare events lexicographically by `(lamport, pid)`: the Lamport stamp comes
first, and when concurrent events tie on the same value, the smaller pid wins.
`make all` must build `test`, and `./test` must print the reference transcript
exactly.

Use Go and the standard library only. The sort must be **stable and
deterministic** and must never rely on input order for ties. The reference
transcript is `expected.txt`, with whitespace normalized; never read the wall
clock.

## Acceptance

The reference transcript is:

```text
e1 e4 e0 e3 e2 e5
```

- `e1` and `e4` both carry Lamport 1, so pid 0 (`e1`) precedes pid 1 (`e4`).
- `e0` and `e3` both carry Lamport 2, so `e0` precedes `e3`.
- Ordering by pid before Lamport is the classic bug and produces a different
  sequence.

## Readings

- **Reading ladder** — a partial order becomes a total order the moment you break ties; the
  paper's construction is the canonical one.
- Lamport, "Time, Clocks, and the Ordering of Events…", §2.3 "Total Ordering of Events":
  https://lamport.azurewebsites.net/pubs/time-clocks.pdf
- Wikipedia, "Total order": https://en.wikipedia.org/wiki/Total_order
