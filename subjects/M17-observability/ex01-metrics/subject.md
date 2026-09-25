# M17-ex01 · Metrics

Observability starts with numeric state that can be updated while a service runs
and read consistently by dashboards. This exercise builds the concurrency-safe
registry behind that view, with deterministic snapshots for later tracing and
telemetry work. You deliver `metrics.go` and the read-write locking discipline
that keeps concurrent updates lossless.

## Shape

Harness-style: the provided `main.go` increments counters from several
goroutines, snapshots them, and prints the result; you write **`metrics.go`**
using only the Go standard library and implement:

```go
// Metrics is a concurrency-safe registry of named counters.
type Metrics struct{ /* your fields */ }

// NewMetrics initializes an empty registry.
func NewMetrics() *Metrics

// Inc increments the named counter by 1 (creating it if absent).
func (m *Metrics) Inc(name string)

// Add increments the named counter by an arbitrary amount.
func (m *Metrics) Add(name string, n int64)

// Get returns the current value of a counter and whether it exists.
func (m *Metrics) Get(name string) (int64, bool)

// Snapshot returns all counters as (name, value) pairs sorted by name.
func (m *Metrics) Snapshot() []KV
```

Use a **read-write lock** (`sync.RWMutex`): reads take `RLock`, while writes
take `Lock`. This is the M4-ex04 rwlock pattern in Go, where snapshots read far
more often than counters change, so readers need not serialize with one
another. Reads (`SnapGet`/`Snapshot`) and writes (`Inc`/`Add`) must be safe from
multiple goroutines. Define `KV` as a `Key, Val` pair where `Val` is a `string`
in this package.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
snapshot (sorted):
conn_ok=5
conn_total=8
req_ok=3
req_total=3
total=19
total a+b=11
```

Three goroutines bump `req_total`, `req_ok`, and `conn_total`; five bump
`conn_total` and `conn_ok`; `total` starts at 19. The final values prove no
concurrent increment is lost, and the sorted snapshot proves stable output. A
non-mutex-safe registry or a nondeterministic snapshot fails the transcript.

## Readings

- **Reading ladder** — OTel metrics primer first:
  https://opentelemetry.io/docs/concepts/signals/metrics/
- Tragically achievable with the Go standard library `sync.Mutex`.
- DDIA, Chapter 9 "Consistency and Consensus" (consistency guarantees for replicated
  measured state — summary for the module).
