# M17-ex01 · Metrics

## Goal

Observability starts with **metrics**: small numeric measurements collected from a running service. A metrics registry tracks named counters and gauges that other components increment and read; a snapshot feeds dashboards and alerting.

Implement `metrics.go`:

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

Use a **read-write lock** (`sync.RWMutex`) so concurrent goroutines can update safely and
readers don't serialize with each other. This is your M4-ex04 rwlock revenge in Go: a
metrics registry is read far more than written (snapshots feed dashboards constantly while
counters are bumped only occasionally), so **reads take `RLock`, writes take `Lock`** — the
exact asymmetry you built in C. The provided `main.go` increments counters from a couple of
goroutines, snapshots them, and prints the result. `make all` must build `test`; `./test`
must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `metrics.go`. Define `KV` (a `Key, Val` pair, where `Val` is a `string`) in this package.
- Reads (SnapGet/Snapshot) and writes (Inc/Add) must be safe to call from multiple goroutines.

## Acceptance

Reference transcript:

```
snapshot (sorted):
conn_ok=5
conn_total=8
req_ok=3
req_total=3
total=19
total a+b=11
```

In the sample, three goroutines bump `req_total`/`req_ok`/`conn_total`, five bump `conn_total`/`conn_ok`, and `total` is seeded to 19.

A registry that loses increments from concurrent goroutines (not mutex-safe), or that fails to snapshot deterministically sorted, is the bug.

## Readings

- Tragically achievable with the Go standard library `sync.Mutex`.
- DDIA, Chapter 9 "Consistency and Consensus" (consistency guarantees for replicated
  measured state — summary for the module).