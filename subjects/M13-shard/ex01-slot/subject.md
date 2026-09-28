# M13-ex01 · Slot

## Goal

The simplest shard placement: a fixed set of slots, and each key is pinned to
one slot by hashing it. The slot owning a key must be **stable** — the same
key always maps to the same slot for a given slot count.

Implement `slot.go`:

```go
// Hash returns the FNV-1a 32-bit hash of key (offset 2166136261, prime
// 16777619), masked to 32 bits.
func Hash(key string) uint32

// SlotOf assigns key to a slot in [0, nSlots).
func SlotOf(key string, nSlots int) int
```

`SlotOf` = `int(Hash(key) % uint32(nSlots))` and must be deterministic across
calls.

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `slot.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- FNV-1a is fully specified above — reproduce it exactly so the slots match.
- Never read the wall clock.

## Acceptance

Reference transcript:

```
apple @4 -> 3
banana @4 -> 0
cherry @4 -> 0
date @4 -> 1
elder @4 -> 3
```

The tell: `banana` and `cherry` collide on slot 0, `apple` and `elder` on slot
3 — two keys hashing to the same slot always land on the same shard.

## Readings

- **Reading ladder** — hashing is the simplest *stable* key-to-slot map; DDIA frames
  partitioning generally.
- *Designing Data-Intensive Applications*, Chapter 6 "Partitioning" — hash partitioning.
- Wikipedia, "Hash function": https://en.wikipedia.org/wiki/Hash_function
