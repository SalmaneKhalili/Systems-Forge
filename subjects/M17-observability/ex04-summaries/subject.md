# M17-ex04 · Latency Summaries

Individual spans expose work, but a glanceable latency view needs aggregation
across many observations. This exercise reduces explicit durations to count,
sum, maximum, and deterministic quantiles such as p50 and p95 for SLOs. You
deliver `summary.go` with the exact index convention the telemetry handoff uses.

## Shape

Harness-style: the provided `main.go` records durations and prints their
summary; you write **`summary.go`** using only the Go standard library and
implement:

```go
// Summary aggregates observed durations in milliseconds.
type Summary struct{ /* your fields */ }

// NewSummary initializes an empty summary.
func NewSummary() *Summary

// Add records one observed duration in ms.
func (s *Summary) Add(durMs int64)

// Count returns the number of observations.
func (s *Summary) Count() int64

// Sum returns the total of all observations.
func (s *Summary) Sum() int64

// Max returns the largest observation, or 0 if empty.
func (s *Summary) Max() int64

// Quantile returns the value at the given quantile (0..1): the smallest
// observation at index ceil(q*n)-1 (1-indexed). For empty data it returns 0.
func (s *Summary) Quantile(q float64) float64

// SummaryValues returns the key/values: count, sum, max, and quantiles at
// 0.5 and 0.95, rendered as strings.
func (s *Summary) SummaryValues() []KV
```

Define `KV` as a `Key, Val string` pair in this package. Quantiles must remain
deterministic, and `SummaryValues` must include the keys `count`, `sum`, `max`,
`p50`, and `p95` in that order. Durations are passed in, never measured from a
wall clock.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
count=5
sum=150
max=50
p50=30
p95=50
```

For durations 10, 20, 30, 40, and 50, the count is 5, sum is 150, maximum is
50, p50 is 30, and p95 is 50. A wrong count, sum, maximum, or quantile breaks
the summary's deterministic contract.

## Readings

- Percentiles and p95/p99 for latency SLOs (DDIA, Chapter 1 "Reliable, Scalable, and
  Maintainable Applications", §1.3 "Describing Performance").
- Quantile rounding (ceil index) convention.
