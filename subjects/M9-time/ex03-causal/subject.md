# M9-ex03 · Causality

## Goal

Vector stamps let a process decide, for any two events, which of the three
causal relations holds:

- **happened before** — `vc(a) <= vc(b)` componentwise and not equal;
- **happened after** — the mirror image;
- **concurrent** — neither direction holds (the stamps are incomparable).

Implement `causal.go`:

```go
// CausallyBefore reports that a happened before b (strict componentwise <=).
func CausallyBefore(a, b [3]int) bool

// Concurrent reports that neither a happened before b nor b before a.
func Concurrent(a, b [3]int) bool
```

The provided `main.go` holds four events from a fixed script and prints one
boolean per relation. `make all` must build `test`; `./test` must print the
reference transcript exactly.

## Constraints

- Go, standard library only; file is `causal.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Never read the wall clock.

## Acceptance

Reference transcript:

```
e0 -> e1: true
e0 -> e2: false
e1 -> e2: false
e2 -> e1: false
e0 -> e3: false
e0 || e1: false
e0 || e2: true
e1 || e2: true
e0 || e3: true
```

`e3` is a copy of `e0`'s stamp: equality is **not** happened-before. And `e2`
outscores others on its own component while trailing elsewhere — the only
correct reads for `e2` are `||`.

A scalar counter cannot express any of this: it has no "both directions".

## Readings

- **Reading ladder** — causality is a *relation between events*, read out of the stamps; the
  notes define the three outcomes precisely.
- Martin Kleppmann, *Distributed Systems* notes — comparing vector stamps (happened-before /
  concurrent):
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
- Lamport's original ordering definitions, §1–§2:
  https://lamport.azurewebsites.net/pubs/time-clocks.pdf
