# M14-ex03 · Gossip

Heartbeats and the suspect lifecycle define one node's view, but membership must
spread through the cluster. Gossip exchanges versioned views pairwise so fresher
facts replace stale ones until every view converges. You deliver `gossip.go` with
a commutative, deterministic merge that gives each member one final entry.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`gossip.go`** using only the Go standard library and implement:

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

`Merge` must be commutative: its result cannot depend on argument order. Keep the
higher version for every member, return one entry per node sorted by node name,
and never read the wall clock. The reference transcript is `expected.txt` with
whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
a: x1 y1
b: x2 y1 z1
merge(a,b): x2 y1 z1
merge(b,a): x2 y1 z1
```

Node `x` is at version 1 in `a` and version 2 in `b`; merging in either order
keeps the fresher version 2, and both views converge on `x2 y1 z1`. Keeping two
entries for one member, or letting a stale version overwrite a fresher one,
fails the merge.

## Readings

- **Reading ladder** — gossip is pairwise membership exchange converging to cluster-wide
  agreement; "infection-style" is SWIM's word for it.
- SWIM paper, infection-style exchange: https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Wikipedia, "Gossip protocol": https://en.wikipedia.org/wiki/Gossip_protocol
