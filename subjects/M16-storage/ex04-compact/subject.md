# M16-ex04 · Compaction

Immutable SSTables accumulate overlapping runs, which makes reads fan out across
more files. This exercise merges those sorted runs into one, resolves duplicate
keys by run recency, and removes tombstones so stale data and deleted keys do not
survive cleanup. You deliver `compact.go`, the merge that bounds the on-disk read
path before the persistent gateway joins the engine.

## Shape

Harness-style: the provided `main.go` merges a small set of overlapping runs and
prints the result; you write **`compact.go`** using only the Go standard library
and implement:

```go
// Compact merges several sorted runs (each a sorted []KV) into one sorted
// run. Each key appears at most once; when a key appears in more than one
// run, the later run (higher index) wins. A KV with Val == "" acts as a
// tombstone and removes the key entirely.
func Compact(runs ...[]KV) []KV
```

Define `KV` as a `Key, Val string` pair in this package. The output must be
strictly sorted by key with no duplicate keys. A tombstone is represented by
`Val == ""` and removes the key entirely, even when a later run contains only
that tombstone. Never read the wall clock.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
after compaction:
a=stale
b=new
d=3
f=tombstoned
g=1
```

In the sample, `a` appears in both runs and the newer `stale` value wins; `c` and
`e` were written and then tombstoned, so they are absent. `b`, `d`, and `f`
survive from the newer run, while `g` survives from the older run. Duplicate
keys, a missing surviving key, or a surviving tombstoned key fails compaction.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2.3 "Performance Optimizations" — compaction strategies and bloom filters).
- Merge of k sorted lists (k-way merge).
