# M13-ex03 · Ring

Range ownership partitions keys into fixed stretches; the ring generalizes that
lookup to a changing node set. A key is owned by the first node at or after its
hash position, wrapping to node 0 when it passes the last node, so most keys stay
in place across membership changes. You deliver `ring.go` with that clockwise
ownership rule.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`ring.go`** using only the Go standard library and implement:

```go
// NodeFor returns the index into nodes (0-based) of the node that owns
// keyHash. nodes is a sorted ascending list of node hash positions. The owner
// is the first node at or after keyHash, wrapping to node 0.
func NodeFor(keyHash int, nodes []int) int
```

Exact hits are inclusive: a key at the same position as a node belongs to it.
Keys past the last node wrap around to node 0. Never read the wall clock. The
reference transcript is `expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
nodes 100 300 500
key 50   -> node 0
key 100  -> node 0
key 150  -> node 1
key 300  -> node 1
key 350  -> node 2
key 600  -> node 0
```

`key 50` and `key 600` wrap to node 0. `key 100` and `key 300` land exactly on
their nodes and prove the inclusive boundary; a lookup that uses `>` instead of
`>=`, or that fails to wrap, fails this transcript.

## Readings

- **Reading ladder** — consistent hashing was invented to make resizes cheap; the original
  paper is short, Wikipedia gives the mechanics.
- Karger et al., "Consistent Hashing and Random Trees" (1997):
  https://www.cs.princeton.edu/courses/archive/fall09/cos518/papers/chash.pdf
- Wikipedia, "Consistent hashing": https://en.wikipedia.org/wiki/Consistent_hashing
