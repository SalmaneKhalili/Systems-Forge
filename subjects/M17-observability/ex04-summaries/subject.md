# M17-ex04 · Latency Summaries

## Goal

Raw spans are too noisy to review at a glance. A **summary** aggregates many observed durations into a small set of statistics: count, sum, maximum, and **quantiles** (percentiles) such as p50 and p95 used for SLOs.

Implement `summary.go`:

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

The provided `main.go` records a set of durations and prints the summary. `make all` must build `test`; `./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `summary.go`. Define `KV` (a `Key, Val string` pair) in this package.
- Quantiles must be deterministic. `SummaryValues` should include keys `count`, `sum`, `max`, `p50`, `p95` in that order.

## Acceptance

Reference transcript:

```
count=5
sum=150
max=50
p50=30
p95=50
```

For recorded durations 10, 20, 30, 40, 50: count 5, sum 150, max 50, p50 30, p95 50. A summary that miscounts, sums, reports the wrong max, or computes a wrong quantile is the bug.

## Readings

- Percentiles and p95/p99 for latency SLOs (DDIA, Chapter 1 "Reliable, Scalable, and
  Maintainable Applications", §1.3 "Describing Performance").
- Quantile rounding (ceil index) convention.