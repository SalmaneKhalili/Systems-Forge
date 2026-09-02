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

---

## So what? (interview / portfolio)

The replicated state machine is the canonical "how a service stays consistent across
replicas" answer, and the append-only log is the data structure at its heart. Being able to
say "every replica replays the same ordered log, so they converge" — and to defend the
commit point — is a core distributed-systems interview skill.

**Interview questions this module arms you for:**
- Why does an append-only log make replicas deterministic and consistent by construction?
- What is the difference between the log length and its committed prefix?
- How does a snapshot let you truncate a log without losing the committed state?
- Why is the order of entries sacred to a replicated state machine?

**Portfolio artifact:** M11-ex05 `replica` — a TCP gateway that serves one consistent
committed log to many clients.