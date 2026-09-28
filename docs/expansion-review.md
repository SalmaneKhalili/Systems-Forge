# Systems-Forge — Expansion Plan vs. Current Curriculum Review

> Companion to [`docs/expansion-plan.md`](expansion-plan.md).
> Reviewed against the actual `subjects/` tree (all 18 modules), `exercise.json` graders,
> `forge/internal/cur` + `forge/internal/methods` (grader capabilities), and `PLAN.md`.
> Method: read-only inspection of subjects + harness/config; `solutions/` and `answers/`
> were not inspected and nothing was graded or modified.
>
> Last updated: 2026-09-17

---

## 1. Executive summary

The plan's **design rule** (keep the existing curriculum as the canonical spine, add practice
exercises + assignments, keep gates as proof of mastery) is sound and aligns with the repo's
own non-negotiable: *"Do not delete or casually rewrite existing exercises"* (PLAN.md §1,
"Completeness over speed").

The plan's **anchors are mostly described correctly** for M1, M2, M3, M4, M5, M6, M8, M9,
M13, M14, M16, M17. It gets several things **wrong or overstated** for M0, M7, M10, M12, M15,
and — critically — it **omits the two most integrative exercises that already exist**
(M12-ex06 real 3-node Raft cluster; M15-ex06 durable txstore). About half of the proposed
"major assignments" are **already substantially built** as today's gates / capstones. The
plan's genuinely valuable contribution is a small set of **real integration gaps** that the
curriculum has never graded: crash-restart recovery, fault-injected Raft, wire-level conflict
repair, MPMC queue closure, and cross-module cumulative assignments.

Verdict: adopt the plan's *structure* (spine + muscle + assignment + cumulative) but **rebase
it onto the actual anchors** before authoring anything, and prioritize the true gaps rather
than re-building existing gates.

---

## 2. Module-by-module: plan claim vs. actual curriculum

### M0 — Tools

| Plan claims | Actual (`subjects/M0-tools`) | Verdict |
|---|---|---|
| ex01 `branch` (Makefile w/ ascii-art) | ex01 `makefile` — Makefile discipline | ⚠️ minor name mismatch (no ASCII-art) |
| ex02 `shell` | ex02 `shell` — shell hygiene (`set -euo pipefail`) | ✅ |
| ex03 `checker` | ex03 `checker` — Python tree checker | ✅ |
| ex04 `commits` (Git) | **ex04 `sanitize`** — AddressSanitizer | ❌ **no Git exercise exists anywhere in M0** |
| ex05 `sanitize` | **ex05 `gate` — mini-forge** (scenario grader) | ❌ swapped |
| gate `tools` | ex05 `gate` = mini-forge | ✅ (but it is ex05, not "gate-tools") |

**Gaps the plan correctly spots:** there is no Git exercise today. Shell argument discipline
(quoting, spaces, empty args) is not covered. Build-failure diagnosis is not covered.

**Notes:** the "M0 currently establishes … Git" bullet is false. The proposed assignment
(M0-A01 reproducible harness requiring "Git history following the project's conventions")
has **no prerequisite exercise** for Git — the plan must add a Git exercise or soften the
assignment requirement. Evidence: `subjects/M0-tools/ex01–ex05` dirs, `exercise.json`, `module.md`.

### M1 — C Foundations

| Plan claims | Actual | Verdict |
|---|---|---|
| strlen, strlcpy, memset, strcmp, strdup gate | exactly `ex01-ft_strlen` … `ex05-ft_strdup` | ✅ exact match (gate = ft_strdup per module.md:29) |

Plan's M1 assignment (mini libc) is genuinely additive (module is 5 libft-style single
functions, each with its own harness). Real gaps: `memcpy/memmove` overlap, allocation
failure, multi-function library. Gate instruction "fresh API rather than reproducing the
assignment" is a good idea and matches the repo's "references are not re-authored"
discipline.

### M2 — Processes

Plan anchors (fork/wait, exec, pipe, signals, pipeline gate) = exact actual anchors
(`ex01-fork` … `ex05-pipeline`). ✅. Plan's runner assignment is a genuine extension;
`dup2` redirection is currently handled in M5-ex03, so M2-E11 overlaps M5 — rebalance or
note the fork-point. `M2-E09 descriptor inheritance` is partially probed by M3-ex02
(copy-on-write / mmap+fork) — fine as an observation exercise.

### M3 — Memory

