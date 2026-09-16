# M17-ex02 · Tracing

## Goal

**Tracing** records the lifecycle of a single request as a tree of **spans**, where each span names a unit of work and carries its duration. A request consists of a root span with nested child spans (e.g. `request` → `db.query`, `cache.get`). Rendering the tree makes latency hot-spots visible.

Durations are **explicit** (passed in), never measured from a wall clock, to keep the grader deterministic.

Implement `trace.go`:

```go
// Span is one unit of work in a request tree.
type Span struct {
	name      string
	durMs     int64
	children  []*Span
}

// NewSpan creates a root span with the given name and duration.
func NewSpan(name string, durMs int64) *Span

// Child adds a child span under s and returns it.
func (s *Span) Child(name string, durMs int64) *Span

// Dump renders the tree depth-first as a slice of lines. Each line is
// indentation (2 spaces per depth) + name + "=" + durMs + "ms".
func (s *Span) Dump() []string

// TotalMs returns the span's own duration, not including children.
func (s *Span) TotalMs() int64
```

The provided `main.go` builds a request tree, dumps it, and prints the transcript. `make all` must build `test`; `./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `trace.go`.
- No wall clock: durations are explicit arguments.
- Dump must visit spans in insertion order and indent by depth.

## Acceptance

Reference transcript:

```
request=120ms
  db.query=40ms
    db.scan=25ms
  cache.get=15ms
  api.out=65ms
root.dur=120
```

The root `request` (120) contains `db.query` (40, itself containing `db.scan` 25), `cache.get` (15), and `api.out` (65). A tracer that renders children in the wrong order, mis-tracks depth, or fabricates durations from a clock is the bug.

## Readings

- OpenTelemetry concepts: spans, trace context.
- DDIA, Chapter 11 (Stream processing) for the disambiguation of distributed request trees.