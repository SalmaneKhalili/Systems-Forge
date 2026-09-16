# M17-ex03 · Spans with Attributes

## Goal

A **span** carries more than a name and duration: it also carries **attributes** (key=value tags) that enrich it, e.g. `db=f'postgres'`, `hit=true`. A tracer must let components attach attributes to spans and then answer queries over the collected tree: how many spans match a name, and which span is the slowest.

Durations remain **explicit** (no wall clock).

Implement `spans.go`:

```go
// Span is a named unit of work with a duration and key=value attributes.
type Span struct {
	name  string
	durMs int64
	attr  map[string]string
	children []*Span
}

// NewSpan creates a root span.
func NewSpan(name string, durMs int64) *Span

// Child adds a child span under s and returns it.
func (s *Span) Child(name string, durMs int64) *Span

// Attr sets an attribute key=value on the span.
func (s *Span) Attr(k, v string)

// Dump renders each span depth-first as "name{dur=Ms,k=v,...}" with 2-space
// indentation per depth.
func (s *Span) Dump() []string

// FindAll returns every span (in any subtree) whose name matches, in
// depth-first order.
func (s *Span) FindAll(name string) []*Span

// MaxDur returns the duration of the slowest span among s and all
// descendants (own durations only).
func (s *Span) MaxDur() int64
```

The provided `main.go` builds a request with two `db.query` spans in different subtrees, attaches attributes, then queries. `make all` must build `test`; `./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `spans.go`.
- No wall clock.
- FindAll must return spans in depth-first order; MaxDur must consider every descendant.

## Acceptance

Reference transcript:

```
request{}=30ms
  db.query{db=mysql}=40ms
    db.scan{}=25ms
  cache.get{hit=true}=15ms
  db.query{db=postgres}=50ms
count db.query=2
slowest=50
```

The tree has two distinct `db.query` spans (one with `db=mysql`, the other with `db=postgres`), so `FindAll` returns 2 and the slowest own-duration is a descendant (50), proving `MaxDur` recurses. A tracer that merges same-named spans, loses attributes, or reports the wrong slowest span is the bug.

## Readings

- **Reading ladder** — OpenTelemetry Span semantics:
  https://opentelemetry.io/docs/concepts/signals/traces/#span
- OpenTelemetry Span attributes and semantic conventions.
- Distributed tracing querying (finding hot spans).