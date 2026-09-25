# M10-ex03 · Leader election

Agreement still needs a single process to drive it, so M10 next resolves a
leader race without consulting a wall clock. The bully algorithm starts from
volunteers and the highest node id; competing logical election stamps decide
which volunteer is newest, and id alone resolves a tie.

## Shape

The **bully algorithm** lets the highest-id volunteer win, but two rounds can
produce candidates with different logical timestamps. The candidate with the
**newest election** (highest timestamp) wins; a tie goes to the **higher id**.
You write **`elect.go`** and implement:

```go
type Candidate struct {
	ID   int
	Stamp int // logical election timestamp
}

// Elect returns the winning candidate, or nil when none ran.
func Elect(cs []Candidate) *Candidate
```

The exercise ships `main.go`, which runs four elections over fixed candidate
sets and prints each winner. `make all` must build `test`, and `./test` must
print the reference transcript exactly. Use Go and the standard library only.
The reference transcript is `expected.txt`, with whitespace normalized. Never
read the wall clock: the timestamp is a logical round, not a time. The result
must not depend on input-slice order.

## Acceptance

The reference transcript is:

```text
e0: leader 3
e1: leader 5
e2: leader 4
e3: leader 2
```

- `e1` races id 3 at a newer stamp against id 5 at an older stamp, so the newer
  stamp selects leader 3, not 5; the higher id would win only on a tie.
- `e2` ties two candidates at the same stamp, so the higher id wins.

## Readings

- **Reading ladder** — the bully algorithm is the classical reference; DDIA frames leader
  election against the split-brain hazard.
- Wikipedia, "Bully algorithm": https://en.wikipedia.org/wiki/Bully_algorithm
- *Designing Data-Intensive Applications*, Chapter 9 — "Leader election and the split brain".
