# M13 · Sharding

A cluster can't store everything on one machine forever. **Sharding** spreads
keys across shards and keeps the data balanced as the cluster grows and
shrinks. This module builds the placement mechanics: fixed slot assignment,
key-range sharding, consistent hashing, and rebalancing.

Milestones:

- ex01 · Slot — assigning each key to a fixed slot by hash.
- ex02 · Range — partitioning keys by contiguous ranges.
- ex03 · Ring — consistent hashing places keys on a ring with minimal moves on resize.
- ex04 · Rebalance — moving keys between shards to restore balance.
- ex05 · The shard gateway — a TCP gateway that routes keys to the right shard.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.

---

## Prerequisites

Before starting M13, you should be comfortable with everything from M0–M12, plus:

- Explain what sharding is: splitting a key space across multiple nodes so each node stores
  only a subset of the data.
- Explain the two sharding strategies:
  - **Fixed-slot:** keys are mapped to slots by `hash(key) % numSlots`, and slots are assigned
    to shards. Growing the cluster requires reassigning slots.
  - **Range:** keys are partitioned by value ranges (e.g., A-M → shard1, N-Z → shard2).
    Growing the cluster requires splitting ranges.
- Explain consistent hashing (ring): hash both keys AND nodes onto a ring. Each key is
  assigned to the next node clockwise. Adding/removing a node only affects its neighbors.
- Explain the boundary problem: adjacent keys (e.g., "foo" and "foo1") hash to different
  positions. Range sharding handles this; fixed-slot sharding does not.
- Explain a shard gateway: a single entry point that routes each request to the correct
  shard based on the key.
- Use `net.Listen`/`net.Dial` and `bufio.Scanner` (from M7/M8) — ex05 is a TCP gateway
  that routes to shard backends.

You do NOT need to know: rebalancing algorithms, virtual nodes, or distributed hash tables
(Kademlia). Those are covered conceptually here but not implemented.

---

## So what? (interview / portfolio)

At any scale beyond a single machine, "where does this key live?" is the first question
you must answer, and every cloud store and cache is a sharding scheme under the hood.
Consistent hashing in particular — minimal key movement on resize — is a fixture of
system-design interviews for distributed caches and databases.

**Interview questions this module arms you for:**
- Hash-slot vs. key-range sharding: what are the trade-offs?
- Why does consistent hashing move so few keys when the ring resizes?
- How do you rebalance keys across shards without pausing the cluster?
- How does a client know which shard holds a given key?

**Portfolio artifact:** M13-ex05 `shard gateway` — a TCP gateway that routes each key to
the correct shard.
