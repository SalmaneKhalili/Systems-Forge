# M8-ex02 · Ordered delivery

A consumer sees a topic's partitions advance unevenly: partition 0's offset 0
can land after its offset 1. Your job is the queue that turns that
deliberately-tangled arrival order into the single allowed delivery order —
holding a partition's later offsets until the missing predecessor shows up.
The ordering obligation you prove here is the same contract the switch (ex03)
and the M9+ replicas depend on.

## Shape

Harness-style: the exercise ships `main.go` (the driver), `go.mod` and a
`Makefile`; you write **`queue.go`** with one constructor and two methods:

```go
func NewQueue() *queue
func (q *queue) Put(partition, offset int)   // one arriving message
func (q *queue) Drain() ([]entry, error)     // the allowed delivery order
```

`entry{partition, offset int}` is defined by `main.go`. The queue must hold:

- **Per-partition FIFO.** Two entries of the same partition must never be
  delivered out of offset order.
- **Gap holding.** An entry whose predecessor has not arrived is held, however
  late its model chases the rest of the topic.
- **Earliest arrival wins.** Whenever more than one entry is deliverable at
  once, the one that **arrived first** is delivered first.
- A drain that ends with an unrecoverable per-partition gap returns a non-nil
  error instead of silently reordering.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
delivered: b0|a0|a1|b1|a2|b2
```

The driver feeds, in this arrival order: `a1, b0, a0, b1, a2, b2`. The only
legal delivery order — holding `a1` until `a0` arrives, then releasing by
earliest arrival — is the line above. A queue that just replays arrival order
prints `a1|b0|a0|b1|a2|b2`; a queue that sorts by partition prints
`a0|a1|a2|b0|b1|b2`; both are wrong.

Stderr must stay empty. A clean exit (0) is part of the grade. Graded
`build` + `stdout` + `quiz`.

## Readings

- **Reading ladder** — the essay first (partition ordering, offsets):
  https://www.confluent.io/blog/log-what-every-software-engineer-should-know-about-real-time-datas-unifying/
- The Log's partition-ordering and offset mechanics — the module's framing.
- Go maps and slices — an ordered consolidation is map-driven index work.

## Quiz

1. What must hold between any two entries of the same partition?
2. What does the queue do with a partition entry whose predecessor has not arrived?
3. When several entries are eligible at once, whose entry is delivered first?