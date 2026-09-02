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

---

## So what? (interview / portfolio)

"Logical time" is the single most at-parity-with-distributed-systems idea you will ever
learn: every leader election, version counter, and total-order total you see later is
Lamport/vector clocks under the hood. Conveying *happens-before* fluently — and knowing why
wall-clock time is useless for causality — is a classic senior-distributed-systems
interview hinge.

**Interview questions this module arms you for:**
- What is the happens-before relation, and why is a wall clock the wrong clock for it?
- Lamport vs. vector clocks: what can each one *prove* about causality?
- How do you turn a partial order into a deterministic total order?
- When are two events concurrent (largely incomparable), and why does that matter?

**Portfolio artifact:** M9-ex05 `sequencer` — a TCP gateway that Lamport-stamps and totally
orders every frame it receives.