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
