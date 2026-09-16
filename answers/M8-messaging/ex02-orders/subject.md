# M8-ex02 — ordered delivery

## Goal

Messages arrive at a consumer **out of order** (partitions advance unevenly:
partition 0's offset 0 can land after its offset 1). A correct consumer holds
back a partition's later offsets until their missing predecessor shows up. You
implement an ordered queue that turns a scripted, deliberately-tangled arrival
order into the single allowed delivery order.

## Shape

Harness-style. The exercise ships `main.go` (the driver), `go.mod`, `Makefile`.
You write **`queue.go`** with one constructor and two methods:

```go
func NewQueue() *queue
func (q *queue) Put(partition, offset int)   // one arriving message
func (q *queue) Drain() ([]entry, error)     // the allowed delivery order
```

`entry{partition, offset int}` is defined by `main.go`. Rules your queue keeps:

- **Per-partition FIFO.** Two entries of the same partition must never be
  delivered out of offset order.
- **Gap holding.** When an entry's predecessor has not arrived yet, it is held,
  however late its model chases the rest of the topic.
- **Earliest arrival wins.** Whenever more than one entry is deliverable at
  once, the one that **arrived first** is delivered first.
- If a drain ends with an unrecoverable per-partition gap, `Drain` returns a
  non-nil error instead of silently reordering.

## Acceptance criteria (all graded)

`make all`, then `./test` must exit 0 and print exactly:

```text
delivered: b0|a0|a1|b1|a2|b2
```

The driver feeds, in this arrival order: `a1, b0, a0, b1, a2, b2`. The only
legal delivery order — holding `a1` until `a0` arrives, then releasing by
earliest-arrival — is the line above. A queue that just replays arrival order
prints `a1|b0|a0|b1|a2|b2`; a queue that sorts by partition prints
`a0|a1|a2|b0|b1|b2`; both are wrong.

Stderr must stay empty. Clean exit (0) is part of the grade.

## Readings

- Jay Kreps, *The Log: What every software engineer should know about
  real-time data's unifying abstraction* — partition ordering, offsets.
- Go maps and slices — an ordered consolidation is map-driven index work.

## Quiz

1. What must hold between any two entries of the same partition?
2. What does the queue do with a partition entry whose predecessor has not arrived?
3. When several entries are eligible at once, whose entry is delivered first?