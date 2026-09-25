# M16 · Storage Engines

The module replaces the in-memory map with the write and recovery path used by
real databases and key-value engines such as RocksDB and Bigtable: a
Log-Structured Merge-tree. The thread runs from a sorted memtable to immutable
SSTables, a crash-safe WAL, and newest-run compaction, then joins those pieces in
a persistent TCP key-value store.

## The build

- **ex01 · Memtable** — keep buffered writes in lexicographic order with update,
  lookup, and copied iteration contracts; the ordered handoff for a disk flush.
- **ex02 · SSTable** — serialize sorted entries into an immutable file, then use
  its sparse index and binary search for boundary-safe lookups.
- **ex03 · WAL** — append each mutation before the volatile memtable and replay
  the log in write order after a crash.
- **ex04 · Compaction** — merge sorted runs, collapse duplicates with the newest
  run winning, and remove tombstoned keys.
- **ex05 · Persistent KV gateway** — **Gate**: combine the storage path in a
  Mini-Capstone TCP store whose `store.log` restores state across restarts.

## Rules

Nothing printed or compared may depend on wall-clock time. All exercises are
transcript-driven and byte-deterministic.

## Prerequisites

Before starting M16, you should be comfortable with everything from M0–M15, plus:

- Explain what an LSM-tree (Log-Structured Merge-tree) is:
  - Writes go to an in-memory buffer (memtable), sorted by key.
  - When the memtable is full, it is flushed to disk as a sorted file (SSTable).
  - Reads check the memtable first, then SSTables from newest to oldest.
  - Compaction merges overlapping SSTables to reduce read overhead.
- Explain what an SSTable is: a sorted, immutable file of key-value pairs. Keys are in
  sorted order, enabling binary search within a file.
- Explain binary search on sorted data: `sort.Search(n, func(i int) bool { return keys[i] >= target })`.
  You must understand the Go `sort.Search` contract.
- Explain the WAL (write-ahead log): a sequential append-only file that records every write
  BEFORE it goes to the memtable. Used for crash recovery.
- Explain compaction: merging multiple SSTables into one, discarding duplicate keys (keep
  the latest version). This is the LSM-tree's way of cleaning up over time.
- Use disk I/O in Go: `os.File`, `file.Write()`, `file.ReadAt()`, `file.Seek()` for
  reading/writing SSTable files.
- Use `encoding/binary.LittleEndian` for on-disk encoding (fixed-width integer keys).
  (M8 used BigEndian for network framing — here we use LittleEndian for local disk.)
- Use `sync.Mutex` for concurrent access to the memtable and active SSTable list.

You do NOT need to know: Bloom filters, level-based compaction, tiered compaction, or
RocksDB internals. Those are advanced storage engine topics outside this module.

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
