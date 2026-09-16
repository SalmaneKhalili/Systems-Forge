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
