# M16-ex01 · Memtable

## Goal

The **memtable** is the write buffer in an LSM-tree (Log-Structured Merge-tree) storage engine. All incoming writes (`Put`) are buffered in memory before being flushed to disk. Crucially, the memtable must keep its keys **sorted** to allow efficient range queries and sequential disk writes during flush.

Implement `memtable.go`:

```go
type KV struct {
	Key, Val string
}

type MemTable struct {
	// your fields
}

// NewMemTable initializes an empty memtable.
func NewMemTable() *MemTable

// Put inserts a key-value pair. If the key already exists, its value is updated.
// The keys must remain sorted lexicographically.
func (m *MemTable) Put(key, val string)

// Get retrieves the value associated with a key, returning (value, true) if found,
// or ("", false) if absent.
func (m *MemTable) Get(key string) (string, bool)

// All returns a copy of all key-value pairs sorted in lexicographical order.
func (m *MemTable) All() []KV
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `memtable.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- No wall clock.
- Memory safety and concurrency: use synchronization primitives if required.

## Acceptance

Reference transcript:

```
keys in order:
apple=1
banana=4
cherry=3
get banana: 4 (ok=true)
get dragon: (ok=false)
```

A memtable that fails to keep keys sorted, or that does not overwrite an existing key, is the bug.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2 "SSTables and LSM-Trees").
- Go `sort` package: https://pkg.go.dev/sort (specifically binary search with `sort.Search`).
