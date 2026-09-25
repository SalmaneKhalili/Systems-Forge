# M14 · Membership

This module builds the failure-detection layer that decides who belongs to a
cluster and which peers are still alive. The thread runs from heartbeat-driven
liveness through suspicion, gossip convergence, and eviction, then combines those
states in a TCP membership gateway.

## The build

- **ex01 · Heartbeat** — keep a peer alive while heartbeats arrive and fail it
  after a deterministic countdown.
- **ex02 · Suspect** — put a two-threshold `alive → suspect → failed`
  lifecycle between silence and removal.
- **ex03 · Gossip** — merge versioned membership views until every node keeps
  the same freshest facts.
- **ex04 · Evict** — remove members whose logical heartbeat sequence exceeds the
  inclusive staleness bound.
- **ex05 · The membership gateway** — **Gate**: serve concurrent TCP clients that
  inspect and mutate one process-wide live set, including clean `SIGTERM`
  shutdown.

## Rules

Nothing printed or compared may depend on wall-clock time. Where a timeout is
needed, it is delivered as an event/step — never read the wall clock. All
exercises are transcript-driven and byte-deterministic.

## Prerequisites

Before starting M14, you should be comfortable with everything from M0–M13, plus:

- Explain what cluster membership is: a node's knowledge of which other nodes are in the
  cluster, and whether each is alive.
- Explain the heartbeat failure-detection model:
  - Each node periodically sends a heartbeat to every peer.
  - If no heartbeat is received from a peer within a timeout, the peer is marked `SUSPECT`.
  - If the suspect node does not refute within a window, it is declared `FAILED`.
- Explain state machines for membership: `ALIVE → SUSPECT → FAILED`. A node can refute
  a suspicion by sending a heartbeat while in `SUSPECT` state.
- Explain gossip protocol: instead of every node talking to every other node directly, nodes
  periodically share their membership state with a random subset of peers. State propagates
  epidemically.
- Explain eviction: removing a `FAILED` node from the membership set. The node must be
  re-added manually (or automatically by an external controller) to rejoin.
- Explain the tradeoff between fast detection and false positives: shorter timeouts detect
  failures faster but may falsely suspect a slow node.
- Use `net.Listen`/`net.Dial` and `bufio.Scanner` — ex05 is a TCP membership gateway.

You do NOT need to know: SWIM protocol details, phi-accrual failure detectors, or
gossip convergence proofs. Those are advanced topics outside this module.

## So what? (interview / portfolio)

"Who is in my cluster, and are they actually alive?" is the subtle failure-detection
problem every distributed system must solve, and the suspect→failed lifecycle plus gossip
convergence are its canonical answers. This maps directly to how real systems (e.g. gossip
protocols) avoid split-brain and stale peers — a strong signal for a platform/infra role.

**Interview questions this module arms you for:**
- How does a heartbeat turn into a suspicion, and why the two-step don't-yell-about-it?
- Why use gossip instead of a central authority to spread membership?
- How do you evict a failed peer without unsafe races on the live set?
- What keeps a small cluster from wrongly declaring healthy peers suspect?

**Portfolio artifact:** M14-ex05 `membership gateway` — a TCP gateway that reports and
updates the live set via heartbeat/gossip/eviction.
