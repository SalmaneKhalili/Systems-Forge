# M13-ex04 · Rebalance

The ring's value appears when membership changes: only keys whose owning node
changes need to move. This exercise turns that claim into a deterministic
comparison of ownership before and after additions and removals, the move count
the shard gateway's final placement story depends on. You deliver
`rebalance.go` with owner lookup and movement counting.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`rebalance.go`** using only the Go standard library and implement:

```go
// OwnerHash returns the hash position of the node owning key on a ring built
// from nodes (sorted ascending). The owner is the first node at or after key,
// wrapping to the first node.
func OwnerHash(key int, nodes []int) int

// Moved returns the number of keys whose owning node hash differs between the
// old and new node sets.
func Moved(oldNodes, newNodes []int, keys []int) int
```

`oldNodes` and `newNodes` are each sorted ascending. Never read the wall clock.
The reference transcript is `expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
owners [100 300 500]: 50:100 150:300 250:300 350:500 450:500 550:100
add 200 -> 1 moved
add 700 -> 1 moved
remove 300 -> 2 moved
```

Adding a node at 200, between 100 and 300, takes over only key `150`. Adding
700, past 500, claims only `550` from the wrap-around owner. Removing 300
redistributes `150` and `250` to 500. Shuffling every key, or counting by
position instead of by node hash, fails the move counts.

## Readings

- **Reading ladder** — the value of consistent hashing appears exactly here: when a node
  changes, only the keys whose owner changed move.
- Wikipedia, "Consistent hashing" (why few keys move): https://en.wikipedia.org/wiki/Consistent_hashing
- *Designing Data-Intensive Applications*, Chapter 6 — "Rebalancing partitions".
