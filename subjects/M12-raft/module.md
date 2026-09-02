# M12 · Raft — the flagship capstone

Everything so far converges here. M9 gave ordering, M10 agreement, M11 a
single replicated log with a commit point. **Raft** is how a leader drives ONE
consistent log across a whole cluster — safely, even when leaders crash and
re-election happens. This module builds the core mechanics of Raft's
safety: the monotonic term, the up-to-date-log vote rule, the log-matching
prefix invariant, and majority leader election.

Milestones:

- ex01 · Term — the monotonically increasing election term; every node obeys it.
- ex02 · Vote — a node votes only for a candidate whose log is at least as
  up-to-date as its own.
- ex03 · Log-match — the AppendEntries consistency check (the matching-prefix
  property that keeps logs in lockstep).
- ex04 · Election — a leader is elected when it holds a majority of votes.
- ex05 · The gate node — a TCP gateway serving one raft node's state to a
  client.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.
