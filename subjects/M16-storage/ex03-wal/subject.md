# M16-ex03 · Write-Ahead Log (WAL)

## Goal

An LSM-tree buffers writes in a memtable (ex01), but memtables are volatile. To survive a crash before the next flush, every mutation is first appended to an append-only **write-ahead log** (WAL) on disk, then applied to the memtable. On restart, the WAL is **replayed** to rebuild the in-memory state.

Implement `wal.go`:

```go
// LogEntry is one recorded operation.
type LogEntry struct {
	Key, Val string
}

// AppendLog appends entries to the log at path in an append-only fashion.
// If the file does not exist it is created. Each entry is written as a
// length-prefixed record so the log can be read back incrementally.
func AppendLog(path string, entries []LogEntry) error

// ReadLog reads back every entry in the log at path in order.
func ReadLog(path string) ([]LogEntry, error)

// Replay returns the map of key -> value reconstructed by applying every
// entry in order (later writes overwrite earlier ones). Because keys must
// stay sorted for a memtable, the result is returned sorted by key.
func Replay(entries []LogEntry) []KV
```

The provided `main.go` appends writes (simulating a crash by opening the file fresh), reads the log back, and replays it. `make all` must build `test`; `./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `wal.go`. Define `KV` (a `Key, Val string` pair) in this package.
- Append-only: never rewrite the log in place; replay order must preserve write order.

## Acceptance

Reference transcript:

```
append ok
log entries: 4
replay:
a=1
b=2
c=3
d=4
```

A replay that applies writes out of order, or that does not preserve the last write to a key, is the bug.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3: "SSTables and LSM-Trees" (The write path and recovery).
- This section: "The write path: memtable + WAL".