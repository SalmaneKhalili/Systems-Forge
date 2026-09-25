# M13-ex01 · Slot

Fixed slots establish the placement primitive that ex02 and ex03 refine into key
ranges and a consistent-hash ring. Here every key is pinned by an exact FNV-1a
hash and slot count, so repeated calls always choose the same slot. You deliver
`slot.go` with the hash and lookup functions used to identify that stable owner.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`slot.go`** using only the Go standard library and implement:

```go
// Hash returns the FNV-1a 32-bit hash of key (offset 2166136261, prime
// 16777619), masked to 32 bits.
func Hash(key string) uint32

// SlotOf assigns key to a slot in [0, nSlots).
func SlotOf(key string, nSlots int) int
```

`SlotOf` must be `int(Hash(key) % uint32(nSlots))` and deterministic across
calls. Reproduce FNV-1a exactly so the slots match, and never read the wall
clock. The reference transcript is `expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
apple @4 -> 3
banana @4 -> 0
cherry @4 -> 0
date @4 -> 1
elder @4 -> 3
```

`banana` and `cherry` collide on slot 0, while `apple` and `elder` collide on
slot 3. Keys that hash to the same slot must always land on the same shard.

## Readings

- **Reading ladder** — hashing is the simplest *stable* key-to-slot map; DDIA frames
  partitioning generally.
- *Designing Data-Intensive Applications*, Chapter 6 "Partitioning" — hash partitioning.
- Wikipedia, "Hash function": https://en.wikipedia.org/wiki/Hash_function
