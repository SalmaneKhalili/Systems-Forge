# M16-ex02 · SSTable

## Goal

When a memtable fills up (ex01), an LSM-tree **flushes** its sorted entries to an immutable on-disk file called an **SSTable** (Sorted String Table). Because SSTables are immutable and sorted, we never rewrite them in place — we compact them (ex04). Lookup uses a sparse index plus binary search rather than scanning the whole file.

Implement `sstable.go`:

```go
// WriteSSTable serializes a sorted slice of KV pairs to a file at path.
// Format: a data section of "key<0x00>value<0x00>" records followed by a
// 2-byte little-endian index offset; at the index offset lives a list of
// records giving, for each block (every 4 entries), the block's first key
// (NUL-terminated) and a 4-byte little-endian offset into the data section.
func WriteSSTable(path string, items []KV) error

// OpenSSTable loads the sparse index into memory.
func OpenSSTable(path string) (*SSTable, error)

// Get returns the value for a key using block index + binary search, or
// ("", false) if absent.
func (s *SSTable) Get(key string) (string, bool)
```

The provided `main.go` writes a temp SSTable, reopens it, and queries keys. `make all` must build `test`; `./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `sstable.go`. `KV` is reused from `memtable.go` (same package).
- Imports (`encoding/binary`, `sort`) are permitted.
- Reference transcript is `expected.txt` (whitespace normalized).
- No wall clock. Looks up must also work near the block boundaries.

## Acceptance

Reference transcript:

```
sstable write ok
get zed -> (v0)
get kappa -> (v1)
get alpha -> ok=true
get amber -> ok=true
get missing -> ok=false
blocks: 3
```

A lookup that mishandles the last entry in a block, or returns the wrong value across a block boundary, is the bug.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2 "SSTables and LSM-Trees").
- Go `encoding/binary`: https://pkg.go.dev/encoding/binary
