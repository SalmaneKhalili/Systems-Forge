# M17 · Observability

A distributed system you cannot see is a distributed system you cannot debug.
**Observability** is the tooling that makes a production service legible:
**metrics** (numeric counters and gauges for dashboards and alerting),
**tracing** (a request's lifecycle as a tree of named spans with durations), and
**summaries** (aggregating many latencies into count/max/percentiles for SLOs).

Milestones:

- ex01 · Metrics — a concurrency-safe registry of counters and gauges with a
  snapshot.
- ex02 · Tracing — a request rendered as a root span with nested child spans.
- ex03 · Spans with attributes — attach key=value tags and query the collected
  tree (match, slowest).
- ex04 · Summaries — aggregate durations into count/sum/max, p50 and p95.
- ex05 · Telemetry gateway — the Mini-Capstone: a TCP service exposing the
  metrics registry, with a live `conn_total`.

Rule: durations are **explicit** (passed in), never measured from a wall
clock — the whole module stays byte-deterministic.

---

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