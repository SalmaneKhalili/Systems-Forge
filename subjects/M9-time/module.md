# M9 · Time & Ordering

Logical time has nothing to do with the wall clock: it is about **causality** —
deciding, from the messages exchanged, which events necessarily happened
before which others, and turning that partial order into something machines can
compare.

Milestones:

- ex01 · Lamport logical clock — a monotonically increasing counter that
  respects the happens-before order.
- ex02 · Vector clocks — per-process counters that can *test* causality in
  both directions.
- ex03 · Causality — decide happened-before and concurrency from vector
  stamps.
- ex04 · Total order — extend the partial order to a deterministic total
  order with a (lamport, pid) tie-break.
- ex05 · Sequencer — a real TCP gateway that assigns a Lamport-stamped,
  totally ordered sequence number to every frame it receives.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.