Plan anchors (arena, CoW, pool, mmap, OOM buffer gate) = exact actual anchors
(`ex01-arena` … `ex05-gate` OOM-safe buffer). ✅. Assignment M3-A01 "bounded allocator
with invalid-operation handling + accounting" is **not** covered by today's exercises
(they test bump/CoW/pool/mmap/OOM-growth, not free/reuse/exhaustion/invalid-op on one
region) — genuine gap. ⚠️ The plan calls the M3 gate "growable-buffer/OOM gate"; actual
is `OOM-safe growable buffer` — same thing. ✅

### M4 — Concurrency

Plan anchors (join, mutex, atomic, rwlock, bounded queue gate) = exact actual anchors
(`ex01-join` … `ex05-gate` **bounded queue**). ✅ (rwlock = "readers and writers").

**Key finding:** M4-ex05 is exactly **one producer + one consumer**, no `queue_close`,
API `queue_push/queue_pop` only (`subjects/M4-concurrency/ex05-gate/subject.md:5–19`).
The plan's M4-A01 (multiple producers, multiple consumers, graceful closure) is a **real
gap** — graded under ThreadSanitizer via `build`, so it is feasible as a harness-shape
exercise.

### M5 — Files & I/O

Plan anchors (descriptors, read-vs-stdio, dup2 redirect, tee, log parser gate) = exact
actual anchors (`ex01-descriptors` … `ex05-gate` **log parser**). ✅. Genuine gap the plan
rightly identifies: streaming parser (M5-ex05 already requires cross-read reassembly —
partial overlap), arbitrary-size input, log rotation/reopen. The **log rotation** mechanism
is genuinely unconvered anywhere.

### M6 — Networking

Plan anchors (echo, line chat, HTTP-ish, read timeout, stateful gate) = exact actual
anchors (`ex01-echo` … `ex05-gate` **stateful server**). ✅. 

**Key finding:** the M6-A01 "stateful TCP service" is **almost exactly today's M6-ex05
gate** (greeting, per-connection state, framing, `PING/ECHO/BAD`, accept loop). The plan
adds malformed-input/timeout extras, but the core is already the gate. Also, **concurrent
per-client state cannot be graded by the current `net` runner**, which drives connections
**sequentially** (`forge/internal/methods/procs.go:210-220`; each session dials a fresh
connection only after the previous completes) — any assignment that needs
simultaneously-served clients requires a new runner shape (scenario/fault driver, or an
extended net runner).

### M7 — Resilience

Plan anchors (backoff, breaker, health, shutdown, supervisor gate) = exact actual anchors
(`ex01-backoff` … `ex05-gate` **mini supervisor**). ✅ for content; the plan's "combines
Python and Go" is accurate (ex01 Python, ex02–ex05 Go; `module.md:16`).

**Key finding:** the M7-A01 "mini supervisor" (start workers, detect, retry, restart,
preserve progress, budgets, graceful shutdown) **is verbatim today's M7-ex05 gate**
(`subjects/M7-resilience/ex05-gate/subject.md:1–51`: three workers, watchdog restart of a
crashed worker with `retry ok`, progress preservation, 2-attempt budget → `giving up`,
graceful shutdown). The plan's assignment duplicates the gate; if adopted, it must be
reshaped (e.g. real process workers vs in-process tasks, or add worker crash *between*
successful tasks).

### M8 — Messaging + Fault Injection

Plan anchors (framing, ordered delivery, transparent relay, fault injection, switch gate)
= exact actual anchors (`ex01-framing`, `ex02-orders`, `ex03-switch`, `ex04-switchfaults`,
`ex05-gate`). ✅. **The fault switch exists and works as described** — dup/drop/hold
per-connection (`ex04-switchfaults/subject.md:5-32`). Genuine gap: **dedup + gap repair in
one reliable-messaging layer** — nothing today detects/repairs dups (the switch *injects*
them; ex02 detects gaps but has no dedup). M8-A01 is a real gap. Note: the switch reuse
contract ("later modules grade through the switch") is **documented only**, never realized
(`M8-messaging/module.md:75-77`; PLAN.md defers it) — the plan's M12 assignment that
"the fault switch from M8 must be usable against it" depends on a capability the repo has
not yet wired up.

### M9 — Time & Ordering

Plan anchors (lamport, vector, causality, total order, TCP sequencer) = exact actual
anchors (`ex01-clock` … `ex05-sequencer`). ✅. The M9-A01 "distributed event ordering
service" (multi-node) would **borrow consensus from M10–M12**, which M9's `module.md` and
principles explicitly fence off ("You do NOT need to know: Raft, 2PC, Paxos … M12 handles
that"). Recommend: keep M9 assignment **single-node** (a deterministic ordering service
per sequencer semantics) or move the multi-node variant to a Cumulative assignment.

### M10 — Commit & Consensus

