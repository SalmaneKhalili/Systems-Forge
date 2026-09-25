# M14-ex04 · Evict

A peer that has remained failed long enough must leave the live set, or the
cluster keeps routing work toward a node it cannot use. This exercise applies an
inclusive staleness bound to logical heartbeat sequences, turning the lifecycle's
failure state into a bounded membership set. You deliver `evict.go` with the
filter that retains every still-valid member in stable order.

## Shape

Harness-style: the provided `main.go` prints the transcript; you write
**`evict.go`** using only the Go standard library and implement:

```go
type Member struct {
	Node string
	Seq  int // sequence of the member's last heartbeat
}

// EvictStale returns the members still live at logical time now: those with
// now-Seq <= staleAfter, sorted by node name.
func EvictStale(members []Member, now, staleAfter int) []Member
```

A member with `now-Seq == staleAfter` is still live because the bound is
inclusive. Return live members sorted by node name. Never read the wall clock;
`now` is a logical sequence, not wall time. The reference transcript is
`expected.txt` with whitespace normalized.

## Acceptance

`make all` must build `test`; `./test` must print the reference transcript
exactly:

```text
now=10 staleAfter=3
members a7 b8 c9 d4
live: a7 b8 c9
evicted: d4
```

At logical time 10 with `staleAfter=3`, `a` at sequence 7, `b` at 8, and `c` at
9 stay within the bound; `d` at sequence 4 is 6 behind and is evicted. Evicting
a member exactly at the bound confuses `>=` with the required `>`, while keeping
a stale member defeats eviction.

## Readings

- **Reading ladder** — eviction keeps the live set bounded; the paper's cleanup rule does the
  age-based pruning.
- SWIM paper, evicting suspected members: https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf
- Datastax, "Failure Detector" (the remove threshold).
