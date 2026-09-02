# M10 · Commit & Consensus

A distributed system only works if its members can agree: on a transaction's
fate, on who leads, on which decisions reached a quorum. This module turns
the ordering ideas from M9 into the machinery of agreement.

Milestones:

- ex01 · The coordinator — the two-phase state machine that drives commit.
- ex02 · Voting — how one participant turns a prepare into a commit/abort
  verdict, and why an abort is final.
- ex03 · Leader election — breaking leader ties with a logical-time stamp.
- ex04 · Quorum — counting votes to decide commit versus abort.
- ex05 · The transaction — a real TCP gateway that runs two-phase commit
  over the wire.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.

---

## So what? (interview / portfolio)

Agreement protocols are the load-bearing wall of every distributed database and
transaction manager. Knowing *why* an abort is final, *how* a coordinator drives a
two-phase commit, and *what* a quorum is lets you speak credibly about distributed
transaction semantics — the kind of answer that turns a "consensus" interview into a
conversation.

**Interview questions this module arms you for:**
- What are the two phases of 2PC, and what makes the coordinator's decision final?
- Why is an abort irreversible once one participant votes abort?
- How is a leader elected when votes tie, without a wall clock to break it?
- What does it mean to say a decision "reached a quorum"?

**Portfolio artifact:** M10-ex05 `transaction` — a TCP gateway that runs true two-phase
commit over the wire.