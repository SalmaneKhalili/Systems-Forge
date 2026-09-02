# M14 · Membership

A cluster must agree on WHO is in it. **Membership** is how a node learns of
its peers, tracks whether they are alive, and eventually drops a peer that has
gone silent. This module builds the failure-detection mechanics: heartbeats,
the suspect/failed lifecycle, gossip propagation, and eviction.

Milestones:

- ex01 · Heartbeat — a node is considered alive while heartbeats arrive.
- ex02 · Suspect — silence moves a node from alive to suspect before failure.
- ex03 · Gossip — membership facts spread node-to-node and converge.
- ex04 · Evict — a failed node is removed from the live set.
- ex05 · The membership gateway — a TCP gateway that reports and updates the live set.

Rule: nothing printed or compared may depend on wall-clock time. Where a
timeout is needed, it is delivered as an event/step — never read the wall
clock. All exercises are transcript-driven and byte-deterministic.

---

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
