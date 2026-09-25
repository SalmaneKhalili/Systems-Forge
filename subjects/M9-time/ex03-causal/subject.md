# M9-ex03 · Causality

The two-process vectors from ex02 become a relation in this exercise. You
compare complete stamps to classify each pair as happened-before,
happened-after or concurrent, exposing the partial order that ex04 extends
without breaking causality.

## Shape

The exercise ships `main.go` with four events from a fixed script. You write
**`causal.go`** in the same package and implement:

```go
// CausallyBefore reports that a happened before b (strict componentwise <=).
func CausallyBefore(a, b [3]int) bool

// Concurrent reports that neither a happened before b nor b before a.
func Concurrent(a, b [3]int) bool
```

A pair is:

- **happened before** when `vc(a) <= vc(b)` componentwise and the stamps are not equal;
- **happened after** in the mirror-image case;
- **concurrent** when neither direction holds and the stamps are incomparable.

`make all` must build `test`, and `./test` must print one boolean per relation
in the reference transcript exactly. Use Go and the standard library only. The
reference transcript is `expected.txt`, with whitespace normalized; never read
the wall clock.

## Acceptance

The reference transcript is:

```text
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

- `e3` copies `e0`'s stamp, so equality is **not** happened-before.
- `e2` leads on its own component while trailing elsewhere, so its only correct
  relations to the shown events are `||`.

A scalar counter cannot express these outcomes because it has no "both
directions".

## Readings

- **Reading ladder** — causality is a *relation between events*, read out of the stamps; the
  notes define the three outcomes precisely.
- Martin Kleppmann, *Distributed Systems* notes — comparing vector stamps (happened-before /
  concurrent):
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
- Lamport's original ordering definitions, §1–§2:
  https://lamport.azurewebsites.net/pubs/time-clocks.pdf
