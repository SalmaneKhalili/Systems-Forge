# M17-ex05 · Telemetry Gateway (Mini-Capstone)

## Goal

Expose the metrics registry (ex01) to the outside world as a **telemetry gateway**: a TCP service that lets operators query and update live metrics, and the module's Mini-Capstone. The server increments a `conn_total` counter every time it accepts a new connection.

The server listens on `127.0.0.1` at the port in `TARGETPORT` (set by the grader) and speaks a line protocol:

- `INC <name> <n>` → `ok` (increments a named counter by `n`)
- `GET <name>` → the current value as an integer, or `0` if absent
- `SNAP` → one `name=value` line per metric, sorted by name

Every accepted connection increments `conn_total` (so a second connection sees the count from the first). Implement `server.go`:

```go
func main() // reads TARGETPORT, serves TCP on 127.0.0.1:port, driving a Metrics
```

`make all` must build `test`. The grader starts `./test`, runs commands across two connections, and checks replies.

## Constraints

- Go, standard library only; file is `server.go`.
- Shared state is required: a second connection must see counters written by the first.
- Read `TARGETPORT` from `os.Getenv`; never use a fixed port. Use a mutex around the registry.
- `SNAP` must list metrics in sorted name order.

## Acceptance

The grader runs these two connections and expects:

```
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

A server that serves each connection from a fresh registry, miscounts connections, or returns the wrong value is the bug.

## Readings

- OpenTelemetry metrics export over the network.
- Reuse `Metrics` from ex01 (mutex-safe registry with sorted snapshot).