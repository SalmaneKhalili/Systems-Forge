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
