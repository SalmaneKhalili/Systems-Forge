# M9-ex02 · Vector clock

A scalar Lamport counter can only say *"this event came later"*; it cannot test
causality in both directions because unrelated events may share a stamp. This
exercise replaces that scalar with one counter per process, so ex03 can compare
complete histories rather than infer causality from one number.

## Shape

The exercise ships `main.go` with a fixed message script. You write
**`vector.go`** in the same package and implement a two-entry vector clock:

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

`make all` must build `test`, and `./test` must print the reference transcript
exactly. Use Go and the standard library only. The reference transcript is
`expected.txt`, with whitespace normalized; never read the wall clock.

## Acceptance

The reference transcript is:

```text
e0 (1,0)
p0 send m0 (2,0)
p1 recv m0 (2,1)
e1 (2,2)
p1 send m1 (2,3)
p0 recv m1 (3,3)
e2 (4,3)
```

- `p1 recv m0` proves the merge happened before the receive advanced p1's own
  component: the stamp is `(2,1)`, not `(2,0)`.
- Forgetting to merge, or advancing the wrong component, changes the vector and
  fails the transcript.

## Readings

- **Reading ladder** — you built the scalar clock; a vector clock is the same idea with one
  counter per process. The notes compare both and give the merge rule.
- Martin Kleppmann, *Distributed Systems* notes — vector clocks:
  https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
- Wikipedia, "Vector clock": https://en.wikipedia.org/wiki/Vector_clock
- Tanenbaum & Van Steen, *Distributed Systems* — vector (Ward) clocks chapter.
