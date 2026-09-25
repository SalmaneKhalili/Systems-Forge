# M13-ex02 · Range

Fixed slots give each key a stable owner, but range sharding turns placement into
contiguous, scan-friendly partitions of the sorted key space. This exercise fixes
the upper-boundary convention that the rebalance work and shard gateway will use.
You deliver `range.go` and an inclusive lookup over already-sorted boundaries.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`range.go`** using only the Go standard library and implement:

```go
// RangeShard returns the shard index owning key, where boundaries is a sorted
// ascending list of shard boundary values. Shard i covers
// (boundaries[i-1], boundaries[i]] (with boundaries[-1] implied to be
// -infinity and the trailing shard extending to +infinity). A key <=
// boundaries[0] goes to shard 0.
func RangeShard(key string, boundaries []string) int
```

`boundaries` is already sorted ascending. `RangeShard` returns the smallest `i`
such that `key <= boundaries[i]`, or `len(boundaries)` if the key exceeds every
boundary; rely on plain string ordering and never read the wall clock. The
reference transcript is `expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
bounds k p u
apple  -> 0
kiwi   -> 1
k      -> 0
mango  -> 1
pear   -> 2
u      -> 2
zebra  -> 3
```

With boundaries `k p u`, `apple` is below `k` in shard 0; `kiwi` and `mango`
sit between `k` and `p` in shard 1; `pear` sits between `p` and `u` in shard 2;
and `zebra` is past `u` in trailing shard 3. The key `k` lands exactly on the
first boundary and routes to shard 0, while `u` lands on the last boundary and
routes to shard 2. A `<` instead of `<=` sends `k` to shard 1 and `u` past all
three boundaries.

## Readings

- **Reading ladder** — contiguous key ranges make scan-friendly shards; the chapter frames the
  trade-off against hotspots.
- *Designing Data-Intensive Applications*, Chapter 6 — "Partitioning by key range".
- Wikipedia, "Shard (database architecture)":
  https://en.wikipedia.org/wiki/Shard_(database_architecture)
