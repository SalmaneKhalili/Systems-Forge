# M17-ex02 · Tracing

The registry shows aggregate state, but tracing follows one request through
named units of work. This exercise models that lifecycle as a root span and
nested child spans, then renders the hierarchy so latency hot-spots remain
visible. You deliver `trace.go` with insertion-ordered depth-first output and
explicit durations.

## Shape

Harness-style: the provided `main.go` builds and dumps a request tree; you write
**`trace.go`** using only the Go standard library and implement:

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

Durations are explicit arguments, never values measured from a wall clock.
`Dump` visits spans in insertion order and indents by depth with two spaces per
level. `TotalMs` returns only the span's own duration, excluding children.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
request=120ms
  db.query=40ms
    db.scan=25ms
  cache.get=15ms
  api.out=65ms
root.dur=120
```

The root `request` (120) contains `db.query` (40, itself containing `db.scan`
at 25), `cache.get` (15), and `api.out` (65). The tree proves insertion order,
depth indentation, and parent-child structure; the final line proves `TotalMs`
excludes child time. Wrong child order, depth, or clock-derived duration fails
the transcript.

## Readings

- **Reading ladder** — OpenTelemetry tracing concepts:
  https://opentelemetry.io/docs/concepts/signals/traces/
- OpenTelemetry concepts: spans, trace context.
- DDIA, Chapter 11 (Stream processing) for the disambiguation of distributed request trees.
