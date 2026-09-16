# M13-ex03 · Ring

## Goal

Consistent hashing places both keys and nodes on a ring of hash positions. A
key is owned by the **next node clockwise** from the key's position — the
first node at or after the key, wrapping around to the first node. This keeps
most keys in place when the node set changes.

Implement `ring.go`:

```go
// NodeFor returns the index into nodes (0-based) of the node that owns
// keyHash. nodes is a sorted ascending list of node hash positions. The owner
// is the first node at or after keyHash, wrapping to node 0.
func NodeFor(keyHash int, nodes []int) int
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `ring.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- Exact hits are inclusive: a key at the same position as a node belongs to it.
- Keys past the last node wrap around to node 0.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
nodes 100 300 500
key 50   -> node 0
key 100  -> node 0
key 150  -> node 1
key 300  -> node 1
key 350  -> node 2
key 600  -> node 0
```

The tell: `key 50` and `key 600` both wrap to node 0; `key 100` and `key 300`
land exactly on their node (inclusive hit). A lookup that uses `>` instead of
`>=`, or that fails to wrap, is the bug.