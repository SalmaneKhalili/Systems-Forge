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