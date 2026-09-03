# M16-ex05 · Persistent KV Gateway (Mini-Capstone)

## Goal

Bring the LSM-tree pieces together into a **persistent key-value store** with a TCP gateway (this module's Mini-Capstone). All mutations are appended to a write-ahead log (ex03) so the store survives a restart; the in-memory map is rebuilt by replaying the log on startup (ex01/ex03), and keys stay sorted (ex01).

The server listens on `127.0.0.1` at the port given in the environment variable `TARGETPORT` (set by the grader) and speaks a simple line protocol:

- `PUT <key> <value>` → `ok`
- `GET <key>` → the value, or `nil` if absent
- `DEL <key>` → `ok` (removes the key)

Every PUT/DEL is appended to a log file named `store.log` **in the working directory** before being applied to memory, so a later restart (including one within the same grader run) replays it.

Implement `server.go`:

```go
func main() // reads TARGETPORT, loads store.log, serves TCP on 127.0.0.1:port
```

`make all` must build `test`. The grader starts `./test`, sends commands over two consecutive TCP connections, and checks the replies.

## Constraints

- Go, standard library only; file is `server.go`.
- Persistence is required: a second connection must still see data written by the first (both are served by the same process; on a real restart the log replay restores it).
- Read `TARGETPORT` from `os.Getenv`; never use a fixed port.
- Handle concurrent connections with a mutex around the store.

## Acceptance

The grader runs these two connections and expects:

```
conn 1:
  PUT alpha 1        -> ok
  GET alpha          -> 1
  PUT beta 2         -> ok
  PUT alpha 9        -> ok
  GET alpha          -> 9
  DEL gamma          -> ok
  GET gamma          -> nil

conn 2:
  GET alpha          -> 9
  PUT gamma 7        -> ok
  GET beta           -> 2
  GET alpha          -> 9
  GET gamma          -> 7
```

A server that serves each connection from a fresh (non-persistent) map, returns stale values, or answers the wrong value, is the bug.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2 "SSTables and LSM-Trees"): the full write path and recovery.
- Reuse ex01 (memtable), ex02 (sstable), ex03 (wal), ex04 (compaction).