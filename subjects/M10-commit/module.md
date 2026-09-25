# M10 · Commit & Consensus

This module turns M9's ordering rules into agreement: a transaction's fate,
one leader and the vote count that makes a decision durable. The exercises move
from coordinator state and participant promises to quorums, then compose those
pieces in a two-phase commit gateway.

## The build

- **ex01 · The coordinator** — a two-phase state machine that prepares every
  participant, chooses one decision and drives the final commit or abort.
- **ex02 · Voting** — the participant promise behind that decision: one
  environment yields a commit/abort verdict, and an abort is final.
- **ex03 · Leader election** — break competing leader candidates with a logical
  election stamp, then its id, producing the deterministic leader the remaining
  decisions rely on.
- **ex04 · Quorum** — count votes for a strict majority so opposing decisions
  cannot both reach a quorum.
- **ex05 · The transaction** — a real TCP gateway that runs two-phase commit
  over the wire.

## Rules

Nothing printed or compared may depend on wall-clock time. All exercises are
transcript-driven and byte-deterministic.

## Prerequisites

Before starting M10, you should be comfortable with everything from M0–M9, plus:

- Explain what a distributed transaction is: a set of writes across multiple nodes that must
  either ALL succeed or ALL fail (atomicity across machines).
- Explain the two-phase commit (2PC) protocol:
  - **Phase 1 (prepare):** coordinator asks all participants "can you commit?" — each replies
    YES or NO.
  - **Phase 2 (commit/abort):** coordinator tells all participants the outcome. If ANY
    participant voted NO or timed out, the coordinator aborts.
- Explain the coordinator's state machine: `INIT → PREPARING → DECIDING → COMMITTED/ABORTED`.
- Explain what a quorum is: in voting, a quorum = ⌊N/2⌋ + 1 votes (majority).
- Explain why a coordinator crash during phase 2 is dangerous: participants are locked in
  PREPARED state and cannot proceed without the coordinator's decision.
- Use `net.Listen`/`net.Dial` and `bufio.Scanner` (from M7/M8) — ex05 is a TCP gateway
  for a coordinator.

You do NOT need to know: Raft, Paxos, or leader election. Those come in M11–M12.

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
