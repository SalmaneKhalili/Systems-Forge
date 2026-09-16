# M14-ex04 · Evict

## Goal

A failed node's membership must eventually be dropped, or the cluster will keep
sending it traffic forever. Each member carries the sequence number of its last
heartbeat; a node whose `now - seq` exceeds a **staleness bound** is evicted
from the live set. Members at or inside the bound stay live.

Implement `evict.go`:

```go
type Member struct {
	Node string
	Seq  int // sequence of the member's last heartbeat
}

// EvictStale returns the members still live at logical time now: those with
// now-Seq <= staleAfter, sorted by node name.
func EvictStale(members []Member, now, staleAfter int) []Member
```

The provided `main.go` prints the transcript. `make all` must build `test`;
`./test` must print the reference transcript exactly.

## Constraints

- Go, standard library only; file is `evict.go`.
- Reference transcript is `expected.txt` (whitespace normalized).
- A member with `now-Seq == staleAfter` is still live (the bound is inclusive).
- Never read the wall clock — `now` is a logical sequence, not wall time.

## Acceptance

Reference transcript:

```
now=10 staleAfter=3
members a7 b8 c9 d4
live: a7 b8 c9
evicted: d4
```

At logical time 10 with `staleAfter=3`, `a` (seq 7), `b` (8) and `c` (9) are
within the bound and stay, while `d` (seq 4) is 6 behind and is evicted. A
sieve that evicts a member exactly at the bound (off-by-one `>=` vs `>`) or
that keeps a stale member is the bug.