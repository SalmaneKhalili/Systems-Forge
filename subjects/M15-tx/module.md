# M15 · Transactions & Chaos

Replicated logs (M11) and Raft (M12) give you **durability of a single write**.
This module is about the harder promise: **making a set of writes behave as one
atomic unit**, and keeping that guarantee honest even when nodes are failing
around you. Transactions stage writes, commit atomically, and roll back cleanly;
a write-ahead log makes crash recovery deterministic; conflict detection stops
silent lost updates; and a chaos drill proves a cluster stays correct as nodes
drop.

Milestones:

- ex01 · Transaction — writes inside a txn are **staged** and atomic-commit.
- ex02 · Replay — a write-ahead log recovers state deterministically,
  last-writer-wins.
- ex03 · Conflict — detect overlapping writes; never silently lose an update.
- ex04 · Drill — a chaos drill: ops only succeed while their node is up.
- ex05 · txstore gateway — the transactional store (ex01–ex04) over a TCP
  line protocol.
- ex06 · durable txstore — the Micro-App: transactions fused with M16 **WAL
  durability** into one artifact.

Rule: nothing printed or compared may depend on wall-clock time. All
exercises are transcript-driven and byte-deterministic.

---

## Prerequisites

Before starting M15, you should be comfortable with everything from M0–M14, plus:

- Explain what a staged write is: a multi-step operation where intermediate results are
  buffered before being committed atomically.
- Explain the two-phase pattern (similar to M10's 2PC but for storage):
  - **Phase 1 (stage):** write each key's new value to a buffer or log.
  - **Phase 2 (commit):** apply all staged values atomically.
  - **Abort:** discard the buffer without applying.
- Explain write-ahead logging (WAL): write the intended change to a log file BEFORE applying
  it to the in-memory state. If you crash mid-write, replay the log on recovery.
- Explain conflict detection: if two concurrent transactions write to the same key, one must
  abort. Timestamp-based ordering (last-writer-wins) is the simplest strategy.
- Explain the chaos drill: inject faults (node crash, network partition) during a multi-key
  write and verify the system recovers to a consistent state.
- Use file I/O in Go: `os.WriteFile(path, data, 0644)` and `os.ReadFile(path)` for WAL files.
- Use `sync.Mutex` for concurrent access to shared state (the buffer and WAL).
- Use `net.Listen`/`net.Dial` and `bufio.Scanner` — ex05 is a TCP gateway for a
  single-node transactional KV store, ex06 is a multi-node cluster version.

You do NOT need to know: MVCC, 2PL (two-phase locking), or ARIES recovery. Those are
database internals topics outside this module.

---

## So what? (interview / portfolio)

Transactions + durability + conflict handling is exactly what a database or
distributed service promises, and interviewers probe it constantly. Your
strongest talking point is ex06 `durable txstore`: a single artifact that stages
atomic writes *and* appends them to a WAL, so you can honestly say you have
built "the storage half of a database" (isolation + durability) and can explain
the trade-off vs. redo/undo logs, 2PC, and optimistic concurrency.

**Interview questions this module arms you for:**
- What does it mean for a write to be atomic, and how does a rollback work?
- Why is a WAL the canonical durability mechanism, and what does replay do?
- What is a lost update, and how do optimistic conflict detection solve it?
- How is "availability" measured under node failures (your chaos drill)?

**Portfolio artifact:** M15-ex06 `durable txstore` — a transactional KV gateway
that commits atomically and durably via an observable write-ahead log.