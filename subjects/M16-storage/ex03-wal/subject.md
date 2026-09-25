# M16-ex03 · Write-Ahead Log (WAL)

The memtable and SSTable are only durable after the disk handoff, so every
mutation first needs an append-only record outside memory. This exercise defines
the length-prefixed WAL and deterministic replay used to reconstruct sorted state
after a crash. You deliver `wal.go`, the recovery record the persistent gateway
writes before applying a change.

## Shape

Harness-style: the provided `main.go` appends writes, opens the file fresh to
simulate a crash, reads the log back, and replays it; you write **`wal.go`** using
only the Go standard library and implement:

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

Define `KV` as a `Key, Val string` pair in this package. The log is append-only:
never rewrite it in place, and preserve write order during replay. Never read the
wall clock.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
append ok
log entries: 4
replay:
a=1
b=2
c=3
d=4
```

The entry count proves the length-prefixed records survive the fresh open, and
the replay proves ordered application and key-sorted output. Applying writes out
of order or failing to preserve the last write to a key breaks recovery.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2 "SSTables and LSM-Trees": the write path and crash recovery via WAL).
- This section: "The write path: memtable + WAL".
