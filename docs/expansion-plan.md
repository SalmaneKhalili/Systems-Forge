# Systems-Forge — Complete Exercise, Assignment, and Gate Expansion Plan

> **Revision 2 (2026-09-17).** Rebased onto the actual curriculum after the review in
> [`docs/expansion-review.md`](expansion-review.md). Key changes: corrected the M0 anchor
> list (no Git exercise exists — it is a planned addition), added the omitted `ex06-cluster`
> (M12) and `ex06-durabletxstore` (M15) anchors, renumbered M12/M15 expansions past the
> existing `ex06` exercises, and added feasibility/scope notes (M6/M7 assignments duplicate
> gates; M9/M11/M14 must respect module fences; M10's gateway is a static votes lookup).

## 0. Design Rule

The existing curriculum is the canonical foundation.

Do **not** delete or casually rewrite existing exercises.

Instead, expand each module from:

```text
Exercise
Exercise
Exercise
Exercise
Gate
```

into:

```text
Foundation exercises
        ↓
Observation / reasoning
        ↓
Variation
        ↓
Failure / debugging
        ↓
Integration exercises
        ↓
Assignment
        ↓
Concept review
        ↓
Gate
```

The existing exercises remain important anchors.

The new exercises exist to make the concepts digestible rather than allowing a learner to encounter a mechanism once and immediately move on.

---

# M0 — Tools of the Trade

## Existing curriculum

M0 currently establishes:

* Make (`all`/`fclean`/`re` discipline)
* shell scripting (`set -euo pipefail` hygiene)
* Python tooling (filesystem/invariant checker)
* sanitizers (AddressSanitizer / use-after-free)
* development hygiene (the mini-forge gate — a miniature of the real grader)
* Git — **not yet in the curriculum**; planned as a new addition below

The current exercises are (all explicitly defined in the curriculum):

1. `ex01-makefile` — Makefile discipline with `all/fclean/re`
2. `ex02-shell` — shell hygiene under `set -euo pipefail`
3. `ex03-checker` — Python filesystem/invariant checker
4. `ex04-sanitize` — AddressSanitizer / use-after-free
5. `ex05-gate` — **Gate: mini-forge** (a real mini grader graded via `scenario`)

## Expanded exercises

### M0-E01 — Makefile Anatomy

Keep the existing Makefile exercise.

Focus:

* targets
* prerequisites
* recipes
* dependency graphs
* `all`
* `fclean`
* `re`
* incremental builds

### M0-E02 — Make Dependency Graph

Given:

```text
main.c
foo.c
bar.c
foo.h
bar.h
```

construct a Makefile that rebuilds only the affected object files.

Concepts:

* timestamps
* prerequisites
* incremental compilation
* dependency correctness

### M0-E03 — Shell Exit Status

Build a script that:

* executes multiple commands;
* stops appropriately on failure;
* propagates meaningful status codes.

Focus:

* `$?`
* `set -e`
* explicit conditionals
* exit status propagation

### M0-E04 — Shell Argument Discipline

Extend shell scripting to:

* positional arguments
* quoting
* empty arguments
* filenames containing spaces
* missing arguments

### M0-E05 — Filesystem Checker

Keep/expand the Python checker.

Require:

* file existence
* directory existence
* permissions
* expected contents
* failure reporting

### M0-E06 — Git State Reconstruction

Given a broken repository history, determine:

* current branch
* commits
* changed files
* merge/rebase state
* recovery path

### M0-E07 — Sanitizer Laboratory

Deliberately diagnose:

* use-after-free
* out-of-bounds access
* leak
* undefined behavior

Require the learner to identify:

```text
symptom
→ sanitizer report
→ offending operation
→ root cause
→ repair
```

### M0-E08 — Build Failure Diagnosis

Give the learner several broken projects:

* missing dependency
* stale object
* incorrect target
* compiler warning
* sanitizer failure

They must determine the failure without being told which subsystem is broken.

## Assignment — M0-A01: Reproducible Project Harness

Build a small project with:

```text
src/
include/
tests/
scripts/
Makefile
README
```

Requirements:

* reproducible build
* `all`
* `fclean`
* `re`
* automated test script
* sanitizer target
* meaningful exit statuses
* Git history following the project's conventions

The learner must be able to clone the project into a clean directory and reproduce the result.

> **Dependency note:** the "Git history following the project's conventions" requirement has
> no prerequisite in the current curriculum — make **M0-E06 (Git State Reconstruction)** the
> required lead-in to this assignment (or soften this bullet). Everything else in the
> assignment has an anchor in ex01–ex05.

## Gate

Existing tooling/hygiene gate.

---

# M1 — C Foundations

The existing module already contains:

1. `ft_strlen`
2. `ft_strlcpy`
3. `ft_memset`
4. `ft_strcmp`
5. `ft_strdup` — gate

These exercises already deliberately test embedded NULs, size-zero behavior, return values, signed-byte behavior, allocation, and memory safety.

## Expanded exercises

### M1-E06 — Pointer Traversal

Implement:

* `ft_memcpy`
* `ft_memmove`

Focus on pointer arithmetic and overlap.

### M1-E07 — Allocation Contracts

Implement a family of functions where the learner must reason about:

* ownership
* NULL
* allocation size
* failure
* freeing responsibility

### M1-E08 — String Construction

Build:

* concatenation
* substring
* duplication
* bounded construction

### M1-E09 — Memory Failure

Modify implementations so allocation failure occurs deterministically.

Learner must ensure:

* no leak
* no invalid access
* correct failure propagation

### M1-E10 — Multi-Function Library

Combine previous functions into a small library and diagnose cross-function bugs.

## Assignment — M1-A01: Mini libc

Build a small coherent libc-style library containing:

* string functions
* memory functions
* allocation helpers

Requirements:

* strict compiler flags
* Makefile
* tests
* ASan
* UBSan
* no leaks
* documented ownership contracts

## Gate

The existing `ft_strdup` gate (ex05) **stays as-is** — it is exactly the memory-safe
allocate/copy/terminate/no-leak test the module already grades under ASan. The mini-libc
assignment builds the library the gate samples from; no fresh gate API is needed.

---

# M2 — Processes

The existing M2 already contains:

1. fork/wait
2. fork/exec
3. pipes
4. signals
5. pipeline gate

The current curriculum explicitly defines the progression this way.

## Expanded exercises

### M2-E06 — Fork State

Observe what changes and what remains the same across:

```text
fork()
```

Require reasoning about:

* PID
* return value
* variables
* address spaces
* inherited descriptors

### M2-E07 — Multiple Children

Create several children and:

* identify them
* wait for each
* collect statuses
* prevent zombies

### M2-E08 — Exec Argument Vectors

Experiment with:

```text
argv
argv[0]
argv[1...]
PATH
execvp
execlp
```

### M2-E09 — Descriptor Inheritance

Create descriptors before `fork()` and determine which process can access them.

### M2-E10 — Pipe EOF Laboratory

Intentionally create:

* correct closure
* missing writer closure
* missing reader closure
* early child exit

Explain each resulting behavior.

### M2-E11 — dup2 Redirection

Redirect:

```text
stdout → file
stdout → pipe
stdin ← file
stdin ← pipe
```

### M2-E12 — Signal State Machine

Practice:

* default disposition
* handler
* ignored signal
* `raise`
* `sig_atomic_t`

## Assignment — M2-A01: Mini Process Runner

Build:

```text
runner command [args...]
```

It must:

* fork
* exec
* report exec failure
* wait
* decode status
* propagate child status

### Assignment extension

Add:

```text
runner command1 | command2
```

using:

* `pipe`
* `fork`
* `dup2`
* `exec`
* `waitpid`

The learner must document every descriptor opened and closed.

## Gate

Existing pipeline gate.

---

# M3 — Memory

Current anchors are:

1. bump arena
2. copy-on-write virtual-memory experiment
3. fixed-slot pool
4. mmap allocator
5. OOM-safe growable buffer gate

The current module explicitly follows the idea of building allocation mechanisms in progressively richer forms.

## Expanded exercises

### M3-E06 — Alignment Laboratory

Implement allocations satisfying:

* alignment
* size
* boundaries

### M3-E07 — Fragmentation Experiment

Construct allocation/free patterns and measure internal/external fragmentation conceptually.

### M3-E08 — Ownership Contracts

Implement structures with explicit:

```text
owner
borrower
lifetime
release
```

### M3-E09 — Virtual Memory Mapping

Explore:

* pages
* protection
* private/shared mappings
* mapping lifetime

### M3-E10 — Allocator Failure Injection

Force:

* exhaustion
* invalid free
* double free
* oversized allocation

and require deterministic behavior.

### M3-E11 — Reallocation

Implement bounded growth while preserving:

* existing data
* alignment
* failure atomicity

## Assignment — M3-A01: Bounded Allocator

Build a bounded allocator using a controlled memory region.

Requirements:

* allocation
* free
* alignment
* reuse
* exhaustion
* invalid operation handling
* accounting

The assignment should explicitly test:

```text
correctness
+
memory safety
+
failure semantics
```

## Gate

Existing growable-buffer/OOM gate.

---

# M4 — Concurrency

Current exercises are:

1. join & exit
2. mutex
3. atomic counter
4. readers/writers
5. bounded queue gate

The existing module already uses ThreadSanitizer and deliberately tests races and FIFO correctness.

## Expanded exercises

### M4-E06 — Race Reconstruction

Given a racy program, identify:

```text
shared state
→ conflicting accesses
→ missing synchronization
```

### M4-E07 — Mutex Ownership

Explore:

* lock/unlock ownership
* critical sections
* deadlock from incorrect locking

### M4-E08 — Condition Variables

Build a minimal producer/consumer synchronization exercise.

### M4-E09 — Lost Wakeup

Intentionally implement incorrect condition-variable logic and diagnose it.

### M4-E10 — Deadlock

Create a two-lock deadlock and repair it.

### M4-E11 — Atomic Memory Operations

Compare:

* mutex
* atomic increment
* unsynchronized increment

### M4-E12 — Queue Lifecycle

Handle:

* producer termination
* consumer termination
* queue closure
* empty/full states

## Assignment — M4-A01: Concurrent Work Queue

Build a bounded worker queue supporting:

* multiple producers
* multiple consumers
* FIFO ordering
* blocking push/pop
* graceful closure
* no data races
* no deadlocks

## Gate

Existing bounded queue.

---

# M5 — Files & I/O

Current anchors:

1. descriptors
2. `read` vs stdio
3. `dup2` redirection
4. tee
5. log parser gate

The existing module already targets descriptor semantics, buffering, redirection and a log-processing artifact.

## Expanded exercises

### M5-E06 — File Descriptor Lifecycle

Practice:

```text
open
→ read/write
→ dup
→ dup2
→ close
```

### M5-E07 — Partial I/O

Handle:

* partial reads
* partial writes
* EOF
* errors

### M5-E08 — Buffering

Compare:

* `read/write`
* stdio buffering
* flushing
* output ordering

### M5-E09 — Redirection Utility

Implement:

```text
input file → stdin
stdout → output file
stderr → error file
```

### M5-E10 — Streaming Parser

Parse arbitrarily sized input without assuming it fits in one read.

### M5-E11 — Log Rotation

Detect a changing log file and reopen appropriately.

## Assignment — M5-A01: Log Processing Utility

Build a streaming log processor that:

* reads from stdin or a file
* parses records
* filters records
* writes results
* handles arbitrary input sizes
* correctly handles partial reads
* reports malformed records

## Gate

Existing log-parser gate.

---

# M6 — Networking

Current anchors:

1. echo
2. line chat
3. HTTP-ish server
4. read timeout
5. stateful networking gate

The curriculum specifically tests accept loops, protocol responses, timeouts and state across connections.

## Expanded exercises

### M6-E06 — Socket Lifecycle

Understand:

```text
socket
→ bind
→ listen
→ accept
→ read/write
→ close
```

### M6-E07 — Client Lifecycle

Build a TCP client with:

* connect
* send
* receive
* close

### M6-E08 — Protocol Framing

Handle:

* delimiters
* partial messages
* multiple messages in one read

### M6-E09 — Connection State

Maintain state independently for multiple clients.

### M6-E10 — Timeouts

Handle:

* idle clients
* read deadlines
* timeout recovery

### M6-E11 — Concurrent Connections

Serve multiple connections without allowing one client to permanently block another.

### M6-E12 — Protocol Failure

Handle:

* malformed request
* premature disconnect
* oversized message
* unexpected input

## Assignment — M6-A01: Stateful TCP Service

Build a small line-oriented TCP service with:

* multiple clients
* per-client state
* protocol framing
* malformed-input handling
* timeout handling
* deterministic responses

> **Rebase note:** the core of this (greeting on accept, per-connection state, line framing,
> `PING/ECHO/BAD`, accept loop spanning connections) **already exists as M6-ex05 gate**.
> Reshape the assignment to add what the gate does not test: malformed-input handling,
> oversized messages, and deterministic failure replies. Also note the `net` grader drives
> connections **sequentially**, so *simultaneous* per-client state cannot be graded today —
> that variant needs a `scenario`/`fault` driver or a new concurrent-net runner shape.

## Gate

Existing stateful TCP gate.

---

# M7 — Resilience

Current anchors:

1. retries/backoff
2. circuit breaker
3. health checks
4. graceful shutdown
5. supervisor gate

The current module already combines Python and Go and uses fake clocks where deterministic timing is required.

## Expanded exercises

### M7-E06 — Failure Classification

Distinguish:

* retryable
* non-retryable
* transient
* permanent

### M7-E07 — Backoff Mathematics

Implement:

* fixed backoff
* exponential backoff
* capped backoff

### M7-E08 — Retry Budgets

Prevent infinite retry storms.

### M7-E09 — Health State

Model:

```text
healthy
→ degraded
→ unhealthy
→ recovering
```

### M7-E10 — Shutdown Coordination

Coordinate:

```text
signal
→ stop accepting work
→ finish current work
→ cleanup
→ exit
```

### M7-E11 — Restart Semantics

Restart failed workers while preserving completed progress.

## Assignment — M7-A01: Mini Supervisor

Build a supervisor that manages several workers.

It must:

* start workers
* detect failure
* retry
* restart
* preserve progress
* enforce retry budgets
* perform graceful shutdown

> **Rebase note:** every bullet here is already tested by **M7-ex05 gate (mini supervisor)** —
> three workers, watchdog restart of a crashed worker with `retry ok`, progress preserved
> across restart, 2-attempt budget (`giving up`), graceful shutdown. Reshape the assignment
> to go beyond the gate: e.g. workers as real subprocesses (not in-process tasks), a crash
> *between* successful tasks, or restart-while-preserving-progress under a supervisor that
> must also re-verify a restarted peer via the `net` `restart` step.

## Gate

Existing supervisor gate.

---

# M8 — Messaging + Fault Injection

Current anchors:

1. framing
2. ordered delivery
3. transparent relay
4. fault injection
5. switch gate

The current module deliberately introduces the fault-injection switch that later modules reuse.

## Expanded exercises

### M8-E06 — Frame Corruption

Handle:

* truncated frame
* invalid length
* oversized frame
* malformed payload

### M8-E07 — Message Identity

Track:

```text
message ID
partition
sequence
sender
```

### M8-E08 — Duplicate Handling

Implement deduplication.

### M8-E09 — Missing Messages

Detect gaps without silently reordering.

### M8-E10 — Fault Composition

Test combinations:

```text
drop + duplicate
delay + duplicate
drop + hold
```

### M8-E11 — Delivery Semantics

Compare:

* at-most-once
* at-least-once
* effectively-once application behavior

## Assignment — M8-A01: Reliable Messaging Layer

Build a small message transport layer supporting:

* framing
* message identity
* ordering
* duplicate detection
* gap handling
* injected faults

The learner must demonstrate behavior under:

```text
normal
drop
duplicate
hold
combined faults
```

## Gate

Existing switch gate.

---

# M9 — Time & Ordering

Current exercises:

1. Lamport clock
2. vector clock
3. causality
4. total order
5. TCP sequencer

These are explicitly the current M9 topics.

## Expanded exercises

### M9-E06 — Event Trace Reconstruction

Given a distributed trace, reconstruct causal relationships.

### M9-E07 — Lamport Counter Failures

Break:

* receive merge
* local increment
* message timestamping

### M9-E08 — Vector Clock Comparison

Implement:

```text
before
after
concurrent
equal
```

### M9-E09 — Causal Broadcast

Prevent delivery of an event before its causal predecessors.

### M9-E10 — Total Ordering

Construct deterministic ordering from:

```text
logical timestamp
+
node ID
+
event ID
```

### M9-E11 — Sequencer State

Maintain sequence state across connections.

## Assignment — M9-A01: Distributed Event Ordering Service

Build a deterministic event-ordering service.

It should accept events from multiple simulated nodes and produce:

* causal relationships
* logical timestamps
* deterministic total order

No wall-clock dependence.

> **Scope note:** keep the "multiple simulated nodes" **in-process** (fed over the M9-ex05
> sequencer's wire protocol or an injected trace). A genuinely multi-node ordering service
> would require consensus, which M9 explicitly fences off ("You do NOT need to know: Raft,
> 2PC, Paxos … those come in M10–M12"). Defer the multi-node variant to a cumulative
> assignment instead.

## Gate

Existing TCP sequencer.

---

# M10 — Commit & Consensus

Current anchors:

1. coordinator
2. vote
3. leader election
4. quorum
5. transaction gateway — a TCP gateway that decides `COMMIT`/`ABORT` from a provided
   `votes.txt` (a **static prepare-allowed lookup**, not yet wire-level 2PC)

The existing module explicitly builds toward a TCP transaction gateway deciding from votes.

## Expanded exercises

### M10-E06 — Quorum Mathematics

Determine whether a proposed decision has sufficient votes.

### M10-E07 — Election Failures

Handle:

* tied votes
* missing voters
* stale candidates
* duplicate votes

### M10-E08 — Coordinator Failure

Model coordinator disappearance during a transaction.

### M10-E09 — Prepare/Commit State Machine

Explicitly model:

```text
INIT
→ PREPARED
→ COMMITTED
```

and abort paths.

### M10-E10 — Decision Safety

Prevent conflicting decisions.

### M10-E11 — Recovery

Recover transaction state after restart.

## Assignment — M10-A01: Deterministic Transaction Coordinator

Build a coordinator managing a small set of participants.

Requirements:

* prepare
* vote
* quorum/decision logic
* commit
* abort
* recovery
* deterministic state transitions

> **Scope note:** participants should be spawned **in-process over real sockets** (this is
> the genuinely missing wire-2PC piece; today's gate only looks up a votes file). The
> **`recovery`** requirement presumes persistent state — pull that from M15's WAL work or
> split the assignment: wire-2PC here, coordinator crash-recovery as part of M15-A01.

## Gate

Existing transaction gateway.

---

# M11 — Replicated Log & Consistency

Current anchors:

1. append
2. committed prefix
3. index arithmetic
4. snapshot
5. TCP replica

The existing module specifically focuses on append logs, commit index, index math, snapshots and serving only committed prefixes.

## Expanded exercises

### M11-E06 — Log Invariants

Enforce:

```text
index monotonicity
term association
prefix consistency
```

### M11-E07 — Commit Advancement

Calculate the largest safely committed index.

### M11-E08 — Conflicting Entries

Repair divergent prefixes.

> **Fence note:** prefix-matching/divergence repair currently lives in **M12-ex03-logmatch**
> (and M11's module.md explicitly says "You do NOT need to know: Raft … M12 handles that").
> Either teach the prefix-matching math here as pure log arithmetic (indices + terms only),
> or move conflict repair to M12 and keep this slot for snapshot-boundary work.

### M11-E09 — Snapshot Boundaries

Translate:

```text
logical index
↔ snapshot-relative index
```

### M11-E10 — Recovery

Reconstruct committed state after restart.

### M11-E11 — Replica Read Semantics

Ensure uncommitted entries are never exposed.

## Assignment — M11-A01: Replicated Log Node

Build a standalone replicated-log state machine supporting:

* append
* conflict resolution
* commit tracking
* snapshotting
* recovery
* committed-prefix reads

## Gate

Existing TCP replica gate.

---

# M12 — Raft

M12 is the flagship capstone.

Current exercises establish:

1. monotonic term (`ex01-term`)
2. voting/up-to-date log rule (`ex02-vote`)
3. AppendEntries prefix matching (`ex03-logmatch`)
4. election/quorum mechanics (`ex04-election`)
5. TCP Raft gate node (`ex05-gatenode`)
6. **`ex06-cluster` — a real 3-node Raft cluster over TCP** (elect/propose/read/log/term
   commands; RequestVote/AppendEntries on real sockets; strict majority; higher-term
   steps down)

The existing curriculum explicitly describes M12 as the flagship Raft capstone.

> **Rebase note:** the plan's original anchor list **omitted ex06-cluster** — the single
> most integrative exercise in the curriculum. It already delivers leader election, voting,
> log replication, commit advancement, and state-machine reads on a real 3-node cluster.
> The M12 assignment below must therefore be scoped to what ex06 does **not** cover:
> crash/restart recovery and fault injection (see assignment note).

Here I would **not merely add lots of tiny exercises**.

Raft requires substantially more integration, and much of that integration already exists.

## Expanded exercises

### M12-E07 — Raft State Machine

Implement:

```text
Follower
Candidate
Leader
```

transitions.

### M12-E08 — Election Safety

Test:

* stale terms
* duplicate votes
* outdated logs
* split votes

### M12-E09 — Leader Replication

Model:

```text
AppendEntries
nextIndex
matchIndex
```

### M12-E10 — Commit Rule

Determine when an entry is safely committed.

### M12-E11 — Crash Recovery

Recover:

* term
* votedFor
* log
* committed state

### M12-E12 — Network Partition Reasoning

Given partition scenarios, determine which side can safely commit.

### M12-E13 — Client Semantics

Handle client requests through a leader while maintaining committed-prefix semantics.

## Assignment — M12-A01: Raft Node

This should be a **major assignment**, not a normal exercise.

Build a functioning Raft node supporting:

* leader election
* terms
* voting
* log replication
* commit advancement
* state machine application
* crash/restart recovery

The fault switch from M8 must be usable against it.

> **Rebase note:** the running node **already exists as M12-ex06-cluster**. Scope this
> assignment to the delta: (1) **crash/restart recovery** of persistent state
> (term/votedFor/log/commit) — gradeable today via the `net` `restart` step, which is
> implemented but currently used by **zero** exercises; and (2) fault tolerance driven by a
> `scenario`/`fault` runner (partition/delay/drop/node kill) with deterministic transcripts.
> Reuse the existing cluster binary and its command protocol rather than building a new one.

## Capstone Extension

Introduce:

* partition
* delayed messages
* dropped messages
* node restart
* leader failure

The learner must demonstrate preservation of Raft safety properties.

> **Feasibility note:** the M8 switch is **not yet wired into any later module** ("later
> exams run the learner's switch" is documented intent, not a shipped capability), and
> `TARGETFAULTS` is not auto-injected into `net` runs. Decide and deliver switch-reuse (or
> an internal fault plane in the cluster binary) before making "the fault switch from M8
> must be usable against it" a hard requirement. Recommended vehicle: the `scenario`/`fault`
> runner (currently barely used: 1 exercise uses `scenario`, 0 use `fault`).

## Gate

Existing Raft TCP gate becomes the formal safety assessment.

---

# M13 — Sharding

Current curriculum:

1. slot assignment
2. key-range routing
3. consistent-hash ring
4. rebalancing
5. shard gateway

These are explicitly the existing M13 artifacts.

## Expanded exercises

### M13-E06 — Key Distribution

Measure how keys map across shards.

### M13-E07 — Range Boundaries

Handle:

* inclusive/exclusive boundaries
* empty ranges
* boundary keys

### M13-E08 — Ring Edge Cases

Test:

* wraparound
* empty ring
* single node
* duplicate positions

### M13-E09 — Rebalancing

Move keys while minimizing movement.

### M13-E10 — Shard Ownership

Track:

```text
key
→ shard
→ owner
```

### M13-E11 — Routing During Movement

Define behavior while ownership changes.

## Assignment — M13-A01: Sharded KV Router

Build a deterministic routing layer supporting:

* hash/range assignment
* shard lookup
* node membership
* rebalancing
* routing during topology changes

## Gate

Existing shard gateway.

---

# M14 — Membership

Current curriculum focuses on:

1. heartbeat countdown
2. alive → suspect → failed
3. gossip merge
4. staleness eviction
5. membership gateway

These behaviors are explicitly documented.

## Expanded exercises

### M14-E06 — Failure Detector State Machine

Model:

```text
ALIVE
→ SUSPECT
→ FAILED
```

### M14-E07 — Heartbeat Accounting

Handle:

* heartbeat refresh
* missed heartbeat
* recovery

### M14-E08 — Gossip Merge

Ensure merge is:

* deterministic
* commutative
* idempotent where appropriate

### M14-E09 — Stale Membership

Evict stale entries correctly.

### M14-E10 — Conflicting Views

Merge two nodes with different membership knowledge.

### M14-E11 — Recovery

Reintroduce a previously failed member.

## Assignment — M14-A01: Gossip Membership Service

Build a membership subsystem with:

* node identity
* heartbeats
* failure detection
* gossip
* merge
* stale-entry removal
* recovery

> **Scope note:** the current module is **single-process / tick-modeled** (gossip is a pure
> commutative `Merge`; the gate tracks one in-process view). Keep this assignment
> in-process with the module's tick style, plus the genuinely missing behaviors — a real
> **down/restart** transition and **failed-member reintroduction** (today's `drop`→`dead`
> has no rejoin path). Defer *networked* gossip between spawned processes to a cumulative
> assignment unless a scenario driver is built for it.

## Gate

Existing membership gateway.

---

# M15 — Transactions + Chaos

Current curriculum combines:

* 2PC
* transactional store
* WAL
* last-writer-wins replay
* lost-update detection/resolution
* chaos testing
* transactional KV gateway

This is explicitly how M15 is currently described.

The concrete anchors are `ex01-transaction` … `ex05-txstore`, plus
**`ex06-durabletxstore`** — a durable transactional KV gateway that appends every committed
`set` to a WAL (`wal` command shows the commit record; `rollback` leaves no WAL trace).
*Correction: ex06 is a single-process gateway, not the "multi-node cluster version" that
module.md's prose claims; adopt it here as the WAL-durability anchor.*

## Expanded exercises

### M15-E07 — Transaction State Recovery

Recover:

```text
prepared
committed
aborted
```

states after restart.

### M15-E08 — WAL Replay

Handle:

* complete records
* incomplete final record
* duplicate records
* ordering

### M15-E09 — Lost Update

Construct and diagnose concurrent update races.

### M15-E10 — Conflict Resolution

Compare:

* last writer wins
* reject
* retry
* merge

### M15-E11 — Chaos Scenarios

Inject:

* node crash
* delayed response
* unavailable participant
* coordinator failure

### M15-E12 — Transactional Recovery

Combine WAL + transaction state + node failure.

## Assignment — M15-A01: Crash-Recoverable Transaction Store

Build a transactional KV service supporting:

* transactions
* WAL
* recovery
* conflict detection
* deterministic conflict resolution
* participant failure

Then run a prescribed chaos scenario suite.

> **Rebase note:** transactions + WAL + durability + ordering are **already built** in
> `ex06-durabletxstore`. What is genuinely missing here: **conflict detection/resolution on
> the wire** (ex03 is a pure-function lost-update resolver), **participant failure**, and a
> *real* crash-restart recovery — note that no exercise today grades restart-replay; the
> `net` `restart` step is the right tool. Import the M16-style WAL replay semantics only to
> the extent the module's forward-reference already assumes them.

## Gate

Existing transactional gateway under fault injection.

---

# M16 — Storage Engines

Current curriculum:

1. memtable
2. SSTable
3. WAL
4. compaction
5. persistent KV gateway

The current implementation specifically uses a sorted memtable, sparse SSTable index, append-only WAL, newest-wins compaction and tombstones.

## Expanded exercises

### M16-E06 — Memtable Ordering

Implement ordered insertion and lookup.

### M16-E07 — SSTable Layout

Construct:

```text
data region
+
sparse index
```

### M16-E08 — WAL Recovery

Replay records into memory.

### M16-E09 — Tombstones

Distinguish:

```text
missing
deleted
never existed
```

### M16-E10 — Compaction

Merge multiple SSTables while preserving newest-wins semantics.

### M16-E11 — Crash During Persistence

Determine what survives after:

* WAL write
* SSTable write
* compaction

### M16-E12 — Persistent Recovery

Restart the engine and reconstruct the correct state.

## Assignment — M16-A01: Persistent LSM KV Store

Build a complete small LSM-style storage engine:

```text
PUT
GET
DELETE
WAL
memtable
flush
SSTable
compaction
recovery
```

Then expose it through a TCP gateway.

## Gate

Existing persistent KV gateway.

---

# M17 — Observability

Current curriculum:

1. metrics
2. tracing
3. spans
4. summaries
5. telemetry gateway

The existing implementation includes concurrency-safe counters, span trees, span attributes, depth-first queries, and deterministic p50/p95 summaries.

## Expanded exercises

### M17-E06 — Counter Semantics

Define and implement:

* counter
* gauge
* snapshot

### M17-E07 — Labels

Handle dimensional metrics without corrupting identity.

### M17-E08 — Trace Trees

Build parent/child relationships.

### M17-E09 — Trace Queries

Support:

* find span
* maximum duration
* subtree traversal

### M17-E10 — Quantiles

Understand and implement:

* p50
* p95
* p99

### M17-E11 — Telemetry Correlation

Connect:

```text
request
→ trace
→ spans
→ metrics
```

### M17-E12 — Failure Observability

Ensure failures produce enough information to diagnose:

* latency
* errors
* request identity
* subsystem

## Assignment — M17-A01: Observable Service

Take a small TCP service and instrument it with:

* metrics
* traces
* span attributes
* latency summaries
* error counters

Then expose the telemetry through a deterministic interface.

## Gate

Existing telemetry gateway.

---

# Cross-Module Assignments

The individual assignments are not enough.

Every few modules, the curriculum should force previously learned concepts back into use.

## Cumulative Assignment 1 — Systems Utility

After M0–M3:

Build a Unix utility combining:

* C
* Make
* memory management
* process creation
* file descriptors

---

## Cumulative Assignment 2 — Concurrent Service

After M4–M7:

Build:

```text
TCP server
+
worker pool
+
timeouts
+
retry policy
+
graceful shutdown
+
supervision
```

This should deliberately require concepts from multiple modules.

---

## Cumulative Assignment 3 — Fault-Tolerant Service

After M8–M10:

Build a service involving:

```text
messaging
+
fault injection
+
logical ordering
+
quorum/commit
```

The learner should not be explicitly told which previously learned mechanisms they are expected to reuse.

---

## Cumulative Assignment 4 — Replicated Service

After M11–M12:

Build a small replicated state machine using the Raft concepts learned in M12.

This becomes the bridge into:

```text
sharding
membership
transactions
```

---

## Cumulative Assignment 5 — Distributed KV Platform

After M13–M17:

Combine:

```text
Raft
+
sharding
+
membership
+
transactions
+
storage engine
+
observability
```

The final artifact should resemble a miniature distributed database/platform system rather than a collection of disconnected exercises.

---

# Assignment Design Rules

Every assignment should contain the following sections.

## 1. Objective

What is being built?

## 2. Required Concepts

Exactly which concepts must be used.

## 3. Required Behavior

What the system must actually do.

## 4. Constraints

What shortcuts are prohibited.

## 5. Failure Cases

What failures must be handled.

## 6. Observability

What behavior the grader can observe.

## 7. Acceptance Criteria

Precise pass/fail requirements.

## 8. Debugging Variant

A deliberately broken version or failure scenario.

## 9. Extension

A meaningful additional requirement.

## 10. Explanation

The learner must explain the important invariants and state transitions.

---

# Exercise Design Rules

Not every exercise should simply ask:

> "Implement X."

Exercises should have distinct pedagogical roles.

Use:

### Mechanism

Learn one primitive.

### Observation

Observe what the primitive actually does.

### Prediction

Predict behavior before running it.

### Variation

Change one important condition.

### Failure

Break the implementation.

### Diagnosis

Determine why the failure occurred.

### Repair

Fix it.

### Composition

Combine several primitives.

### Reconstruction

Build something without being given the implementation strategy.

---

# Recommended Exercise Counts

The original curriculum's roughly 4–8 exercise structure is a good **starting point**, but it should no longer be treated as a hard limit. The current blueprint explicitly describes approximately 4–8 exercises plus a closing gate.

Use conceptual density instead.

| Module | Recommended exercises | Major assignment        |
| ------ | --------------------: | ----------------------- |
| M0     |                     8 | Project harness         |
| M1     |                  9–10 | Mini libc               |
| M2     |                 10–12 | Process runner/pipeline |
| M3     |                 10–12 | Bounded allocator       |
| M4     |                 10–12 | Concurrent work queue   |
| M5     |                 10–11 | Log processor           |
| M6     |                 10–12 | Stateful TCP service    |
| M7     |                 10–11 | Mini supervisor         |
| M8     |                 10–11 | Reliable messaging      |
| M9     |                 10–11 | Ordering service        |
| M10    |                 10–11 | Transaction coordinator |
| M11    |                 10–11 | Replicated log          |
| M12    |                   12+ | Raft node               |
| M13    |                 10–11 | Sharded KV router       |
| M14    |                 10–11 | Membership service      |
| M15    |                 10–11 | Transactional store     |
| M16    |                 11–12 | LSM KV store            |
| M17    |                 10–11 | Observable service      |

These numbers are **targets, not quotas**.

If a concept is genuinely mastered after two exercises, do not create five artificial exercises.

If a concept is extremely difficult, create more.

---

# The Most Important Change

The curriculum should stop treating:

```text
5 exercises completed
```

as equivalent to:

```text
module understood
```

Instead, module mastery becomes:

```text
Can implement it
        ↓
Can explain it
        ↓
Can predict it
        ↓
Can break it
        ↓
Can diagnose it
        ↓
Can repair it
        ↓
Can modify it
        ↓
Can compose it
        ↓
Can recall it later
```

That is especially important for the later distributed-systems modules.

A learner who can implement a Lamport clock once is not necessarily someone who understands distributed ordering.

A learner who can reproduce a Raft vote rule is not necessarily someone who understands Raft safety.

A learner who can implement an SSTable is not necessarily someone who understands storage-engine recovery.

The assignments are what force that transition.

---

# Final Curriculum Architecture

The complete Systems-Forge learning loop should therefore become:

```text
                    ┌──────────────┐
                    │   Reading    │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │   Exercise   │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │   Observe    │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │    Vary      │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │    Break     │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │   Diagnose   │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │  Integrate   │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │  Assignment  │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │    Review    │
                    └──────┬───────┘
                           ↓
                    ┌──────────────┐
                    │     Gate     │
                    └──────┬───────┘
                           ↓
                  ┌───────────────────┐
                  │ Cumulative Review │
                  └───────────────────┘
                           ↓
                  ┌───────────────────┐
                  │ Cumulative Project│
                  └───────────────────┘
```

The current curriculum already has a strong systems ladder from tooling → C → processes → memory → concurrency → I/O → networking → resilience → messaging → distributed ordering → consensus → Raft → sharding → membership → transactions → storage → observability.

**The change is not the roadmap.**

The change is giving each conceptual jump enough room for the learner to actually digest it.

The existing 5-exercise modules become the **spine**.

The additional exercises become the **muscle**.

The assignments become the **integration layer**.

The Gates remain the **proof of mastery**.