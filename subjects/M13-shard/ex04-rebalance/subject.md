# M13-ex04 · Rebalance

## Goal

The value of consistent hashing shows up on resize: when a node joins or
leaves, only the keys whose owning node changes actually move. Counting those
moves is the clean way to verify a sharder stays balanced with minimal churn.

Implement `rebalance.go`:

```go
// OwnerHash returns the hash position of the node owning key on a ring built
// from nodes (sorted ascending). The owner is the first node at or after key,
// wrapping to the first node.
func OwnerHash(key int, nodes []int) int

// Moved returns the number of keys whose owning node hash differs between the
// old and new node sets.
func Moved(oldNodes, newNodes []int, keys []int) int
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `rebalance.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `oldNodes`/`newNodes` are each sorted ascending.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
owners [100 300 500]: 50:100 150:300 250:300 350:500 450:500 550:100
add 200 -> 1 moved
add 700 -> 1 moved
remove 300 -> 2 moved
```

Adding a node at 200 (between 100 and 300) takes over only key `150`; adding
700 (past 500) claims only `550` from the wrap-around owner; removing 300
redistributes `150` and `250` to 500. Shuffling every key, or counting by
position instead of by node hash, is the bug.