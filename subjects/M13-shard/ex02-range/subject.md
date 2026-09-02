# M13-ex02 · Range

## Goal

Instead of hashing, key-range sharding splits the sorted key space into
contiguous ranges. Each shard owns everything between two boundaries. A key is
routed to the shard whose boundary is the greatest one `<=` the key.

Implement `range.go`:

```go
// RangeShard returns the shard index owning key, where boundaries is a sorted
// ascending list of shard boundary values. Shard i covers
// (boundaries[i-1], boundaries[i]] (with boundaries[-1] implied to be
// -infinity and the trailing shard extending to +infinity). A key <=
// boundaries[0] goes to shard 0.
func RangeShard(key string, boundaries []string) int
```

`boundaries` is sorted ascending; `RangeShard` returns the smallest `i` such
that `key <= boundaries[i]`, or `len(boundaries)` if the key exceeds every
boundary.

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `range.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- `boundaries` is already sorted ascending; rely on plain string ordering.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
bounds k p u
apple  -> 0
kiwi   -> 1
k      -> 0
mango  -> 1
pear   -> 2
u      -> 2
zebra  -> 3
```

With boundaries `k p u`: `apple` is below `k` (shard 0); `kiwi` and `mango` sit
between `k` and `p` (shard 1); `pear` between `p` and `u` (shard 2); `zebra` is
past `u` (shard 3, the trailing shard). The key `k` lands exactly ON the first
boundary → shard 0 (inclusive), and `u` on the last boundary → shard 2 — a
`<`-instead-of-`<=` comparator that excludes the boundary would route `k` to
shard 1 and `u` past all three. That is the bug to catch.