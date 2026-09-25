# M17 · Observability

This module makes the finished distributed stack legible through numeric metrics,
request traces, and latency summaries. The thread runs from a concurrency-safe
registry to nested spans, attributed-span queries, and quantile aggregation,
then exposes the live registry through a TCP telemetry gateway.

## The build

- **ex01 · Metrics** — maintain named counters behind a read-write lock and
  return deterministic name-sorted snapshots.
- **ex02 · Tracing** — model one request as a root span with nested children and
  render the tree in insertion order.
- **ex03 · Spans with attributes** — attach key=value tags, query matching spans,
  and find the slowest own-duration in the collected tree.
- **ex04 · Summaries** — aggregate explicit durations into count, sum, maximum,
  p50, and p95 for SLO reporting.
- **ex05 · Telemetry gateway** — **Gate**: serve live registry updates and sorted
  snapshots over TCP while counting every accepted connection in `conn_total`.

## Rules

Durations are **explicit** (passed in), never measured from a wall clock — the
whole module stays byte-deterministic.

## Prerequisites

Before starting M17, you should be comfortable with everything from M0–M16, plus:

- Explain what observability is: the ability to understand a system's internal state from
  its external outputs (metrics, logs, traces).
- Explain the three pillars:
  - **Metrics:** numeric counters and gauges (e.g., "requests per second", "queue depth").
  - **Logs:** timestamped text records of events.
  - **Traces:** a tree of spans showing the path of a single request through the system.
- Explain a metric type: `Counter` (monotonically increasing — only goes up) vs `Gauge`
  (can go up or down — e.g., queue depth, active connections).
- Explain the summary/aggregation problem: every node collects its own metrics; the summary
  endpoint must aggregate them into a cluster-wide view.
- Explain the tradeoff of push vs pull: push = nodes send metrics to a central server on
  a timer; pull = central server queries nodes. Here we use push.
- Explain what OpenTelemetry is: a vendor-neutral standard for metrics/traces. You will NOT
  implement OTel here — you will build the same concepts from scratch.
- Use `net.Listen`/`net.Dial` and `bufio.Scanner` — ex05 is a TCP telemetry gateway.
- Use `sync.RWMutex` (from M4) for concurrent access to the metrics registry. Reads are
  fast (RLock); writes are exclusive (Lock).
- Use goroutines for concurrent metric updates (multiple connections updating metrics
  simultaneously).

You do NOT need to know: Prometheus, Grafana, Jaeger, or any external tooling. This module
builds the concepts from scratch.

## So what? (interview / portfolio)

Distributed systems interviews and on-call reality revolve around "how do you
know it's broken and where?" — this module is your answer. ex05 `telemetry`
gives you a runnable gateway that emits live counters over TCP, and you can
speak fluently to the metrics/traces/logs triad, percentiles vs. averages, and
why quantile (p50/p95) beats mean for latency SLOs.

**Interview questions this module arms you for:**
- Metrics vs. traces vs. logs: what is each for, and how do they complement?
- Why does monitoring use percentiles (p50/p95) instead of the average?
- What is a concurrency-safe registry, and why does it matter?
- How would you surface an SLO violation from live telemetry?

**Portfolio artifact:** M17-ex05 `telemetry` — a TCP telemetry gateway exposing
a concurrency-safe metrics registry with live `conn_total` counters.
