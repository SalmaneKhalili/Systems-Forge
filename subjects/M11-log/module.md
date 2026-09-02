# M11 · Replicated Log & Consistency

A replicated state machine drives every replica through the same ordered
sequence of entries — the **log** — so each converges to the same state. This
module builds the log mechanics: append-only entries, the commit point, index
math, snapshots, and a wire gateway that exposes a single consistent log to
clients.

Milestones:

- ex01 · Append — the log is append-only state, one entry at a time.
- ex02 · Committed — marking the committed prefix of the log.
- ex03 · Index — translating an index and length to the committed window and
  the next write slot.
- ex04 · Snapshot — compacting the log without losing the committed prefix.
- ex05 · The replica — a TCP gateway that serves one consistent committed
  log to many clients.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.