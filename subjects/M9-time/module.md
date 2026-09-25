# M9 · Time & Ordering

This module turns message exchange into explicit causal facts. A scalar logical
clock becomes a per-process vector, vector stamps expose happened-before and
concurrency, and a deterministic tie-break turns the resulting partial order
into the sequence assigned by a TCP sequencer.

## The build

- **ex01 · Lamport logical clock** — a monotonically increasing counter that
  respects happens-before; the scalar clock that hands ex02 its first vector
  component.
- **ex02 · Vector clocks** — per-process counters that can *test* causality in
  both directions; the stamps that hand ex03 a relation to decide.
- **ex03 · Causality** — classify two stamps as happened-before, happened-after
  or concurrent; the relation that exposes the partial order ex04 must totalize.
- **ex04 · Total order** — extend the partial order to a deterministic total
  order with a `(lamport, pid)` tie-break.
- **ex05 · Sequencer** — a real TCP gateway that assigns a Lamport-stamped,
  totally ordered sequence number to every frame it receives.

## Rules

Nothing printed or compared may depend on wall-clock time. All exercises are
transcript-driven and byte-deterministic.

## Prerequisites

Before starting M9, you should be comfortable with everything from M0–M8, plus:

- Explain what a Lamport logical clock is: a counter attached to each event, incremented on
  local events and updated on message receipt (`max(local, received) + 1`).
- Explain what a vector clock is: an array of counters, one per node, updated element-wise
  on events and message receipt (`for each i: vc[i] = max(local[i], received[i])`).
- Explain causality: event A "happened before" event B if there is a chain of messages
  from A to B (or A is a local predecessor of B).
- Explain total order vs partial order: total order means every event pair is comparable;
  partial order means some pairs are concurrent (neither caused the other).
- Explain what a sequencer is: a designated node that stamps every message with a global
  sequence number, giving total order without tracking causality.
- Use `net.Listen`/`net.Dial` and `bufio.Scanner` (from M7/M8) — you will build a TCP
  sequencer server in ex05.
- Write a goroutine that handles one TCP connection (ex05 requires spawning a handler
  goroutine per client).

You do NOT need to know: Raft, 2PC, Paxos, or distributed consensus. Those come in M10–M12.

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
