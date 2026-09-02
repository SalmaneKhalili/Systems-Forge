# M16-ex04 · Compaction

## Goal

Over time an LSM-tree accumulates many overlapping SSTables (ex02), and lookups must check them in order. **Compaction** merges a set of sorted runs into fewer, larger sorted runs. During the merge, keys that appear in multiple runs collapse to a single entry — the newest run (last in the input order) wins, so stale overwritten values are dropped.

Implement `compact.go`:

```go
// Compact merges several sorted runs (each a sorted []KV) into one sorted
// run. Each key appears at most once; when a key appears in more than one
// run, the later run (higher index) wins. A KV with Val == "" acts as a
// tombstone and removes the key entirely.
func Compact(runs ...[]KV) []KV
```

The provided `main.go` merges a small set of overlapping runs and prints the result. `make all` must build `test`; `./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `compact.go`. Define `KV` (a `Key, Val string` pair) in this package.
- Output must be strictly sorted by key, with no duplicate keys.

## Acceptance

Reference transcript:

```
after compaction:
a=stale
b=new
d=3
f=tombstoned
g=1
```

In the sample, `a` appears in both runs (the newer `stale` wins), `c` and `e` were written then tombstoned (absent), `b`, `d`, `f` survive from the newer run, and `g` survives from the older run. A compaction that leaves duplicate keys, drops a surviving key, or fails to drop a tombstoned key, is the bug.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3: "SSTables and LSM-Trees" (Compaction and performance).
- Merge of k sorted lists (k-way merge).