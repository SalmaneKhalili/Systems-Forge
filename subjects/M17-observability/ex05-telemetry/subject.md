# M17-ex05 · Telemetry Gateway (Mini-Capstone)

The concurrency-safe registry now needs an external surface for live updates and
snapshots. This Mini-Capstone lets operators mutate and read named counters over
TCP while the service itself accounts for every accepted connection. You deliver
`server.go`, the shared-state telemetry endpoint the module culminates in.

## Shape

Whole-program. You write **`server.go`** using only the Go standard library:

```go
func main() // reads TARGETPORT, serves TCP on 127.0.0.1:port, driving a Metrics
```

The server listens on `127.0.0.1` at the port in the `TARGETPORT` environment
variable and speaks this line protocol:

- `INC <name> <n>` → `ok` (increments a named counter by `n`)
- `GET <name>` → the current value as an integer, or `0` if absent
- `SNAP` → one `name=value` line per metric, sorted by name

Every accepted connection increments `conn_total`, so a second connection sees
the first connection in that count. Read `TARGETPORT` with `os.Getenv`, never a
fixed port, and guard the shared registry with a mutex. A second connection must
see counters written by the first. `make all` must build `test`; the grader
starts `./test`, runs commands over two connections, and checks the replies.

## Acceptance

`make all` must build `test`; the grader starts `./test` and requires these two
connections to produce exactly:

```text
conn 1:
  INC http_ok 3      -> ok
  INC http_total 10  -> ok
  GET http_ok        -> 3
  GET http_missing   -> 0
  SNAP               -> conn_total=1\nhttp_ok=3\nhttp_total=10\n

conn 2:
  GET conn_total     -> 2
  INC http_ok 4      -> ok
  GET http_ok        -> 7
  INC http_total 1   -> ok
  SNAP               -> conn_total=2\nhttp_ok=7\nhttp_total=11\n
```

The second connection's inherited values prove shared state, and the two
`conn_total` snapshots prove accept-time accounting. A fresh registry per
connection, a miscounted connection, an unsorted snapshot, or a wrong value
fails the transcript.

## Readings

- **Reading ladder** — OpenTelemetry metrics export:
  https://opentelemetry.io/docs/concepts/signals/metrics/
- OpenTelemetry metrics export over the network.
- Reuse `Metrics` from ex01 (mutex-safe registry with sorted snapshot).