Plan anchors: coordinator, vote, leader election, quorum, **transaction/2PC gateway**.
Actual: `ex01-coordinator`, `ex02-vote`, `ex03-elect`, `ex04-quorum`, `ex05-transaction`.

⚠️ **Overstated:** M10-ex05 is a **static lookup** on a provided `votes.txt` — reply
`COMMIT`/`ABORT` per `tx <id>`; "a missing entry is an abort" (`ex05-transaction/subject.md:5-18`).
There is **no wire-level prepare/vote/decision fan-out**. The plan's M10-A01 (deterministic
2PC coordinator: prepare/vote/abort/**recovery**, deterministic transitions) goes well
beyond the current anchor and its *recovery* sub-requirement presumes persistence (WAL)
that belongs to M15/M16. If adopted, split: a **wire-2PC** exercise fits M10 (participants
spawned in-process, real sockets), and **coordinator crash recovery** should wait for M15.

### M11 — Replicated Log

Plan anchors (append, committed prefix, index arithmetic, snapshot, TCP replica) = exact
actual anchors (`ex01-append` … `ex05-replicagate`). ✅. Note `ex05-replicagate` already
enforces "read never reveals uncommitted" (`ex05-replicagate/subject.md:14-15,43-45`).

⚠️ **Prerequisite fence:** the plan's M11-E08 "conflicting entries / repair divergent
prefixes" pulls in the AppendEntries matching-prefix check, which the curriculum places in
**M12-ex03-logmatch**, and M11's `module.md` explicitly says "You do NOT need to know:
Raft … M12 handles that". Fix: teach prefix matching derivation within M11 (as log math)
or bump conflict-repair to M12.

### M12 — Raft ⚠️ biggest correction

| Plan claims | Actual |
|---|---|
| anchors: term, vote, logmatch, election, TCP gate node | ex01-term, ex02-vote, ex03-logmatch, ex04-election, ex05-gatenode — ✅ |
| — (omitted) | **ex06-cluster: a real 3-node Raft cluster over TCP** (elect/propose/read/log/term; peers on `TARGETPORT+1..+3`; strict majority; higher-term steps down) |

**The plan omits the single most integrative exercise in the whole curriculum.**
`subjects/M12-raft/ex06-cluster/subject.md:5-60` is exactly the "functioning Raft node"
the M12-A01 demands (leader election, terms, voting, log replication, commit
advancement, state machine), minus crash/restart and fault injection. Findings:

- What M12-A01 adds that is genuinely new: **crash/restart recovery**, **fault injection**
  against the cluster with re-election (partition/delay/drop/node kill). Nothing today
  grades these and **no fault machinery exists in M12** (`TARGETFAULTS` appears only in
  M8 files). Feasibility: the `net` runner provides a `restart` step (kill → wait → respawn
  → re-dial) that is **currently unused anywhere** (`forge/internal/methods/procs.go:191-197,256-305`),
  and the `scenario`/`fault` runners exist; `scenario` is used by M0-ex05 only, `fault` is
  registered but never used (`cur.go:35,84`; `methods.go:116`). This is the natural vehicle
  for fault-injected Raft.
- The plan's "Capacity Extension" (partition, delayed/dropped messages, restart, leader
  failure) is precisely what the repo's PLAN.md defers ("M12 subprocess spawn — risks
  flagship raft-capstone determinism", PLAN.md:1079-1080). Flag as **needs design**, not
  just authoring.

### M13 — Sharding

Plan anchors (slot, range, ring, rebalance, shard gate) = exact actual anchors
(`ex01-slot` … `ex05-shardgate`). ✅. **Genuine gap:** routing **during** ownership changes
(migration-time dual-read/redirect). M13-ex04 rebalance is pure move-count math;
`ex05-shardgate` routes to 3 static in-process shards. M13-A01 is feasible with a
harness/scenario grader.

### M14 — Membership

Plan anchors (heartbeat, suspect, gossip, eviction, membership gate) = exact actual anchors
(`ex01-heartbeat` … `ex05-membershipgate`). ✅ (ex04 = "evict"; gate = `membershipgate`).
**Genuine gaps the plan correctly targets:** real node **down/restart** and
**failed-member reintroduction** — today's `drop`→`dead` has no rejoin path
(`ex05-membershipgate/subject.md:47-52`), and ex03-gossip is a pure commutative `Merge`
function. ⚠️ "Gossip membership service" as a *networked* assignment exceeds today's
single-process view; if adopted, keep it in-process tick/gossip (already the module's
style) or move networked version to a Cumulative assignment.

### M15 — Transactions + Chaos

Plan claims M15 "combines 2PC, transactional store, WAL, LWW replay, lost-update,
chaos, tx KV gateway". Actual: `ex01-transaction`, `ex02-replay`, `ex03-conflict`,
`ex04-drill`, `ex05-txstore`, **ex06-durabletxstore** (WAL + commit, `wal` command,
rollback-no-side-effect). Content mostly matches; **4 small corrections**:

1. **The plan omits ex06-durabletxstore** — the existing "crash-durable" WAL artifact
   (ex06-durabletxstore/subject.md:5-61) already does what M15-A01 sketches (transactions,
   WAL, durability, ordering); the assignment adds conflict detection + recovery + noise.
2. `module.md` mislabels ex06 as "multi-node cluster version" while it's a single process
   — the plan should not inherit that claim.
3. **Crash-restart recovery is never actually graded** as of today (see §4, finding G1).
4. M15 already forward-references M16 (`ex06-durabletxstore/subject.md:74-76`) — a plan
   that adds WAL-replay rigor in M15 is fine, but avoid implying hardened WAL (fsync
   semantics) before M16.

### M16 — Storage Engines

Plan anchors (memtable, sstable, WAL, compaction, persistent gateway) = exact actual
anchors (`ex01-memtable` … `ex05-persistent`). ✅ compact = newest-wins + tombstones
(`ex04-compact/subject.md:5-15`). Genuine gaps the plan correctly targets: **crash-during-
flush / crash-during-compaction / recovery-after-restart**. Today's gateway asserts WAL
replay **in prose only** — the grader never restarts the process, only SIGTERM-exit-0
(`ex05-persistent/exercise.json`; `subject.md:13` "a later restart … replays it" is not
tested). A reboot-verify exercise is feasible right now via the existing **unused**
`restart` step.

### M17 — Observability

Plan anchors (metrics, tracing, spans, summaries, telemetry) = exact actual anchors
(`ex01-metrics` … `ex05-telemetry`). ✅ (ex04 = summaries: deterministic p50/p95).
M17-A01 "observable service" is essentially ex05 plus correlation; genuine additive value
is the **request→trace→span→metrics correlation** (M17-E11) and failure observability
(M17-E12), which no exercise couples today.

---

## 3. What the plan gets right (worth keeping)

1. **Design rule** — spine (existing exercises) + muscle (practice) + integration
   (assignments) + proof (gates). Faithful to the repo's canon-and-no-rewrite discipline.
2. **Correct anchor descriptions** for M1, M2, M3, M4, M5, M6, M8, M9, M13, M14, M16, M17.
3. **Real gaps correctly identified** (§4 list below) — these are the valuable part.
4. **Assignment + exercise design rules** (objective/concepts/behavior/constraints/
   failure/observability/acceptance/debug-variant/extension/explanation; mechanism→…→
   reconstruction roles) map cleanly onto the existing per-estimate contract the repo
   already uses in subject.md and `exercise.json` `answers[]`/`readings[]` — adoptable
   as an authoring standard (cf. `docs/readings-standard.md`).
5. **"Targets, not quotas"** matches the repo's quality bar (completeness over speed).

---

## 4. The genuine gaps (adopt these)

Proven by evidence to be **uncovered by any current exercise**:

- **G1. Crash-restart recovery is never graded, anywhere.** The `net` `restart` step exists
  and is proven by fixtures (`PLAN.md:387-389`) but **zero subjects use it**. Targets:
  M16 (persistent LSM replay), M15 (durable txstore replay), M14 (member restart), M12
  (Raft node term/votedFor/log/commit recovery).
- **G2. Fault injection against a real cluster with re-election.** M12-ex06 has no fault
  plane; no later module runs TCP through the M8 switch. `scenario`/`fault` runners exist
  and are barely used (scenario: 1 use; fault: 0).
- **G3. Wire-level log/conflict repair between nodes** (M11→M12 seam; needs fence fix).
- **G4. 2PC coordinator/participant failure + in-doubt recovery** — beyond the static
  votes.txt gateway.
- **G5. Routing during ownership transitions** (M13) and **failed-member reintroduction**
  (M14).
- **G6. MPMC work queue with graceful closure** (M4 — today 1P1C, no close).
- **G7. Concurrent per-client service state** (M6) — blocked by sequential `net` sessions;
  needs a new grader shape (scenario/fault driver or extended runner).
- **G8. Cross-module cumulative assignments** — the repo has none (closest: M15-ex06 /
  M16-ex05 micro-capstones that fuse two modules). CA1–CA5 are genuinely new and valuable.
- **G9. Log rotation & reopen (M5), Git workflow (M0), build-failure diagnosis (M0)**.

---

## 5. Misleading/required corrections before adopting

- **P1.** "Gates" are a **convention, not a flag**: `exercise.json` has no `gate` field and
  `DisallowUnknownFields` would reject one (`cur.go:338`). "Existing X gate" in the plan is
  prose-only. If gates should become machine-visible, add an explicit `kind`/`gate` field
  and update `cur.Exercise` — a platform change, not a subject change.
- **P2.** M0 anchor list is wrong (no Git; ex04= sanitizer, ex05=gate/mini-forge). Rebase
  §M0 "Existing curriculum".
- **P3.** M12: add `ex06-cluster` (real 3-node Raft) to the anchor list; the M12-A01 must
  be scoped to *crash/fault* additions, not a from-scratch Raft node.
- **P4.** M15: add `ex06-durabletxstore`; M15-A01 ≈ ex06 generalized, not new.
- **P5.** M7-A01 ≈ today's M7-ex05 gate; reshape or drop.
- **P6.** M6-A01 ≈ today's M6-ex05 gate; reshape with the concurrency caveat (G7).
- **P7.** M10 "2PC gateway" is a static lookup; wire-2PC is additive, recovery waits for M15.
- **P8.** M11/M9 assignments must respect the module fences (`M11 module.md:43-44`,
  `M9 module.md:45`) or the plan must move those pieces.
- **P9.** M12 "fault switch must be usable" currently depends on an **unbuilt capability**
  (switch reuse deferred). Decide: build switch reuse as part of this work, or drop that
  requirement.
- **P10.** Per-exercise count targets (8/10-12 per module, 12+ for M12) are aspiration;
  the repo's "targets, not quotas" and both-ways-fixture discipline should govern.

---

## 6. Grader feasibility for the proposed assignments

Capabilities that exist and are proven / available but unused:

| Capability | Status | Enables |
|---|---|---|
| `net` `signal` + `wait_exit` steps | in use (M14-ex05, M15-ex05, M16-ex05) | graceful-shutdown verification |
| `net` `restart` step (+`TARGETRESTART=1`) | **available, zero subjects use it** | G1 crash-restart assignments |
| `scenario` runner (`FORGE_RESULT`, `FORGE_ROOT`, `FORGE_EXERCISE_DIR`) | 1 use (M0-ex05) | multi-process drivers, G2/G7 |
| `fault` runner | registered, 0 uses | G2 chaos |
| `lincheck` runner | 0 uses | linearizability for M4 queue / M13 router |
| `TARGETPORT`/`TARGETHOST` injection | auto in `net` | standard |
| `build` + sanitizers (ASan/UBSan/TSan) | standard | M4 TSan, C assignments |

Limits (blockers if unaddressed):

- **`net` sessions are sequential** — no concurrent-clients grading (G7) without a new
  runner shape or a scenario driver.
- **`net` starts exactly one process** — a cluster must be a single binary spawning peers
  (as M12-ex06 does on `TARGETPORT+1..+3`).
- **`TARGETFAULTS` is not auto-injected** into `net` runs — per-exercise `StartEnv` would be
  needed; no exercise sets it today.

---

## 7. Recommended adoption path

1. **Rebase plan §M0, §M7, §M10, §M12, §M15** onto the actual anchors (P2–P5, P7, P3, P4).
2. **Priority 1 (cheap, high-value, infra ready):** crash-restart recovery exercises using
   the existing `restart` step — M16 (reboot-verify the LSM), M15 (durable txstore
   reboot), M12 (Raft node state recovery). Follow the both-ways fixture discipline.
3. **Priority 2 (main new engineering):** fault-injected Raft (G2) via a `scenario`/`fault`
   driver, plus decide/deliver M8-switch reuse (P9).
4. **Priority 3 (grader-shape work):** MPMC queue closure (G6, TSan `build` — feasible
   now), and a concurrency-capable runner for M6 (G7).
5. **Cumulative assignments (G8):** CA1–CA5 are the plan's biggest structural addition and
   have no conflict with the existing tree; start with CA1 (after M3) and CA4 (Raft bridge).
6. **Before any new subject ships:** author → solve → prove both pass/fail → add fixtures →
   grow selftest (repo's standing rule, PLAN.md §3/§9).

---

## 8. Conclusion

The expansion plan is a **good structural proposal with an outdated map**. Rebased onto the
real anchors it becomes a strong roadmap: keep the spine, add the muscle, and — most
importantly — target the nine genuine gaps (crash-restart recovery chief among them) with
assignments the current grader **already has the machinery to grade**. The gates remain the
proof of mastery; the plan's new contribution is making mastery *stick* through integration
and cumulative projects rather than re-listing the five-exercise modules we already have.