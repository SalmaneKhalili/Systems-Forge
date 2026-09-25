# M16-ex02 · SSTable

Once the memtable is full, its sorted entries can be flushed without another
in-place update path. This exercise serializes them into an immutable SSTable and
builds a sparse block index for boundary-safe lookup, the on-disk unit ex04 later
compacts. You deliver `sstable.go` with the file format, writer, opener, and
indexed `Get`.

## Shape

Harness-style: the provided `main.go` writes a temporary SSTable, reopens it,
and queries keys; you write **`sstable.go`** using only the Go standard library
and implement:

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

`KV` is reused from `memtable.go` in the same package. Imports from
`encoding/binary` and `sort` are permitted. Lookups must also work near block
boundaries, and never read the wall clock. The reference transcript is
`expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
sstable write ok
get zed -> (v0)
get kappa -> (v1)
get alpha -> ok=true
get amber -> ok=true
get missing -> ok=false
blocks: 3
```

The successful writes and reads prove serialization, reopening, and indexed
lookup; the three-block count proves the sparse index's block cadence.
Mishandling the last entry in a block or returning the wrong value across a
boundary breaks the file contract.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2 "SSTables and LSM-Trees").
- Go `encoding/binary`: https://pkg.go.dev/encoding/binary
