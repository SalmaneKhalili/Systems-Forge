# M16-ex01 · Memtable

The LSM write path starts with a sorted buffer that can later flush entries in
order to an SSTable. This exercise makes every mutation visible through an ordered
in-memory view with overwrite and absence semantics. You deliver `memtable.go`,
the component ex02 serializes and the WAL later rebuilds.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`memtable.go`** using only the Go standard library and implement:

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

`Put` overwrites an existing key and preserves lexicographic order. `Get`
returns `("", false)` when absent, while `All` returns a copy rather than the
memtable's backing storage. Never read the wall clock; use synchronization
primitives if memory safety and concurrency require them. The reference
transcript is `expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
keys in order:
apple=1
banana=4
cherry=3
get banana: 4 (ok=true)
get dragon: (ok=false)
```

The ordered output proves that every key remains lexicographically sorted, and
the lookup lines prove the value/presence contract. A memtable that loses that
order or fails to overwrite an existing key fails the transcript.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2 "SSTables and LSM-Trees").
- Go `sort` package: https://pkg.go.dev/sort (specifically binary search with `sort.Search`).
