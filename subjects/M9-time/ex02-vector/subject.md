# M9-ex02 · Vector clock

## Goal

A scalar Lamport counter can only say *"this event came later"* — it cannot
test causality in both directions, because two unrelated events may share a
stamp. A **vector clock** fixes that: each of the two processes keeps one
counter per process, advances its *own* component on local events, and folds
*received* vectors in elementwise before advancing.

Implement `vector.go` with a two-entry vector clock:

```go
type VC struct {
	v [2]int
}

// Local ticks component i for a local event and returns the new vector.
func (vc *VC) Local(i int) [2]int

// Merge folds a received vector in elementwise (max per component).
func (vc *VC) Merge(other [2]int)

// Now returns the current vector without changing it.
func (vc *VC) Now() [2]int
```

The provided `main.go` plays a fixed message script and prints the transcript.
`make all` must build `test`; `./test` must print the reference transcript
exactly.

## Constraints

- Go, standard library only; file is `vector.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Never read the wall clock.

## Acceptance

Reference transcript:

```
e0 (1,0)
p0 send m0 (2,0)
p1 recv m0 (2,1)
e1 (2,2)
p1 send m1 (2,3)
p0 recv m1 (3,3)
e2 (4,3)
```

Watch `p1 recv m0`: after merging (2,0), p1's own component advances — the
receive is an event of p1 — so it reads (2,1), not (2,0). A clock that forgets
to merge (or advances the wrong component) prints a different vector.

## Readings

- **Reading ladder** — you built the scalar clock; a vector clock is the same idea with one
  counter per process. The notes compare both and give the merge rule.
- Martin Kleppmann, *Distributed Systems* notes — vector clocks:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
- Wikipedia, "Vector clock": https://en.wikipedia.org/wiki/Vector_clock
- Tanenbaum & Van Steen, *Distributed Systems* — vector (Ward) clocks chapter.
