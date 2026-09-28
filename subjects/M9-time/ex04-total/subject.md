# M9-ex04 · Total order

## Goal

Lamport's happens-before is a **partial** order: concurrent events are
incomparable. A system that must process events in a single deterministic
sequence — a log, a ledger, a replicated state machine — needs a **total**
order that still respects causality.

The classic construction extends each stamp with the sender's process id and
orders by `(lamport, pid)` lexicographically: compare stamps first; when two
events tie on the same Lamport value (they are concurrent), the smaller pid
wins.

Implement `total.go`:

```go
type Event struct {
	Lamport int
	Pid     int
}

// TotalOrder returns the indices of es sorted by (lamport, pid) so the
// partial happens-before order is extended to a total order.
func TotalOrder(es []Event) []int
```

The provided `main.go` feeds six events (two ties at Lamport 1 and 2) and
prints the ordered event ids. `make all` must build `test`; `./test` must
print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `total.go`.
- The sort must be **stable and deterministic** — never rely on the input
  order for ties.
- Reference transcript is `expected.txt` (whitespace normalized).
- Never read the wall clock.

## Acceptance

Reference transcript:

```
e1 e4 e0 e3 e2 e5
```

`e1` and `e4` both carry Lamport 1 → pid 0 (`e1`) before pid 1 (`e4`). `e0`
and `e3` both carry Lamport 2 → `e0` before `e3`. Ordering by pid first
instead of Lamport first is the classic bug and yields a different sequence.

## Readings

- **Reading ladder** — a partial order becomes a total order the moment you break ties; the
  paper's construction is the canonical one.
- Lamport, "Time, Clocks, and the Ordering of Events…", §2.3 "Total Ordering of Events":
  https://lamport.azurewebsites.net/pubs/time-clocks.pdf
- Wikipedia, "Total order": https://en.wikipedia.org/wiki/Total_order
