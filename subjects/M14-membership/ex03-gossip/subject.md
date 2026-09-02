# M14-ex03 · Gossip

## Goal

No single node knows the whole cluster at first. **Gossip** spreads membership
facts node-to-node: each node carries a view, and when two nodes exchange they
**merge** their views. Every member keeps the entry with the higher version, so
fresher facts win and all views converge to the same set.

Implement `gossip.go`:

```go
type Member struct {
	Node    string
	Version int
}

// Merge folds b into a: for every member seen in either view, keep the entry
// with the higher version. The result has one entry per member, sorted by node
// name.
func Merge(a, b []Member) []Member
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `gossip.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `Merge` must be commutative (result independent of argument order) and
  deterministic; result sorted by node name.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
a: x1 y1
b: x2 y1 z1
merge(a,b): x2 y1 z1
merge(b,a): x2 y1 z1
```

Node `x` is at version 1 in `a` and version 2 in `b`; merging — in either
order — keeps the fresher version 2. Both nodes converge on `x2 y1 z1`. Keeping
two entries for one member, or letting a stale version overwrite a fresher one
(non-commutative merge), is the bug.