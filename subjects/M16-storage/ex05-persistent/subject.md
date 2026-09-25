# M16-ex05 · Persistent KV Gateway (Mini-Capstone)

The module's storage path now reaches a network boundary: ordered memtable
entries, immutable SSTables, and pre-apply WAL records must behave as one
persistent key-value service. This Mini-Capstone serves TCP clients, logs each
mutation before applying it, and rebuilds memory from `store.log` on startup.
You deliver `server.go`, the runnable LSM gateway.

## Shape

Whole-program. You write **`server.go`** using only the Go standard library:

```go
func main() // reads TARGETPORT, loads store.log, serves TCP on 127.0.0.1:port
```

The server listens on `127.0.0.1` at the port supplied in the `TARGETPORT`
environment variable. Its line protocol is:

- `PUT <key> <value>` → `ok`
- `GET <key>` → the value, or `nil` if absent
- `DEL <key>` → `ok` (removes the key)

Append every PUT/DEL to `store.log` **in the working directory** before applying
it to memory, then rebuild the in-memory map by replaying that log on startup;
keys remain sorted. Read `TARGETPORT` with `os.Getenv`, never a fixed port, and
protect the store with a mutex while serving concurrent connections.
Persistence is required: a second connection still sees the first connection's
data, and a real restart restores it through log replay. `make all` must build
`test`; the grader starts `./test` and checks replies over two consecutive TCP
connections.

## Acceptance

`make all` must build `test`; the grader starts `./test` and requires these two
connections to produce exactly:

```text
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

The grader then sends **SIGTERM** to the process and requires it to exit **0**.
This re-deploys the M7-ex04 graceful-shutdown skill: a crash would kill the
process with a non-zero exit and could leave the WAL unsynced. A correct
server installs a SIGTERM handler, stops accepting new connections, flushes
and closes `store.log`, and exits 0.
```

A fresh non-persistent map per connection, a stale value, or the wrong reply
fails the transcript. A server that dies on `SIGTERM` with a non-zero exit also
fails the graceful-shutdown step.

## Readings

- *Designing Data-Intensive Applications* (DDIA), Chapter 3 "Storage and Retrieval"
  (§3.2 "SSTables and LSM-Trees"): the full write path and recovery.
- Reuse ex01 (memtable), ex02 (sstable), ex03 (wal), ex04 (compaction).
