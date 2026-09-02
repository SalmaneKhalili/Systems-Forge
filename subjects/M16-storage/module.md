# M16 · Storage Engines

Everything you've persisted so far used a trivial in-memory map. Real storage
servers — databases, key-value engines like RocksDB, Bigtable — are built on a
**Log-Structured Merge-tree (LSM)**: buffer writes in a sorted in-memory
**memtable**, flush them to immutable sorted **SSTables** on disk, log every
mutation to a crash-safe **write-ahead log**, and periodically **compact**
overlapping runs so reads stay fast and stale data is reclaimed.

Milestones:

- ex01 · Memtable — the sorted in-memory write buffer that keeps keys ordered.
- ex02 · SSTable — flush to an immutable on-disk file; look up via sparse index
  + binary search.
- ex03 · WAL — append every mutation durably, then replay it after a crash.
- ex04 · Compaction — merge sorted runs; the newest run wins, tombstones delete.
- ex05 · Persistent KV gateway — the Mini-Capstone: an LSM key-value store with
  a TCP line protocol, durable across restarts.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.

---

## So what? (interview / portfolio)

The LSM-tree is the storage architecture behind some of the most-used engines
in the industry, so this module makes your resume's "storage" claim concrete
rather than hand-wavy. ex05 `persistent-kv` is a genuinely runnable mini-database
you can demo: put/get over TCP, WAL-backed recovery, sorted iteration. That's a
portfolio item an infrastructure reviewer recognises instantly, and the memtable
→ SSTable → WAL → compaction chain is a classic systems design story to tell.

**Interview questions this module arms you for:**
- Walk through the write path of an LSM engine (memtable → WAL → SSTable).
- Why are SSTables immutable, and what problem does compaction solve?
- How does a write-ahead log make a volatile memtable crash-safe?
- LSM vs. B-tree: when would you pick which?

**Portfolio artifact:** M16-ex05 `persistent-kv` — a TCP persistent KV gateway
(LSM: memtable + SSTable + WAL + newest-wins compaction).