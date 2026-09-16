# M10-ex03 · Leader election

## Goal

Replicated systems pick a single **leader**. A simple and robust scheme — the
**bully algorithm** — says: the node with the highest id among those that
volunteered wins. But rounds may race: two nodes can volunteer with
different logical timestamps. The rule that keeps a single winner:

- the candidate with the **newest election** (highest timestamp) wins;
- ties are broken by **higher id**.

Implement `elect.go`:

```go
type Candidate struct {
	ID   int
	Stamp int // logical election timestamp
}

// Elect returns the winning candidate, or nil when none ran.
func Elect(cs []Candidate) *Candidate
```

The provided `main.go` runs four elections over fixed candidate sets and prints
the winner of each. `make all` must build `test`; `./test` must print the
reference transcript exactly.

## Constraints

- Go, standard library only; file is `elect.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Never read the wall clock; the timestamp is a logical round, not a time.
- The result must not depend on the order of the input slice.

## Acceptance

Reference transcript:

```
e0: leader 3
e1: leader 5
e2: leader 4
e3: leader 2
```

`e1` is the tell: two candidates race, id 3 at a newer stamp and id 5 at an
older stamp — the *newer* stamp wins (so the leader is 3, not 5; the higher
id would win only on a tie). `e2` ties two candidates on the same stamp and
the higher id wins.