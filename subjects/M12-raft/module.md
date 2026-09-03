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
- ex06 · The cluster — a real 3-node raft cluster: peers exchange
  RequestVote/AppendEntries over real TCP and replicate one log.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.

---

## Prerequisites

Before starting M12, you should be comfortable with everything from M0–M11, plus:

- Explain Raft leader election:
  - A node starts an election by incrementing its term and sending `RequestVote` RPCs.
  - A node votes for at most one candidate per term (first-come-first-served).
  - A candidate wins if it receives votes from a majority (⌊N/2⌋ + 1).
  - If no winner, a new election starts with a higher term.
- Explain Raft log replication:
  - The leader appends entries to its log, then sends `AppendEntries` RPCs to followers.
  - Each follower appends the entry and acknowledges.
  - Once a majority has persisted the entry, the leader commits it (advances commit index).
  - Followers learn the commit index from the next `AppendEntries` heartbeat.
- Explain log matching: if two logs have the same index and term, all entries before that
  index are identical. This is how the leader detects missing/conflicting entries.
- Explain what `term` is: a monotonically increasing counter. Each term has at most one
  leader. A higher term always trumps a lower term.
- Explain the TCP peer-to-peer protocol: nodes talk to each other directly (not through
  a central switch). Each node has its own listener; you connect to peers by address.
- Use goroutines for concurrent TCP connections (ex05/06 each require handling multiple
  peers simultaneously).

You do NOT need to know: log compaction/snapshotting, cluster membership changes, or
linearizability proofs. Those are advanced topics outside this module.

---

## So what? (interview / portfolio)

Raft is the *flagship capstone* and the single most askable algorithm in distributed
systems interviews. Being able to explain — and in ex06 actually *run* — a real 3-node
raft cluster (election via majority, step-down on a higher term, append-entries
replication) is a headline portfolio claim no infrastructure candidate should be without.

**Interview questions this module arms you for:**
- How does a raft term stay monotonic, and what "steps down" a leader?
- Why must a voter only back a candidate with an at-least-as-up-to-date log?
- What invariant does the matching-prefix check in AppendEntries preserve?
- How many votes elect a leader of N nodes, and why a strict majority?

**Portfolio artifact:** M12-ex06 `raft cluster` — a real 3-node raft cluster whose three
peers agree on one replicated log over real TCP sockets.
