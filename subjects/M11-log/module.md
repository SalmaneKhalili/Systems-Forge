# M11 · Replicated Log & Consistency

A replicated state machine drives every replica through the same ordered
entries, so this module builds the log those replicas share. The sequence moves
from append-only storage and its commit boundary through index arithmetic and
compaction, then exposes one process-wide committed log over TCP.

## The build

- **ex01 · Append** — append-only state that assigns one stable public index per
  command; the ordered log ex02 marks as durable.
- **ex02 · Committed** — track the commit index and expose only the durable
  prefix; the boundary whose arithmetic ex03 isolates.
- **ex03 · Index** — translate an index and length into the next write slot and
  uncommitted count; the arithmetic that keeps public indexes stable for ex04.
- **ex04 · Snapshot** — compact the committed prefix without losing or
  renumbering the surviving suffix.
- **ex05 · The replica** — a TCP gateway that serves one consistent committed
  log to many clients.

## Rules

Nothing printed or compared may depend on wall-clock time. All exercises are
transcript-driven and byte-deterministic.

## Prerequisites

Before starting M11, you should be comfortable with everything from M0–M10, plus:

- Explain what a replicated state machine is: multiple replicas process the same ordered
  sequence of commands, so they all converge to the same state.
- Explain what a log is in this context: an append-only sequence of entries, each containing
  a command. The log IS the state machine's input.
- Explain the commit point: the index up to which all replicas have the entry persisted.
  A committed entry will NOT be lost, even if a leader crashes.
- Explain the relationship between log index and state: entry at index 1 is applied first,
  then index 2, etc. Applying out of order = incorrect state.
- Explain what a snapshot is: a compacted representation of state up to index N, so you
  can discard log entries before N.
- Explain index arithmetic: log is 1-indexed (index 0 = sentinel/empty). `lastLogIndex` =
  length of the log (or length - 1 if 0-indexed). Be precise about off-by-one.
- Use `net.Listen`/`net.Dial` and `bufio.Scanner` (from M7/M8) — ex05 is a TCP gateway
  for a single replica.

You do NOT need to know: Raft, leader election, or log replication across nodes. M12 handles
that. Here you build the log mechanics for ONE node.

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
