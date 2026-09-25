# M17-ex03 · Spans with Attributes

A rendered request tree becomes searchable when each span also carries key=value
attributes. This exercise adds those tags, preserves distinct same-named spans,
and answers tree-wide match and slowest-span queries. You deliver `spans.go`
with attributed rendering, depth-first lookup, and recursive duration analysis.

## Shape

Harness-style: the provided `main.go` builds two `db.query` spans in different
subtrees, attaches attributes, and queries the result; you write **`spans.go`**
using only the Go standard library and implement:

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

Durations remain explicit; never read the wall clock. `FindAll` returns every
matching span in depth-first order, and `MaxDur` considers the current span and
every descendant's own duration.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
request{}=30ms
  db.query{db=mysql}=40ms
    db.scan{}=25ms
  cache.get{hit=true}=15ms
  db.query{db=postgres}=50ms
count db.query=2
slowest=50
```

The tree contains two distinct `db.query` spans, one tagged `db=mysql` and one
`db=postgres`, so `FindAll` returns 2. The slowest span is a descendant with own
duration 50, proving `MaxDur` recurses. Merging same-named spans, losing an
attribute, or reporting the wrong slowest span fails the contract.

## Readings

- **Reading ladder** — OpenTelemetry Span semantics:
  https://opentelemetry.io/docs/concepts/signals/traces/#span
- OpenTelemetry Span attributes and semantic conventions.
- Distributed tracing querying (finding hot spans).
