# M12 · Raft — the flagship capstone

Everything so far converges here: M9 supplies ordering, M10 agreement and M11
a replicated log with a commit point. Raft drives one consistent log across a
cluster through a monotonic term, an up-to-date vote rule, a matching-prefix
check and majority election, ending with a real three-node network.

## The build

- **ex01 · Term** — the monotonically increasing election term every node
  obeys; the version rule that gates the vote rule in ex02.
- **ex02 · Vote** — grant a vote only when the candidate's log is at least as
  up-to-date as the node's own; the election restriction that keeps a stale
  candidate from overtaking a fresh one.
- **ex03 · Log-match** — enforce AppendEntries' matching-prefix check; the
  consistency rule that keeps follower logs in lockstep.
- **ex04 · Election** — elect a leader only after a majority of votes; the
  quorum rule that completes one node's safety logic.
- **ex05 · The gate node (Gate)** — a TCP gateway serving one Raft node's
  term and committed log; the complete single-node contract handed to the
  cluster.
- **ex06 · The cluster** — a real three-node Raft cluster whose peers exchange
  RequestVote/AppendEntries over TCP and replicate one log.

## Rules

Nothing printed or compared may depend on wall-clock time. All exercises are
transcript-driven and byte-deterministic.

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
