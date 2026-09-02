# systems-forge — Master Plan & Progress Tracker

> **Purpose:** keep the full curriculum, platform contract, and progress on disk so nothing is
> lost across sessions/compaction. Update this file every time you finish a chunk of work.
>
> Canonical repo: `/home/salmane/GolandProjects/systems-forge`
>
> Last updated: 2026-09-02

---

## 1. Mission

Build a **complete, offline, self-graded systems engineering → distributed systems
curriculum** ("piscine" style) that carries a total beginner from zero to a job-ready
Infra/Platform Software Engineer over ~18–24 months.

Core principles (user-set, non-negotiable):

1. **Completeness over speed.** No half-arsed exercises. Everything is a *real artifact*:
   actual servers, actual allocators, actual protocols — never toy stand-ins.
2. **Everything is graded automatically and deterministically.** Grader passes *and* fails
   are validated with fixtures before an exercise ships.
3. **42-piscine discipline** where it applies (C foundation): `-std=gnu11 -Wall -Wextra -Werror`,
   sanitizers, Makefile discipline (`all`/`fclean`/`re`).
4. **Readings with comprehension checkpoints.** Every exercise points at exact chapters /
   articles and auto-grades a quiz on them.
5. Static, offline, local: no server daemon; all state in a local SQLite DB.

---

## 2. Operating Constraints (user decisions)

| Topic | Decision |
|---|---|
| Start level / time | Beginner, ~15–25 hrs/week, 18–24 month runway |
| Target job | Infra / Platform Software Engineer |
| Delivery form | Static offline package + **TUI** (accepted; CLI verbs must also exist) |
| Languages | Go (platform + distributed capstones), C (foundation exercises), Python stdlib (graders/harness), Bash (glue/Makefiles) |
| Repo name | **systems-forge** (renamed from `piscine-predictor`, approved) |
| Go deps | bubbletea, bubbles, viewport, glamour, lipgloss (TUI); `modernc.org/sqlite` (pure-Go, no cgo) |
| Toolchain | User: "download it if you have to" ⇒ Go 1.27 installed user-local |
| Curriculum | Full systems ladder; Raft is the flagship capstone |
| Jobcraft | Job-search engineering track starts ~month 10–12 |
| Work style | One thing thoroughly at a time; maximum effort; no half-assing |

---

## 3. Status Snapshot (2026-08-30)

### Done (verified)
- [x] Toolchain verified: gcc 13.3.0, clang 18.1.3, GNU Make 4.3, python3 3.12.3, git 2.43.0.
      `shellcheck` **not** installed (could not install offline — revisit).
- [x] Go 1.27.0 linux/amd64 installed to `~/.local/go`; `PATH` export appended to `~/.bashrc`.
- [x] Repo created + `git init -b main`.
- [x] Platform core compiles: `go build ./...` green; `go vet ./...` green.
- [x] Binary `bin/forge` builds and runs (`help`, `list`, `score`, `selftest` work).
- [x] **Grader self-test: 23 fixtures, 0 mismatches** — all 10 method types validated
      both-ways (pass + fail). See §9.
- [x] Fixed during bring-up: missing `cur` imports (procs.go, scenario.go, lin.go);
      `normalize` param shadowing package func in exec.go; unused vars; slice-bounds panic
      in `run()` when invoked with no args; **every runner was ignoring the exercise
      working dir** — `Ctx.Opts()` now sets `Dir: c.Dir`, build/verification runs use
      `buildDir`.

- [x] Root `Makefile` (setup/selfcheck/init/list/show/check/score/selftest/tui) + `tools/selfcheck.sh` — green.
- [x] `forge init` mirroring + `score` validated against `M0-tools` (2-ex exercise set).
- [x] M0 "tools" module authored (`subjects/M0-tools/ex01–ex05`); M0-ex01 solved in `answers/` → PASS.
      (M0-ex03/04/05 shipped later — see tracker item #23 and §14 2026-09-02 → M0 5/5.)
- [x] **TUI layout reworked to side-by-side two-column** (module tree | spec pane) — verified under tmux:
      header with progress, tree with ✓/·/✕ marks, spec renders subject.md, `G`/`n`/arrows
      refresh the spec, `browse`/`score`/`help` views intact.
- [x] **M1 "C gate" module authored + solved** — `subjects/M1-clib/ex01–ex05` all `build`+`quiz`
      green; both-ways determinism proven; `forge score` → `M1-clib 5/5`.
- [x] **M2 "Processes" module authored + solved** — `subjects/M2-procs/ex01–ex05` all `build`+`quiz`
      green; broken no-handler signal caught (`exit -1`) and unclosed write-end caught
      (`timed out`); `forge score` → `M2-procs 5/5`, `TOTAL 11/15`.
- [x] **M3 "Memory" module authored + solved** — `subjects/M3-memory/ex01–ex05` (bump arena,
      copy-on-write via mmap+fork, fixed slot pool, mmap allocator, OOM-safe gate) all
      `build`+`quiz` green; broken solutions proven to FAIL (overlapping pool slots,
      `MAP_SHARED` misuse, refusal ignored); `forge score` → `M3-memory 5/5`,
      `TOTAL 16/20`.
- [x] **M4 "Concurrency" module authored + solved** — `subjects/M4-concurrency/ex01–ex05`
      (join & exit, mutex bomb, atomic counter, readers & writers, bounded-queue gate) all
      graded under ThreadSanitizer (`sanitizers:["thread"]`, runs wrapped in
      `setarch x86_64 -R` for TSan); all `build`+`quiz` green; broken solutions proven to
      FAIL — plain-increment bombs (ex02/03), lock-free rwlock client (ex04), mutex-free
      queue (ex05) all exit 66 (TSan `data race`); wrong exit-value plumbing (ex01) → `sum 15`;
      LIFO queue (ex05) → `fifo violated`; `pthread_cond_wait(…, NULL)` rejected by
      `-Werror=nonnull` at build time; `forge score` → `M4-concurrency 5/5`, `TOTAL 21/25`.
- [x] **M5 "Files & I/O" module authored + solved** — `subjects/M5-filesio/ex01–ex05`
      (descriptors, read-vs-stdio, dup2 redirect, tee, log-parser gate) all `build`+`quiz`
      green (ex04 also `artifact`); one run on two fixtures each (ex02/ex05); broken
      solutions proven to FAIL — fd miscount (`3 and 3`), buffer-size counting (`read 49`
      + `mismatch`), missing stdio flush (`captured: ` empty), stdout-only tee (`readback
      0`), wrong delimiter (`parse failed`); quiz trap found + handled: question text with
      an internal `:` breaks `parseKV`'s first-colon split, so ex04 Q2 was reworded;
      `forge score` → `M5-filesio 5/5`, `TOTAL 26/30`.
- [x] **M6 "Networking" module authored + solved** — `subjects/M6-networking/ex01–ex05`
      (echo, line chat, HTTP-ish, read timeout, stateful gate) all `build`+`net`+`quiz`
      green — **first curriculum use of the `net` grader**; the `net` runner grew
      **multi-connection sessions** (`connections: [{steps}, …]`, one server process,
      one fresh TCP dial per session) so accept-loop persistence and cross-connection
      state are gradeable deterministically; fixture `net/{pass, fail}/ex02-*` added
      (23 fixtures, 0 mismatches). Broken solutions proven to FAIL — wrong port (ex01,
      `start` refused), no accept loop (ex02, `conn 2` refused), 404 answered with 200
      (ex03, byte mismatch), missing `SO_RCVTIMEO` (ex04, deadline), no greeting-on-accept
      (ex05, `conn 1 step 1` no data); `forge score` → `M6-networking 5/5`,
      `TOTAL 31/35`.

- [x] **M7 "Resilience" module authored + solved** — `subjects/M7-resilience/ex01–ex05`
      (backoff, circuit breaker, health checks, graceful shutdown, supervisor gate). ex01 is
      the **first `lang: python` exercise, graded via the `stdout` method**; ex02–ex05 are
      Go (shared Makefile; breaker uses an injected fake clock). All `build`/`stdout`+`net`
      +`quiz` green; broken solutions proven to FAIL (five bug classes: no retry loop,
      breaker that never trips, `/down` that never flips state, no graceful path, flaky task
      failing every retry); `forge score` → `M7-resilience 5/5`, `TOTAL 36/40`.

- [x] **M8 "Messaging + the Switch" module authored + solved** —
      `subjects/M8-messaging/ex01–ex05` (framing, ordered delivery, transparent switch
      relay, switch fault injection, switch gate). ex01/ex02 harness-shape Go graded
      `build`+`stdout`+`quiz`; ex03–ex05 whole-program `switch.go` + provided counter
      backend, graded `build`+`net`+`quiz` — **no grader changes needed** (switch spawns
      its own backend; `StartEnv` carries `TARGETFAULTS`; fault effects proven at the
      data level). Broken solutions proven to FAIL (five bug classes, see §10 M8
      detail); new fixtures `net/{pass, fail}/ex05-*` (25 fixtures, 0 mismatches);
      `forge score` → `M8-messaging 5/5`, `TOTAL 41/45`.

- [x] **2026-08-31 bug-fix + hardening pass** — (a) M0-ex02 now graded via
      `bash solve.sh` (reference uses `set -euo pipefail`; `sh` on Debian is `dash`
      and rejects it) → M0-tools 2/5; (b) `artifact.wipe` config added and applied
      to M5-ex04's `out.txt`, closing stale-file inheritance (both-ways re-proven);
      (c) progress.db history rebuilt — every authored exercise re-checked, all PASS.
      `forge score` → `TOTAL 42/45`. See §14 2026-08-31 entry.

- [x] **M9 "Time & Ordering" module authored + solved** —
      `subjects/M9-time/ex01–ex05` (lamport clock, vector clock, causality, total
      order, sequencer). ex01–ex04 harness-shape Go graded `build`+`stdout`+`quiz`;
      ex05 whole-program `srv.go` TCP gateway, graded `build`+`net`+`quiz` (two
      connections, byte-exact replies). All rules honor the no-wall-clock law.
      Broken solutions proven to FAIL (five bug classes — ignored receive merge,
      ignored vector fold, non-strict happened-before, pid-major ordering, no-stamp
      merge in the sequencer); refs restored immediately after each. New net
      fixtures `pass/ex05-clock` + `fail/ex05-reset` (sequencer across connections)
      → selftest **27 fixtures, 0 mismatches**; `forge score` → `M9-time 5/5`,
      `TOTAL 47/50`.

### In progress
- [x] **M1 "C gate" module authored + solved** — `subjects/M1-clib/ex01–ex05`, all graded
      `build`+`quiz`. See §10 (M1 detail). Reference solutions in `answers/`, all PASS;
      broken `ft_strlen` verified to FAIL (both-way determinism).
- [x] **M2 "Processes" module authored + solved** — `subjects/M2-procs/ex01–ex05`
      (fork/wait, exec-spawn, pipes, signals, pipeline gate), graded `build`+`quiz`.
      Verified: all PASS; broken signal (no handler) caught as `exit -1`; broken pipe
      (unclosed write end) caught as `timed out`.
- [x] **M3 "Memory" module authored + solved** — `subjects/M3-memory/ex01–ex05` (bump arena,
      CoW virtual-memory experiment, fixed-slot pool, mmap allocator, OOM-safe growable
      buffer gate), graded `build`+`quiz`. Verified: all PASS; broken solutions proven to
      FAIL (three bug classes: slot overlap, wrong mapping flag, ignored refusal).
- [x] **M4 "Concurrency" module authored + solved** — `subjects/M4-concurrency/ex01–ex05`
      (join & exit, mutex bomb, atomic counter, readers & writers, bounded-queue gate),
      graded `build`+`quiz` under ThreadSanitizer. Verified: all PASS; broken solutions
      proven to FAIL (six bug classes, see §10 M4 detail).
- [x] **M5 "Files & I/O" module authored + solved** — `subjects/M5-filesio/ex01–ex05`
      (descriptors, read vs stdio, dup2 redirect, tee, log-parser gate), graded
      `build`+`quiz` (+ `artifact` on ex04's `out.txt`). Verified: all PASS; broken
      solutions proven to FAIL (five bug classes, see §10 M5 detail).
- [x] **M6 "Networking" module authored + solved** — `subjects/M6-networking/ex01–ex05`
      (echo, line chat, HTTP-ish, read timeout, stateful gate), graded
      `build`+`net`+`quiz` — first curriculum use of the `net` grader; multi-connection
      `connections` sessions added to the grader for accept-loop/state proofs. Verified:
      all PASS; broken solutions proven to FAIL (five bug classes, see §10 M6 detail).
- [x] **M7 "Resilience" module authored + solved** — `subjects/M7-resilience/ex01–ex05`
      (backoff, circuit breaker, health checks, graceful shutdown, supervisor gate), graded
      `stdout`/`build`+`net`+`quiz` — ex01 is the first `lang: python` exercise, ex02–ex05
      Go on the shared Makefile. Verified: all PASS; broken solutions proven to FAIL (five
      bug classes, see §10 M7 detail).
- [x] **M8 "Messaging + the Switch" module authored + solved** —
      `subjects/M8-messaging/ex01–ex05`, graded `stdout`/`build`+`net`+`quiz`. ex01/ex02
      harness-shape (frame codec, ordered queue); ex03–ex05 the switch (relay, fault
      injection, gate) with a provided counter backend. Verified: all PASS; broken
      solutions proven to FAIL (five bug classes, see §10 M8 detail).
- [x] **M9 "Time & Ordering" module authored + solved** —
      `subjects/M9-time/ex01–ex05` (lamport clock, vector clock, causality, total
      order, sequencer), graded `build`+`stdout`/`net`+`quiz`. ex01–ex04 harness-
      shape Go; ex05 a whole-program TCP sequencer. Verified: all PASS; broken
      solutions proven to FAIL (five bug classes, see §10 M9 detail).
- [x] **M10 "Commit & Consensus" module authored + solved** —
      `subjects/M10-commit/ex01–ex05` (coordinator, vote, elect, quorum, transaction),
      graded `build`+`stdout`/`net`+`quiz`. ex01–ex04 harness-shape Go; ex05 a
      whole-program 2PC TCP gateway deciding from a `votes.txt` file. Verified: all
      PASS; broken solutions proven to FAIL (five bug classes, see §10 M10 detail).
- [x] **M11 "Replicated Log & Consistency" module authored + solved** —
      `subjects/M11-log/ex01–ex05` (append, committed, index, snapshot, replica),
      graded `build`+`stdout`/`net`+`quiz`. ex01–ex04 harness-shape Go; ex05 a
      whole-program TCP replica serving only the committed prefix across
      connections. Verified: all PASS; broken solutions proven to FAIL (five bug
      classes, see §10 M11 detail); new fixtures `pass/ex05-replica` +
      `fail/ex05-replicaleak` → **31 fixtures, 0 mismatches**;
      `forge score` → `M11-log 5/5`, `TOTAL 57/60`.
- [x] **M12 "Raft — flagship capstone" module authored + solved** —
      `subjects/M12-raft/ex01–ex05` (term, vote, logmatch, election, gatenode),
      graded `build`+`stdout`/`net`+`quiz`. ex01–ex04 harness-shape Go (monotonic
      term, up-to-date-log vote rule, matching-prefix AppendEntries check,
      strict-majority quorum); ex05 a whole-program TCP raft gate node serving a
      monotonic term + committed log across connections. Verified: all PASS;
      broken solutions proven to FAIL (five bug classes, see §10 M12 detail);
      new fixtures `pass/ex05-gatenode` + `fail/ex05-gateregress` →
      **33 fixtures, 0 mismatches**; `forge score` → `M12-raft 5/5`, `TOTAL 62/65`.
- [x] **M13 "Sharding" module authored + solved** —
      `subjects/M13-shard/ex01–ex05` (slot, range, ring, rebalance, shardgate),
      graded `build`+`stdout`/`net`+`quiz`. ex01–ex04 harness-shape Go (FNV slot
      assignment, key-range routing, consistent-hash ring lookup, move-count
      rebalancing); ex05 a whole-program TCP key-range shard gateway with three
      isolated shards. Verified: all PASS; broken solutions proven to FAIL (five
      bug classes, see §10 M13 detail); new fixtures `pass/ex05-shardgate` +
      `fail/ex05-shardwrong` → **35 fixtures, 0 mismatches**;
      `forge score` → `M13-shard 5/5`, `TOTAL 67/70`.
- [x] **M14 "Membership" module authored + solved** —
      `subjects/M14-membership/ex01–ex05` (heartbeat, suspect, gossip, evict,
      membershipgate), graded `build`+`stdout`/`net`+`quiz`. ex01–ex04
      harness-shape Go (heartbeat countdown, alive→suspect→failed lifecycle,
      commutative version-merge gossip, staleness eviction — all tick-modeled,
      no wall clock); ex05 a whole-program TCP membership gateway tracking
      alive/suspect/dead per member. Verified: all PASS; broken solutions proven
      to FAIL (five bug classes, see §10 M14 detail); new fixtures
      `pass/ex05-membership` + `fail/ex05-membershipdrop` → **37 fixtures,
      0 mismatches**; `forge score` → `M14-membership 5/5`, `TOTAL 72/75`.

### Blocked / deferred
- `shellcheck` unavailable offline (M0 exercise that uses it will note the alternative).
- Virtual switch `switch/` is M8 student work, not Stage 1 platform work.
- Fixture binaries under `tools/fixtures/**/hello` — git-ignored.
- M0-ex03/04/05 and the learner-written deliverable pattern: references for M0's
  student-authored exercises are not authored yet (M0 stands at 2/5: ex01 + ex02
  after the `bash` fix). Deliberate — M0 refs get completed alongside the M9+ push or
  when M0 is revisited.

---

## 4. Toolchain Cheat-Sheet

| Tool | Version / path | Notes |
|---|---|---|
| Go | go1.27.0 linux/amd64 @ `~/.local/go/bin/go` | PATH export in `~/.bashrc`; new shells may not have it — use `export PATH="$HOME/.local/go/bin:$PATH"` |
| gcc | 13.3.0 | C build grader |
| clang | 18.1.3 | alternative toolchain |
| make | GNU Make 4.3 | build discipline |
| python3 | 3.12.3 | stdlib-only graders / drivers |
| git | 2.43.0 | |
| shellcheck | **missing** | revisit; do not depend on it |
| `file`, `nm` | present | used by artifact grader |

---

## 5. Repo Layout

```
systems-forge/
├── PLAN.md                  ← this tracker
├── README.md                ← user-facing usage contract
├── Makefile                 ← (planned) setup/list/show/check/score/selfcheck
├── bin/forge                ← built binary (git-ignored)
├── forge/                   ← Go module `forge` (go 1.27)
│   ├── go.mod
│   ├── cmd/forge/           ← main.go + cli.go (verbs, usage, renderReport, cmdSelftest)
│   └── internal/
│       ├── cur/cur.go       ← exercise.json model, strict JSON, defaults, tree loader
│       ├── sandbox/         ← ulimit+setsid+timeout wrapper, Background handles
│       ├── store/store.go   ← SQLite (modernc), runs table, progress queries
│       ├── methods/         ← 10 runner funcs + Ctx/Part/Result/Dispatch
│       ├── check/check.go   ← Engine: FindExercise, WorkDir mirror, Check orchestration
│       └── tui/app.go       ← bubbletea tree/spec/score/check UI
├── subjects/                ← curriculum tree (modules/, exercises/)   [mostly empty — M0 next]
│   └── M0-tools/{module.md, exNN-slug/exercise.json, subject.md, …}
├── caps/                    ← capstone specs (queued, M8+)
├── tools/
│   └── fixtures/            ← <method>/<expect>/<case>/…   (21 fixtures, all green)
├── answers/                 ← student workspace (mirrored from subjects/, git-ignored)
└── progress.db              ← run history (git-ignored)
```

### Exercise dir contract
- Module dir: `NN-slug` (e.g. `M0-tools`); optional `module.md` (first `# ` line = title).
- Exercise dir: `exNN-slug`; must contain `exercise.json` + `subject.md`.
- Grading config parsed by `cur.LoadExercise` with `DisallowUnknownFields` (typos are errors).
- `exercise.json` JSON schema — see §7.

---

## 6. Platform Architecture (Stage 1 core — SHIPPED)

```
forge (TUI / CLI)
  │
  ├── cur.Load(subjects/)            → curriculum in memory
  ├── check.Engine
  │     ├── FindExercise(id)         → module + exercise
  │     ├── WorkDir(ex)              → answers/<rel> mirror (never clobber student files)
  │     └── Check(ctx, ex)           → run methods in order, stop on first fail
  │           └── methods.Dispatch → runner per method type
  │                                   each runner uses sandbox.Options (Dir, Env, ulimits, timeout)
  │                                   via Ctx.Opts() (Dir defaults to exercise workdir)
  └── store.Store (progress.db, WAL)
        └── runs(exercise,status,detail,stdout,duration,runtime)
        └── LastRun / RunHistory / ModuleProgress / OverallProgress
```

Data flow for one check:
1. `cur.Load` parses all `exercise.json`s (strict).
2. `WorkDir` mirrors scaffolding from `subjects/<mod>/<ex>` to `answers/<mod>/<ex>` —
   only when subject file is newer (student edits never overwritten).
3. `methods.Dispatch` runs each declared method in order; first failing method stops the run.
4. `store.RecordRun` persists the outcome; TUI/CLI render from the store.

### Sandbox rules (every command)
- `LC_ALL=C TZ=UTC TERM=dumb`; PATH/HOME inherited (deterministic env overlay).
- Fresh session via `setsid` (pgid == child pid); kill-tree on timeout; grace SIGTERM→SIGKILL.
- Optional `ulimit -v/-f/-u/-n` from `exercise.json`.
- Timeout from `exercise.json` (default 20s); background tasks return `Handle` (must Stop).

---

## 7. Grader Method Contract (10 runner types)

| type | runner | config block | verdict |
|---|---|---|---|
| `build` | `exec.go` | `build{make,clean,fclean? sub,runs[]}` | `make fclean` (failure tolerated, timeout=fail) → `make <target>` under CFLAGS (default `-std=gnu11 -Wall -Wextra -Werror` + sanitizers) → each `runs[]` cmd diffed vs `expect` file (relative path in exercise dir) with `normalize` mode |
| `stdout` | `exec.go` | `stdout{cmd,input?,expect,normalize,stderr_ok}` | single run, `expect` file diff (normalize strip/blank/trailing) |
| `artifact` | `exec.go` | `artifact{path,exists,is_a,min_size,max_size,symbols[],contains[]}` | existence (both directions), `file -b` prefix, size bounds, `nm` symbols, byte-substring scan of the file |
| `process` | `procs.go` | `process{start,wait_ms,target_env,probes[],kill}` | spawn stale target, wait, run probes with `TARGET_PID` (pgid) env, expect exit std-out substring |
| `net` | `procs.go` | `net{start,port,host,wait_ms,timeout_ms,start_env,steps[]}` | spawn server (`TARGETPORT`/`TARGETHOST` env), retry-dial, per-step send bytes / expect bytes (full or substring), Go-style escapes + `*_hex` |
| `quiz` | `exec.go` | `quiz{file=quiz.txt}` + exercise `answers[]` | parse `key: value`/`key=value`, case-insensitive match, per-answer part |
| `report` | `exec.go` | `report{cmd,reference,tolerate{key:rel}}` | run cmd, parse k=v, compare to reference file, optional relative tolerance on numeric fields |
| `scenario` | `scenario.go` | `scenario{driver,env,timeout_s}` | run driver; **contract**: exit 0 ⇒ pass; `FORGE_RESULT: pass|fail` overrides; `FORGE_META: k=v` captured; driver gets `FORGE_ROOT` + `FORGE_EXERCISE_DIR` |
| `fault` | `scenario.go` | `fault{driver,env,timeout_s}` | same contract as scenario (kept distinct for future divergence) |
| `lincheck` | `lin.go` | `lincheck{history}` | glob from exercise dir, backtracking linearizability search (`maxLinOps=220`, `searchBudget=2_000_000`), real-time order via `End<Start`, every read must match last write; witness on fail; budget exceed = error (never a false verdict) |

Normalize modes: `strip` (concatenate fields), `blank` (trim trailing ws, collapse blank runs), `trailing`/`` (trim trailing ws, drop final blanks), else exact.

Decode of byte fields: values wrapped in `"…"` are unquoted (Go escapes: `\n`, `\xNN`…); else raw.

### Determinism rule (NON-NEGOTIABLE)
Every grader method must be able to *both* pass a correct solution **and** fail a broken one,
proven by automated fixtures before the exercise is accepted. Verified per method in §9.

---

## 8. Store / Progress Model

- DB: `progress.db` (SQLite via modernc, WAL, busy_timeout), single `runs` table.
- Row: exercise id, status (`pass|fail|error`), method detail JSON, stdout scratch, duration ms, timestamp.
- Queries used by UI: `LastRun(id)`, `RunHistory(id, n)`, `ModuleProgress(ids)`, `OverallProgress(ids)`.
- Note: `ModuleProgress` currently just re-reads `LastRun` per exercise (fine at curriculum scale).

---

## 9. Fixture Inventory (tools/fixtures) — ALL GREEN 2026-08-30

Layout: `tools/fixtures/<method>/<expect>/<case>/exercise.json` (+ artifacts). `expect` = `pass|fail`;
only used by `forge selftest`; selftest runs methods in the fixture dir directly (no answers mirror).

| method | pass | fail |
|---|---|---|
| stdout | `pass/ex01-hi` prints "hello forge" (strip-normalized diff) | `fail/ex01-wrong` prints "hello wrong" |
| build | `pass/ex01-hello` real `Makefile{C}`, compiles under Werror, `./hello` → expected.txt | `fail/ex01-broken` missing `;` never compiles |
| artifact | `pass/ex01-flag` ASCII text present, populated, contains "FORGE" | `fail/ex01-missing` absent file; `fail/ex02-unexpected` forbidden file exists |
| quiz | `pass/ex01-ok` both answers right | `fail/ex01-wrong` one answer wrong |
| report | `pass/ex01-ok` output matches reference (tolerate on writes) | `fail/ex01-drift` writes=99 vs 3 (out of tolerance) |
| process | `pass/ex01-alive` `sleep 30`, probe `kill -0 $TARGET_PID` | `fail/ex01-gone` target exits immediately |
| net | `pass/ex01-echo` server replies PONG on one connection (17001); `pass/ex02-sessions` two fresh `connections` on one server (17003) | `fail/ex01-wrong` server replies NOPE (17002); `fail/ex02-once` server serves one connection then exits (17004) |

Module-gate `net` fixtures (ex05 gateways), all both-ways:
`ex05-membership`(pass)/`-membershipdrop`(fail), `ex05-persistent`(pass)/`-persistonce`(fail, M16 LSM gateway),
`ex05-telemetry`(pass)/`-telemetryfresh`(fail, M17 telemetry gateway), `ex05-txstore`(pass)/
`-txstoreleak`(fail, M15 txstore gateway), `ex06-nodewnly`(fail, M12 raft cluster — no
replication to followers), `ex06-walded`(fail, M15 durable txstore — writes never reach the WAL),
plus M4/M8/M11/M12/M13/M15 gateways → selftest **45 fixtures, 0 mismatches**.
| scenario | `pass/ex01-ok` exit 0 + FORGE_META | `fail/ex01-ko` exit 1 |
| fault | `pass/ex01-ok` exit 1 + `FORGE_RESULT: pass` (tests override) | `fail/ex01-ko` exit 0 + `FORGE_RESULT: fail` (tests override) |
| lincheck | `pass/ex01-lin` linearizable single-register history | `fail/ex01-nonlin` real-time contradiction (read of v1 must follow write of v2) |

**To do:** extend fixtures as method features grow (e.g. net `contains`/hex/sleep steps, artifact `is_a`/sizes, process kill-mode, stdin, `expect_exit`).

---

## 10. Curriculum Blueprint (M0 → M14)

Piscine-style modules; each module ≈ 4–8 exercises + a closing "gate" exercise. Every exercise ships:
`subject.md` (goal, constraints, acceptance criteria, readings), `exercise.json` (grader), any scaffolding.

| Module | Focus | Languages | Key artifacts (real things, not toys) |
|---|---|---|---|
| **M0** | Tools of the trade | sh, Make, git, Python | Make-discipline exercices (`all/fclean/re`), shell with `set -euo pipefail`, tiny checker script, git workflows, sanitizer flags |
| **M1** | C gate (42 piscine style) | C | libft-style reimplementation bank, strings/memory utils under `-Wall -Wextra -Werror`, ASan/UBSan clean |
| **M2** | Processes | C | fork/exec/wait, pipelines, signal handling, a tiny shell |
| **M3** | Memory | C | malloc arena, virtual-memory experiments, `mmap` allocator, OOM safety |
| **M4** | Concurrency | C, Go | threads, mutex/condvar, atomics, producer–consumer, race-free queues |
| **M5** | Files & I/O | C, Python | file descriptors, redirection, `open/read/write` abstractions, a mini log parser |
| **M6** | Networking | C, Go | sockets, echo/chat server, HTTP-ish server, timeouts and non-blocking I/O |
| **M7** | Resilience | Go, Python | retries/backoff, circuit breaker, health checks, graceful shutdown, watchdog |
| **M8** | Messaging + the Switch | Go | RPC-ish transport, framing, **the fault-injection switch** (`switch/`) the learner builds and reused by every later grader |
| **M9** | Time & ordering | Go | clocks, vector clocks, ordering, logical time experiments |
| **M10** | Commit & Consensus | Go | coordinator, vote, leader election, quorum, 2PC transaction gate |
| **M11** | Replicated Log & Consistency | Go | append log, commit index, snapshot, TCP replica (read committed prefix only) |
| **M12** | **Raft (flagship capstone)** | Go | full leader election + log replication + safety under the switch |
| **M13** | Sharding | Go | hashing/key-range sharding, rebalancing |
| **M14** | Membership | Go | gossip membership, failure detection |
| **M15** | Transactions + chaos | Go | 2PC/transactions, chaos drills against a full cluster |
| **M16** | Storage Engines | Go | LSM-tree: memtable, SSTables, WAL, compaction, disk-persistent KV gateway |
| **M17** | Observability | Go | Metrics registry, tracing (spans/attrs), latency summaries, telemetry gateway |

Jobcraft track (resume/LinkedIn/practice/apply pipeline) is **scheduled starting ~month 10–12**, interleaved.

### Grading language split
- M0–M4 externals may be auto-graded by `build`+`stdout`+`artifact`+`quiz`+`report`.
- M2+ process/net exercices use `process`/`net`.
- M5+ race-free and M8+ distributed exercices use `scenario`/`fault` drivers + `lincheck`.
- Determinism: keep lincheck windows small; drivers sample histories when big.

### M0 detail (immediate next milestone — write these next)
1. `ex01-branch` … Makefile with `all/fclean/re` + ASCII-art discipline. (Grader: `build` with `make re`)
2. `ex02-shell` … a 20-line shell script with `set -euo pipefail`, arg handling. (Grader: `stdout` + `quiz`)
3. `ex03-checker` … a Python stdlib script that checks file/dir invariants (mini-forge clone). (Grader: `stdout` + `artifact`)
4. `ex04-commits` … git history shape / convention exercise. (Grader: `scenario` or `report`)
5. `ex05-sanitize` … compile the same program under `-fsanitize=address` and make a real UAF/detect via artifact/symbols. (Grader: `build` + `process`)
6. `gate-tools` … the M0 gate: combined hygiene + tooling exam.

Each with exact readings (chapters of e.g. "The Linux Programming Interface", man pages, `shellcheck`), and a quiz.

### M1 detail (authored 2026-08-30) — module `M1-clib`, all green
Shape: student writes the implementation `.c` + a `Makefile` (`all`→`./test`, `fclean`, `re`);
`forge` injects `-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined` via CFLAGS and
compiles `main.c` (provided harness) + their file, then diffs `./test` stdout vs `expected.txt`
(normalize `strip`) and runs the `quiz`.

1. `ex01-ft_strlen` — counting loop incl. embedded-NUL case. **Known trap**: `"a\x00b"` is a
   2-byte string (`\x00b` is one hex escape); the harness now uses `"a\x00" "b"`.
2. `ex02-ft_strlcpy` — bounded copy returning full src length; size-0 must not write.
   **Known trap**: `strncmp` does NOT stop at a NUL, so the harness compares only the copied
   region via `memcmp(dst, src, copy)`.
3. `ex03-ft_memset` — fill + untouched-region sentinel checks; return-value check.
4. `ex04-ft_strcmp` — signs only; high-bit byte case (`\xff` vs `\x00` ⇒ positive) catches
   signed-char bugs.
5. `ex05-ft_strdup` — **gate**: own malloc/copy loop, must leak nothing and be independent;
   ASan/UBSan enforce heap overflows and leaks (LeakSanitizer reports fail the run).

Validated: reference solutions in `answers/M1-clib/*` all PASS via `forge check`; a broken
`ft_strlen` (overcounts by 1) FAILs on both ways (correct→pass, broken→fail);
`forge list` shows `M1-clib 5/5`.
Readings per exercise are man-page + K&R + TLPI §2.5/§7.1 + cppreference links (see subject.md).

### M2 detail (authored 2026-08-30) — module `M2-procs`, all green
Shape: the *whole program* is student work — one `main.c` + `Makefile` (`all`→`./test`,
fclean/re). No scaffold `main.c` provided. Graded `build`+`quiz`, strict
`-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined`. Determinism rule: no PIDs
in output, parent prints strictly after `waitpid` returns; signals tested in-process via
`raise` (async-signal-safe: handler only sets a `volatile sig_atomic_t` flag, NO printf in
handler — this is also ASan/UBSan-safe).

1. `ex01-fork` — fork + `waitpid` + WIFEXITED/WEXITSTATUS; child exits 7, parent must verify.
2. `ex02-exec` — fork→`execvp`→wait spawn idiom with TWO runs (success: `/bin/echo hello`;
   failure: bogus cmd → `exec failed for …` + parent propagates exit 127, graded via
   `expect_exit`). **Bug caught authoring**: parent must propagate the child's status or the
   fail run exits 0 and the `expect_exit:127` check fails.
3. `ex03-pipe` — child writes 19 bytes, parent reads until EOF; EOF requires closing the
   parent's write end — forgotten close hangs → grader times out (20s) = feedback.
4. `ex04-signal` — `raise(SIGUSR1)` after installing handler; default disposition kills the
   process → grader reports `exit -1 (want 0)`.
5. `ex05-pipeline` — **gate**: two children + pipe + `dup2` + `execlp echo|tr`, parent reaps
   both; two runs (`hello`, `42 school`). Enforces `dup2`/close discipline.

**Grading-method decision**: graded exclusively with proven `build`+`stdout`+`quiz`.
`process`/`net` runners exist (both fixture-green) but cross-process timing/zombie-state
subtleties make graceful-shutdown tests flaky right now; deferring `process`/`net`
curriculum use until M6 (networking) / M8 (switch), where they are the natural fit.

### M3 detail (authored 2026-08-30) — module `M3-memory`, all green
Shape: **implementation-from-file** ("harness shape") for ex01/ex03/ex04/ex05 — provided
header + harness `main.c`, student writes the impl `.c` + `Makefile`; ex02 is *whole-program*
like M2 (observational, needs a `fork` in `main`). Graded `build`+`quiz`, strict
`-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined`. Determinism: no dynamic
`malloc` anywhere — every allocator owns a `static _Alignas(8)` arena/pool/registry, so
"OOM" is a deterministic boundary, all harness assertions are exact and ASan-safe.

The module storyline: `malloc` is not magic — build it four ways; the gate proves you can
build a *bounded* one that refuses instead of dying.

1. `ex01-arena` — bump arena: 512-byte static region, 8-aligned base, offset-only alloc.
   Harness proves exact 512-byte exhaustion (64 allocations), per-block byte tags (no
   overlap), `arena_left()`, reset, and `NULL` on full.
2. `ex02-virtual` — copy-on-write: mmap two anonymous pages (MAP_PRIVATE + MAP_SHARED),
   fork, child writes `B`, parent `waitpid`s and reads both bytes; whole-program. Verdict
   lines pin the CoW behavior; a student on the wrong flag flips exactly one predicate.
3. `ex03-pool` — fixed slot pool: 8×32B static pool, bitset bookkeeping; `NULL` when
   exhausted; `pool_put` must silently reject foreign/NULL pointers. Harness checks 9th-get
   NULL, isolation tags, accounting, invalid-put, and 100-round churn.
4. `ex04-mmap` — mmap allocator: anonymous private mappings, page-rounded, zeroed; static
   live-region registry; `mm_allocated()` exact accounting. Harness checks alignment,
   disjointness, zero-fill, accounting (incl. churn), silent invalid-free, leak-free exit.
5. `ex05-gate` — OOM-safe growable buffer: `_Alignas(8)` static 1024 arena, `cap` starts 64
   and doubles to `BUF_MAX`; `buf_put` refuses `-1` with NO side effects when full; growth
   must preserve all earlier bytes. Broken classes proven to FAIL: ignored refusal
   (`r=0` bug), slot overlap (pool), `MAP_SHARED` on the private page.

**Determinism evidence (both-ways)**: reference solutions all PASS; three broken solutions
proven to FAIL deterministically — pool handing out the same slot (29 duplicate-slot FAIL
lines + churn failure), copy-on-write with `MAP_SHARED` on the private page (`FAIL: private
page shows B (want A)`), gate ignoring refusal (`FAIL refusal (r=0 len=1024)`).

### M4 detail (authored 2026-08-30) — module `M4-concurrency`, all green

Shape: whole-program (student `main.c` + `Makefile`) for ex01–ex04, harness shape for the
ex05 gate (provided `queue.h` + `main.c`, student writes `queue.c`); graded `build`+`quiz`.
**All five runs under ThreadSanitizer**: `"sanitizers": ["thread"]` (compile+link get
`-fsanitize=thread`) and every run command is wrapped `["setarch","x86_64","-R","./test"]`
— stock gcc TSan crashes on this kernel's ASLR (`FATAL: ThreadSanitizer: unexpected memory
mapping`), `setarch -R` disables ASLR and makes it deterministic. Policy: **data race =
build error**, not "usually works" — TSan aborts (exit 66) on the first race, and the
grader's `stderr_ok:false` + `exit` check turns that into FAIL. `-pthread` appears on the
cc/link lines in the student Makefile, never in `CFLAGS`. Determinism: all verdict lines
printed by `main` strictly after `pthread_join`; every aggregate is an exact number.

1. `ex01-join` — join & exit: 5 threads, value `i*2`; thread 2 exits via `pthread_exit`,
   the rest via `return`; `main` joins, prints `thread N returned V` in index order,
   `join errors 0`, exact `sum 30`.
2. `ex02-mutex` — mutex bomb: two threads × 200,000 increments of `volatile long tally`
   under `PTHREAD_MUTEX_INITIALIZER`; exact `total 400000`. A plain `tally++` races → TSan.
3. `ex03-atomic` — atomic counter: same 400,000 total through `atomic_ulong` +
   `atomic_fetch_add`, no lock at all; exact `total 400000`. The lesson: the lock free
   alternative; TSan recognizes `_Atomic` as synchronization, so the correct build is silent.
4. `ex04-rwlock` — readers & writers: 4 readers × 50 + 2 writers × 10 over `{a,b}` with
   invariant `b == 2*a` plus `version`; readers hold `rdlock`, writers `wrlock`; the
   `reads`/`consistent` counters are guarded by their *own* plain mutex (never the rwlock,
   which readers hold concurrently). Exact `consistent 200`, `version 20`.
5. `ex05-gate` — bounded queue: `queue_new(16)` + blocking `queue_push`/`queue_pop` with a
   condvar pair (`not_full`, `not_empty`), the canonical `while(…) pthread_cond_wait`
   loop; harness runs producer(1000)+consumer(1000); expected exact FIFO + `checksum 499500`.

**Grading-method decision**: `build`+`quiz` only (like M1–M3); race detection is done by
TSan inside the normal `build` run, and `process`/`net` drivers remain deferred to M6/M8.

**Determinism evidence (both-ways)**: all five references PASS; six broken classes proven to
FAIL under `setarch x86_64 -R` with strict flags — (a) ex01 wrong exit value (`sum 15
(want 30)`, exit 1), (b) ex02 mutex omitted (TSan `data race`, exit 66), (c) ex03 plain
`long tally++` (TSan `data race`, exit 66), (d) ex04 no locks (TSan `data race`, exit 66),
(e) ex05 mutex-free queue (TSan `data race`, exit 66), (f) ex05 LIFO pop (`fifo violated`,
exit 1). Bonus: calling `pthread_cond_wait(…, NULL)` (forgetting the mutex) is itself a
build error under `-Werror=nonnull` — the platform refuses it before it can ever run.

### M5 detail (authored 2026-08-30) — module `M5-filesio`, all green

Shape: whole-program (student `main.c` + `Makefile`) for ex01–ex04; harness shape for the
ex05 gate (provided `fix.h` + `main.c`, student writes `log.c`); graded `build`+`quiz`
under `-std=gnu11 -Wall -Wextra -Werror` + ASan/UBSan (no TSan here — nothing shared).
ex02 and ex05 are **multi-run** (two fixtures, per-run `expect`); ex04 adds an `artifact`
method pinning `out.txt` to exactly 24 bytes. Determinism: verdicts are exact byte counts
and exact strings printed after every fd is closed and every child is reaped.

1. `ex01-descriptors` — the fd table: stdin/stdout/stderr via `fileno`, then
   `open("/dev/null")` (first → 3), `close` + reopen (reuse → 3), two simultaneous opens
   (3 and 4), lowest-available again. No files beyond `/dev/null`; exact numbers from
   `open`, never hard-coded strings.
2. `ex02-readwrite` — `read()` (7-byte chunks) vs `fread()` (11-byte chunks) over the same
   fixture with matching totals + byte checksums; runs on `input.txt` (44 B) and
   `second.txt` (11 B). Teaches short reads/EOF.
3. `ex03-redirect` — `dup2` + fork: child repoints fd 1 at `out.txt`, prints, **flushes**,
   `_exit`s; parent reaps it, reads back, prints `captured: <bytes>`. The stdio-flush trap
   is the lesson (`_exit` skips it).
4. `ex04-tee` — stdin → stdout **and** `out.txt` in one `read`/`write` loop; `wrote`/
   `readback` must agree; `artifact` asserts the 24-byte file. Printed verdicts come after
   the copied bytes.
5. `ex05-gate` — log parser over an fd: `parse_log(int fd, Tokens *t)` reads via `read()`,
   buffers across partial reads, splits lines on `=`, skips blanks, records names in order,
   returns `-1` on malformed/overflow. Fixtures `config.txt`, `apps.txt` → exact token
   lists.

**Grading-method decision**: `build`+`quiz` (+ `artifact` on ex04). The `process`/`net`
drivers stay deferred to M6/M8; ex03's fork+dup2 stays deterministic via `waitpid`
because the child is reaped before the parent reads.

**Quiz parsing trap found**: `parseKV` splits a quiz line at the **first** `:` — question
text containing an internal colon (e.g. `True or false: …`) breaks the key. Lesson for
future modules: quiz question text must contain no `:` (and no `=` outside the delimiter).
ex04 Q2 reworded from `True or false: …` to `Read() and write() return the number of bytes
actually transferred.`

**Determinism evidence (both-ways)**: all five references PASS; five broken classes proven
to FAIL — (a) ex01 reused one fd number for both "opens" (`two opens 3 and 3`, diff FAIL),
(b) ex02 counted a full buffer on every `read` (`read 49` vs `fread 44` + `mismatch`,
exit 1), (c) ex03 forgot the flush (`captured: ` empty, exit 1 — the `_exit` trap the
reference demonstrates correctly), (d) ex04 skipped the file write (`wrote 24` /
`readback 0` + `mismatch`, and `out.txt` 0 bytes fails the artifact min_size; note the
broken proof must run with stdout separated from `out.txt`, or the shell-redirect noise
masks it), (e) ex05 split on `:` instead of `=` (`parse failed` on stderr, exit 1).

### M6 detail (authored 2026-08-30) — module `M6-networking`, all green

Shape: whole-program (student `main.c` + `Makefile`) for ex01–ex05; graded
`build`+`net`+`quiz` under `-std=gnu11 -Wall -Wextra -Werror` + ASan/UBSan. **First
curriculum use of the `net` grader**: the runner spawns the student's `./test`, retry-dials
until it accepts (`TARGETPORT`/`TARGETHOST` injected), and plays scripted send/expect steps
with a bounded deadline — verdicts are byte-exact transcript matches, not timing probes, so
the both-ways rule holds. Server stdout/stderr is ignored (the grader closes those); the
only graded voice is the socket.

1. `ex01-echo` — the five-call dance `socket → setsockopt(SO_REUSEADDR) → bind → listen →
   accept` on the injected loopback port, then an echo read/write loop per connection,
   forever (`SIGPIPE` ignored). Scripted: `hello\n`↔`hello\n`, `ping\n`↔`ping\n`.
2. `ex02-linechat` — per-connection line framing: reply `<lineno>:<line>\n` counting
   complete lines within one connection. Two **connections** in one run (conn 2 restarts at
   `1:`) prove the accept loop survives EOF and state is per-connection.
3. `ex03-http` — parse until `\r\n\r\n`, route `GET /ping` → exact
   `HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\npong\n`, anything else → exact
   `HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n`; two requests on one connection
   (loop, don't close after the first response).
4. `ex04-timeout` — `SO_RCVTIMEO` 400 ms on the listening fd: echoed complete lines, but a
   client that sends `half` (no newline) and stalls gets `TIMEOUT\n` before the grader's
   1200 ms sleep ends; server re-accepts afterwards.
5. `ex05-gate` — stateful protocol server: greeting `HELLO <n>\n` **sent on accept**, then
   `PING`→`PONG`, `ECHO <t>`→`ECHO <t>`, else `BAD`; the connection counter is shared
   across connections of one process (conn 2 sees `HELLO 2`) — the skeleton of the M8
   switch that grades later modules.

**Grading-engine change**: `NetSpec` gained `connections: [{steps: [...]}, …]` — each
session dials a **fresh TCP connection against the same server process** (the earlier draft
used two `net` methods, which `normalize()` rightly rejects as duplicate types). The
single-`steps` form is unchanged; the existing `net/pass|fail/ex01-*` fixtures still pass.
New fixtures `net/pass/ex02-sessions` and `net/fail/ex02-once` (server that serves exactly
one connection then exits) lock the path in. Selftest is now 23 fixtures.

**Determinism evidence (both-ways)**: all five references PASS; five broken classes proven
to FAIL — (a) ex01 binds a hard-coded port (grader's dial to `TARGETPORT` never connects:
`start` FAIL), (b) ex02 serves one connection then exits (`conn 2 step 1` FAIL), (c) ex03
answers 404 with the 200 body (byte mismatch on step 2), (d) ex04 never sets `SO_RCVTIMEO`
(step 2 never receives `TIMEOUT\n`; deadline FAIL), (e) ex05 only speaks after the first
request line (greeting-with-`""`-send step times out: `conn 1 step 1` FAIL).

### M7 detail (authored 2026-08-30) — module `M7-resilience`, all green

Shape: **first two-language mix**. ex01 is Python — the **first `lang: python` exercise and
the first curriculum use of the `stdout` grader method** (`python3 main.py`, `normalize:
strip`, stderr must stay empty). ex02–ex05 are Go, built by the module's **shared Makefile**
(`PATH := $(HOME)/.local/go/bin:$(PATH)`, `go build -o test .`, `all`/`fclean`/`re`). M7
rule: **nothing printed may depend on wall-clock time** — retries/sleeps are real but never
printed; ex02's breaker runs against an injected fake clock; transcripts are event-driven,
so both-ways determinism holds.

1. `ex01-backoff` (Python, `stdout`) — retry-connect to a flaky peer: first two attempts end
   in EOF/RST, the third answers. `attempt()` classifies a reset/close **before a response**
   as a failed attempt; the driver loops 3 tries with exponential backoff (`0.01 * 2^(n-1)`,
   slept but never printed). *Bring-up trap*: the flaky server closes with the client's PING
   still unread → the peer sends RST, so `recv` raises `ConnectionResetError`, not a clean
   EOF — both classes must count as a failed attempt.
2. `ex02-breaker` (Go, harness-shape) — provided `main.go` drives 12 scripted calls
   (`results = [ok,ok,fail,fail,fail,ok,ok,ok]`, consumed only by allowed calls so
   fast-fails never shift the transcript) against the student's `breaker.go`
   (`NewBreaker(2,3,now)`): Closed / Open / HalfOpen, trip after 2 consecutive failures,
   cooldown 3 ticks, one-probe-decides, fast-fail while open. *Bring-up fix*: probing must
   resume at the cooldown boundary inclusively (`!now.Before(openUntil)`), else the probe
   drifts one call late.
3. `ex03-health` (Go, whole-program, `net` grader) — health-check service with **cross-
   connection state**: 5-session dialogue (port 17301) = `/health`→200 `ok`, `/down`→200
   then strobes unhealthy, `/health`→503 with `Content-Length: 0`, `/up`→200, `/health`→200.
   First `net` exercise where later connections assert state set by earlier ones.
4. `ex04-shutdown` (Go, whole-program) — worker completes 5 tasks (ordered by an
   unbuffered channel), then the program installs a SIGTERM handler, SIGTERMs itself,
   prints `received SIGTERM` / `shutdown complete` and exits 0. Killing-without-a-handler
   must be out: the handled-signal transcript + exit 0 is the graded contract.
5. `ex05-gate` (Go, whole-program) — mini supervisor: three workers a(3)/b(3)/c(2),
   round-robin steps; b's task 2 crashes on its first attempt, the watchdog restarts the
   worker **preserving completed progress** and the retry succeeds; `maxCrash = 2`
   consecutive failures → `giving up on worker <name>` + exit 1 (safety valve, untriggered).
   Ends `all workers complete` / `shutdown complete`.

**Determinism evidence (both-ways)**: references all PASS; five broken classes proven FAIL —
(a) ex01 no-retry/single-attempt (missing `attempt 2/3` lines), (b) ex02 breaker that never
trips (never prints `open`/`half-open` lines), (c) ex03 `/down` never flips state (conn 3
`/health` still returns 200), (d) ex04 no graceful path (missing `received SIGTERM`), (e)
ex05 flaky task fails on *every* retry (restart budget exhausted → `giving up on worker b`,
exit 1). Hygiene lesson: the both-ways loop must restore references *immediately* after each
broken run — the swap-grader-again-restore script overwrote answers' references and they had
to be rebuilt from spec (verified PASS again). `forge score` → `M7-resilience 5/5`,
`TOTAL 36/40`; selftest 23 fixtures, 0 mismatches; no orphaned `./test` after the ex03 net
run.

### M8 detail (authored 2026-08-30) — module `M8-messaging`, all green

Shape: full Go. ex01/ex02 are harness-shape (provided `main.go` + `go.mod` +
`Makefile`; student writes `frame.go` / `queue.go`), graded `build`+`stdout`+`quiz`.
ex03–ex05 are whole-program: the student writes **`switch.go`**; the scaffold ships a
provided **backend** (line server replying `R<seq> <content>` per connection, built as
`./svc`) the switch spawns itself. M8 rule: nothing printed may depend on wall-clock
time — every fault effect is proven **at the data level**.

1. `ex01-framing` — length-prefixed frame codec (4-byte big-endian length + payload;
   `PutFrame` loops; `GetFrame` uses `io.ReadFull` so one-byte-at-a-time readers can't
   break it; clean EOF vs `io.ErrUnexpectedEOF` mid-frame). Driver proves roundtrip,
   one-byte-chunk decode, and truncated-frame error in one transcript.
2. `ex02-orders` — ordered delivery: per-partition FIFO **with gap holding** (a
   partition's later offsets wait for their missing predecessor), earliest-arrival wins,
   gap-guard error instead of silent reorder. The tangled arrival order
   `a1, b0, a0, b1, a2, b2` has exactly one legal delivery order: `b0|a0|a1|b1|a2|b2`.
3. `ex03-switch` — the relay core: bind `TARGETPORT`; spawn `./svc` on an ephemeral
   loopback port; dial a **fresh** backend connection per client; relay newline-framed
   bytes both directions at once; on client disconnect tear the pair down and keep
   serving. Graded net: conn 1 `ping/alpha/beta` → `R1..R3`, conn 2 → `R1 one` (reset).
4. `ex04-switchfaults` — fault plane from `faults.txt` (`dup:N`, `drop:N`, `hold:N`;
   per-connection frame numbers; `TARGETFAULTS` env overrides the path): dup → one send
   gets two replies; drop → the next relayed frame answers with the *earlier* R number;
   hold → the frame goes after the following frame, so reply order swaps. The backend's
   numbered-and-echoed replies (`R5 M3`, not numbers alone) make every scenario
   byte-distinct.
5. `ex05-gate` — mixed faults in one stream (`dup:1 drop:2 hold:4`), a mid-scenario
   disconnect, then a fresh connection with faults re-applied and numbering reset — the
   switch in the shape every later grader will use as its injector.

**Engine note: zero grader changes.** `NetSpec.StartEnv` already forwards
`TARGETPORT`/`TARGETHOST`/`TARGETFAULTS`; the switch spawns its own backend so the
runner stays single-process; M6's `connections` sessions carry the multi-conn proofs.
Three findings: (a) **make binary/source-dir collision** — a `backend:` binary target
with a `backend/` source dir never runs (make sees the target file as up-to-date);
binary renamed `svc`; (b) **exact-bytes union wants** — when a dup produces two replies,
a single step expects *both* (`R1 dupme\nR2 dupme\n`), so `netTransaction`'s exact match
can never be over-read by coalesced packets; (c) spec hygiene — frame indexes must be
recounted against each connection's own numbering when handwriting multi-conn scenarios
(two early drafts miscounted `hold`/`drop` frames and failed the check; corrected in the
exercise.json, not the reference).

**Determinism evidence (both-ways)**: references all PASS; five broken classes proven
FAIL — ex01 single-`Read` codec (garbage length/payload → wrong stdout), ex02 arrival-
order replay (`a1|b0|…`), ex03 one shared backend connection (conn 2 gets `R4` instead
of `R1`), ex04 hold-as-drop (released reply never arrives), ex05 hold released *before*
the following frame (replies `R4 M3 / R5 M4`, swapped). References restored immediately
after each broken run (M7 hygiene, now standard). New fixtures `net/pass/ex05-child` +
`net/fail/ex05-norelay` (server spawns a worker subprocess; the fail variant never
relays) — selftest now **25 fixtures, 0 mismatches**. `forge score` → `M8-messaging
5/5`, `TOTAL 41/45`; no orphaned `switch`/`svc` processes after the multi-conn runs.

### M9 detail (authored 2026-08-31) — module `M9-time`, all green
Five exercises on **logical time**, every one transcript-driven with no wall-clock
dependence. The module rules: nothing printed or compared may read the clock; every
exercise is independently byte-deterministic.

1. `ex01-clock` — harness-shape Go; student writes `clock.go` with `Lamport{Tick,Now,Add}`.
   The receive rule (`c = max(c, s) + 1`) is the point: p1's stamp must jump past p0's
   message (`p1 recv m0 t3` after `t2`), and p0 jumps past p1's later stamp on receive
   (`p0 recv m1 t6` after being at `e2 t3`) — the tell that local and remote events
   interleave correctly.
2. `ex02-vector` — student writes `vector.go` `VC{Local,Merge,Now}` over `[2]int`. The
   two subtle rules: a receive merges the vector **then** ticks its own component
   (`p1 recv m0 (2,1)`, not `(2,0)`), and a send *is* a local event (`p0 send m0 (2,0)`).
   Broken vector prints a visibly different stamp grid on the first receive.
3. `ex03-causal` — `causal.go` `CausallyBefore` (strict componentwise `<=`) and
   `Concurrent` (neither direction). The trap: stamps that are *equal* are NOT
   happened-before (`e0 -> e3: false`), and `e2`'s own-component-heavy stamp is
   concurrent with everything — a scalar counter cannot express any of this.
4. `ex04-total` — `total.go` `TotalOrder` returns indices sorted by `(lamport, pid)`
   via `sort.SliceStable`. Two Lamport ties (`e1/e4` at 1, `e0/e3` at 2) are broken by
   pid so the sequence `e1 e4 e0 e3 e2 e5` is fully deterministic. pid-major ordering
   is the classic bug → `e1 e0 e2 e4 e3 e5`.
5. `ex05-sequencer` — whole-program `srv.go`, TCP on `TARGETPORT`; the grader dials two
   connections and requires **byte-exact** replies. Each frame line `<content> <stamp>`
   gets `seq = max(seq, stamp) + 1`; the counter is **process-lifetime** (never reset
   per conn). Conn 2 opens at `R12` (after `R11`), so a per-connection reset is caught.
   This is the single-writer total-order primitive replicated systems build on.

**Grade shape**: ex01–ex04 `build`+`stdout`+`quiz` (shared Makefile, `go build -o test .`,
`PATH`-prefixed); ex05 `build`+`net`+`quiz`. No grader changes needed.

**Determinism evidence (both-ways)**: references PASS; five broken classes proven FAIL —
ex01 `Add` ignoring the incoming stamp (`p1 recv m0 t0`), ex02 `Merge` ignoring the
vector (`p1 recv m0 (0,1)`), ex03 non-strict happened-before (`e0 -> e3: true`), ex04
pid-major order, ex05 no-stamp-merge sequencer (conn 2 `R2 b 9` not `R12`). Quizzes:
manual actuals proven for ex01 (`explicit`, `Timestamp`, `happened-before`),
ex02/ex03/ex05 defaults (`yes`, `no`, `persist`, `a total order`). References restored
immediately after each broken run. New fixtures `net/pass/ex05-clock` + `net/fail/
ex05-reset` → selftest **27 fixtures, 0 mismatches**. `forge score` → `M9-time 5/5`,
`TOTAL 47/50`; no orphaned `srv` processes.

### M10 detail (authored 2026-08-31) — module `M10-commit`, all green
Turning M9's ordering into agreement. Same discipline: transcript-driven, byte-
deterministic, never read the wall clock.

1. `ex01-coordinator` — harness-shape Go; student writes `coordinator.go`
   `Coordinator{Run(vote Voter) string}`. The rule: `commit` iff EVERY participant
   votes commit, else `abort`. `main.go` prints the phase-one prepare lines from
   `votes`, then the decision and per-participant actions. The tell: with votes
   `yes,yes,no` the decision is `abort` and EVERYONE is aborted — partial commit is
   the banned bug.
2. `ex02-vote` — `vote.go` `VoteFor(canCommit bool) string`. Trivial branch, but the
   quiz carries the real content: a commit vote is a binding promise, and an
   abort-after-commit is forbidden because the coordinator may already have told other
   participants to commit. Broken = always commit (or always abort) flips the
   transcript.
3. `ex03-elect` — `elect.go` `Elect([]Candidate) *Candidate`. Bully election: newest
   logical stamp wins, ties broken by higher id, result independent of input order
   (`e0` newest-stamp wins; `e1` id5 at an OLDER stamp loses to id3; `e2` tie → higher
   id). Id-major is the classic bug → `e1: leader 5`.
4. `ex04-quorum` — `quorum.go` `QuorumSize(n) = n/2+1` and `HasQuorum(n, yes)`. The
   point: a 2-node group needs both (quorum 2), a 4-node group needs 3 (2 is a tie,
   not a majority). `floor(n/2)` quorum is the classic bug → 4-group accepts 2.
5. `ex05-transaction` — whole-program `gate.go`, TCP on `TARGETPORT`; the grader dials
   two connections and requires byte-exact replies. Participants prepared-votes come
   from a scaffold `votes.txt`; `tx <id>` replies `COMMIT` only when the id is
   present, else `ABORT`. The gateway keeps the prepared-set for the process lifetime
   (an absent id is an abort even on a later connection).

**Grade shape**: ex01–ex04 `build`+`stdout`+`quiz` (shared Makefile, `go build -o test .`,
`PATH`-prefixed); ex05 `build`+`net`+`quiz` (separate port 17410, starts `./gate`). No
grader changes.

**Determinism evidence (both-ways)**: references PASS; five broken classes proven FAIL —
ex01 majority-commit (`DECISION commit` when one participant said no), ex02 always-commit
(`t1: commit`), ex03 id-major (`e1: leader 5`), ex04 `n/2` quorum (`n=4 quorum 2` +
`has 2?` wrong), ex05 commit-every-tx (`conn 1 step 2` got `COMMIT` for `ABORT`).
References restored immediately after each broken run. New fixtures `net/pass/ex05-tx` +
`net/fail/ex05-txabort` → selftest **29 fixtures, 0 mismatches**. `forge score` →
`M10-commit 5/5`, `TOTAL 52/55`; no orphaned `gate` processes.

### M11 detail (authored 2026-08-31) — module `M11-log`, all green
The replica tier: a replicated log and the consistency rule that makes it usable.
Same discipline — transcript-driven, byte-deterministic, never read the wall clock.

1. `ex01-append` — harness-shape Go; `log.go` `Log{entries}`, `Append` returns the
   index, `Len`, `All` (defensive copy). Append is a tail push, never a head insert;
   `all=set|add|del` in append order is the tell.
2. `ex02-committed` — `log.go` `NewLog()`, `CommitIndex()`, `Commit(idx)` (pure setter
   returning the new index, never lowers), `Committed()` returns ONLY the durable
   prefix up to the commit index. Commit starts at -1 (nothing durable). Leaking the
   uncommitted tail to `Committed()` is the banned bug.
3. `ex03-index` — `index.go` pure arithmetic on `State{Commit,Len}` (ex02's fields
   lifted out): `NextIndex` = `Len` (never `Len-1`), `Uncommitted` = `Len-Commit-1`.
   `next(c=2,l=3)=3` with `uncommitted=0` is the tell; `len-1` mixes up index and count.
4. `ex04-snapshot` — `log.go` `Log{entries,base}` compaction: `Snapshot(upto)` drops
   the prefix at/below `upto` and returns the count removed, keeping a `base` offset
   so remaining entries keep their original public indexes (`c` stays index 2 after
   dropping 0 and 1; next append lands at 4). Renumbering the tail after compaction is
   the banned bug (leads to index collision/corruption).
5. `ex05-replica` — whole-program `replica.go`, TCP server on `TARGETPORT`; one log for
   the process lifetime shared across connections. Commands: `append <cmd>` → `ok <idx>`,
   `commit <idx>` → `commit <idx>`, `read` → committed prefix joined by `|`. `read` never
   reveals an entry beyond the commit index: after appends x,y then `commit 1`, conn 2
   appending z (index 2, uncommitted) must still `read` only `x|y`.

**Grade shape**: ex01–ex04 `build`+`stdout`+`quiz` (shared Makefile, `go build -o test .`,
`PATH`-prefixed); ex05 `build`+`net`+`quiz` (port 17415, starts `./replica`), a single net
method with two connections. No grader changes.

**Determinism evidence (both-ways)**: references PASS; five broken classes proven FAIL —
ex01 prepend-not-append (`all=del|add|set`), ex02 `Committed` leaks uncommitted
(`committed=a|b|c` before any commit), ex03 `NextIndex=Len-1`/wrong `Uncommitted`
(`next(c=2,l=3)=2,u=1`), ex04 renumbers-after-compaction (no base → `all=0:c|1:d`), ex05
read leaks uncommitted (`conn 2 step 2` got `x|y|z` not `x|y`). References restored
immediately after each broken run. New fixtures `net/pass/ex05-replica` +
`net/fail/ex05-replicaleak` → selftest **31 fixtures, 0 mismatches**. `forge score` →
`M11-log 5/5`, `TOTAL 57/60`; no orphaned `replica` processes.

### M12 detail (authored 2026-08-31) — module `M12-raft`, all green
The flagship capstone, as the blueprint intends: Raft's core safety mechanics.
M9 gave ordering, M10 agreement, M11 a single replicated log; M12 is how a
leader drives ONE consistent log across a cluster safely. All no-wall-clock,
byte-deterministic.

1. `ex01-term` — `node.go` monotonic term. `NewNode()` (term 0, no vote),
   `Term()`, `Voted()`, `SeeTerm(t)` (steps down to follower + resets the vote
   ONLY on a strictly higher term), `StartElection()` (term++, vote self).
   `see 3` at term 6 is ignored (still 6, vote kept) — the tell. A term that
   regresses, or a vote reset on a stale/equal term, is the bug.
2. `ex02-vote` — `vote.go` the up-to-date-log vote rule. `LogInfo{Term,Idx}`,
   `Voter{term,voted}`, `RequestVote(candTerm, cand, mine)`. Grants only when
   term not stale AND not already voted AND cand log at least as up-to-date
   (higher last-term wins; equal term → higher/equal last-index). case d
   (`cand term1 idx5 vs mine term0 idx8` → true) — the higher term wins even
   though the candidate has fewer entries; voting on length alone, or voting
   twice, is the bug.
3. `ex03-logmatch` — `match.go` the matching-prefix AppendEntries check.
   `AppendEntries(log, prevIndex, prevTerm, ent)` appends only when
   `log[prevIndex].Term == prevTerm`; else returns nil (REJECT). `mismatch
   prev(0,3) -> REJECT` — appending anyway is the bug.
4. `ex04-election` — `elect.go` strict-majority quorum. `Quorum(n)=n/2+1`,
   `HasMajority(n,yes)=yes>=Quorum(n)`. `n=2 quorum=2`, `n=4 quorum=3`; a
   `floor(n/2)` quorum that lets a 2-node win on 1 (or a 4-node on 2) is the
   bug.
5. `ex05-gatenode` — whole-program `node.go`, TCP server on `TARGETPORT`;
   one node state (monotonic term + committed log) for the process lifetime.
   Commands: `see <t>` → `term <t>` (never regresses), `append <cmd>` → `ok
   <idx>` or `not leader` while term 0, `read` → committed prefix joined by
   `|`. conn2 shares the state and sees only committed entries.

**Grade shape**: ex01–ex04 `build`+`stdout`+`quiz` (shared Makefile, `go build -o test .`,
`PATH`-prefixed); ex05 `build`+`net`+`quiz` (port 17420, starts `./node`), a single net
method with two connections. No grader changes.

**Determinism evidence (both-ways)**: references PASS; five broken classes proven FAIL —
ex01 term-regress (`see3 -> term=3 voted=false` after elect), ex02 vote-on-length-only
(cases c/d flip to false), ex03 append-despite-mismatch (`mismatch ... -> b2 c3 d4 e5`
not REJECT), ex04 `n/2` quorum (`n=2 quorum=1 has1=true`), ex05 stale-see-regresses
(`conn 1` `see 2` got `term 2` not `term 3`). References restored immediately after each
broken run. New fixtures `net/pass/ex05-gatenode` + `net/fail/ex05-gateregress` →
selftest **33 fixtures, 0 mismatches**. `forge score` → `M12-raft 5/5`, `TOTAL 62/65`;
no orphaned `node` processes.

### M13 detail (authored 2026-08-31) — module `M13-shard`, all green
Spreading keys across a cluster: placement mechanics and rebalancing. All
no-wall-clock, byte-deterministic.

1. `ex01-slot` — `slot.go` fixed-slot assignment. `Hash` = FNV-1a 32 (offset
   2166136261, prime 16777619), `SlotOf = int(Hash % nSlots)`. The transcript
   prints five fruit keys; `banana`/`cherry` collide on slot 0 and
   `apple`/`elder` on slot 3 — two keys hashing to the same slot always land on
   the same shard. A constant/`%0`-style broken hash flattens every slot to 0.
2. `ex02-range` — `range.go` key-range sharding. `RangeShard(key, boundaries)`
   = smallest `i` with `key <= boundaries[i]`, else `len(boundaries)`. With
   boundaries `k p u`: `apple→0, kiwi→1, k→0, mango→1, pear→2, u→2, zebra→3`.
   The on-boundary keys `k`→0 and `u`→2 lock the **inclusive** semantics: a
   `<`-instead-of-`<=` comparator routes `k`→1 and `u`→3 (proven FAIL).
3. `ex03-ring` — `ring.go` consistent hashing. `NodeFor(keyHash, nodes)` = the
   first node at or after the key, wrapping to node 0. `key 50`/`key 600` wrap
   to node 0; `key 100`/`key 300` are exact (inclusive) hits. A lookup using
   `>` or that fails to wrap routes `key 600`→node 2 (proven FAIL).
4. `ex04-rebalance` — `rebalance.go` move counting. `OwnerHash` (first node at
   or after, wrapping) + `Moved(old, new, keys)` counting keys whose owning node
   hash changed. Add 200 moves only key `150` (1); add 700 moves only `550` (1);
   remove 300 moves `150`+`250` (2). Counting by index instead of by node hash
   over-reports (proven FAIL: `add 200 -> 4 moved`).
5. `ex05-shardgate` — whole-program `gate.go`, TCP server on `TARGETPORT`; three
   isolated key→value shards with sorted boundaries `m`,`t`. `set <k> <v>` →
   `set <n>`, `get <k>` → `<n>:<v>` (or `<n>:?`). apple→0, mango→1, zebra→2,
   melon→1 (absent → `1:?`). Shard state persists across connections. Routing a
   key to the wrong shard is the bug (proven FAIL: everything → 0).

**Grade shape**: ex01–ex04 `build`+`stdout`+`quiz` (shared Makefile, `go build -o test .`,
`PATH`-prefixed); ex05 `build`+`net`+`quiz` (port 17425, starts `./gate`), a single net
method with two connections. No grader changes.

**Determinism evidence (both-ways)**: references PASS; five broken classes proven FAIL —
ex01 constant-hash (`@4 -> 0` for all), ex02 exclusive-`<` boundary (`k->1 u->3`), ex03
no-wrap (`key 600 -> node 2`), ex04 count-by-index (`add 200 -> 4 moved`), ex05 wrong-
shard-routing (`conn 1 step 2` got `0:?`). References restored immediately after each
broken run. New fixtures `net/pass/ex05-shardgate` + `net/fail/ex05-shardwrong` →
selftest **35 fixtures, 0 mismatches**. `forge score` → `M13-shard 5/5`, `TOTAL 67/70`;
no orphaned `gate` processes.

### M15 detail (authored 2026-08-31) — module `M15-tx`, all green
Transactional ACID primitives and fault injection (chaos). Determinism rule: no wall clock, tick-based events.
1. `ex01-transaction` — `store.go` transaction begin/commit/rollback harness; snapshot isolation + atomic application.
2. `ex02-replay` — `log.go` write-ahead log append and replay, last-writer-wins state recovery.
3. `ex03-conflict` — `conflict.go` detection/resolution for lost-update; my-write wins on merge.
4. `ex04-drill` — `drill.go` deterministic fault-injection: operations succeed/fail based on tick-based node availability.
5. `ex05-txstore` — `txstore.go` TCP transactional KV-store gateway; multi-conn isolation.

### M16 detail (authored 2026-09-02) — module `M16-storage`, all green
Disk-based persistent LSM-tree KV engine. Determinism rule: no wall clock; explicit durations/steps.
1. `ex01-memtable` — `memtable.go` in-memory put/get/All write buffer; keys kept sorted lexicographically.
2. `ex02-sstable` — `sstable.go` flush to immutable on-disk sorted file with a trailing sparse index (every 3 records); lookup via index + bounded scan.
3. `ex03-wal` — `wal.go` append-only length-prefixed write-ahead log; replay rebuilds sorted state.
4. `ex04-compact` — `compact.go` k-way merge, newest run wins, tombstone (`Val==""`) drops keys.
5. `ex05-persistent` — **Mini-Capstone:** `server.go` TCP LSM-tree gateway (`persistent-kv`) with WAL-backed shared persistent state; PUT/GET/DEL line protocol; `net` grader (2 conns). Both-ways fixtures added → selftest 41.

### M17 detail (authored 2026-09-02) — module `M17-observability`, all green
Observability for distributed systems. Determinism rule: no wall clock; explicit durations.
1. `ex01-metrics` — `metrics.go` concurrency-safe (mutex) counter registry with sorted snapshot.
2. `ex02-tracing` — `trace.go` span tree with explicit durations; depth-first indented render.
3. `ex03-spans` — `spans.go` spans with attributes; FindAll (depth-first) + MaxDur over subtrees.
4. `ex04-summaries` — `summary.go` count/sum/max/p50/p95 via ceil-index quantiles, deterministic.
5. `ex05-telemetry` — **Mini-Capstone:** `server.go` TCP telemetry gateway exposing the metric registry; INC/GET/SNAP line protocol + per-connection `conn_total`; `net` grader (2 conns). Both-ways fixtures added → selftest 41.

---

## 11. Capstones & Gates (12 graded artifacts, ranked)

Learner builds big standalone artifacts; each is graded by scenario/fault drivers + lincheck typically.

1. **Raft KV replica set** (flagship; M11)
2. **Chaos KV** (chaos-injecting KV under the switch; M14)
3. **Sharded KV** (hash + range shards; M12)
4. **Lock service** (distributed mutex on a replicated log; M10/11)
5. … remaining gates to be detailed as curriculum is written (see checklist below).

Gate checklist to flesh out while writing M8–M14: repro from a clean checkout; run under switch with partitions;
crash recovery; determinism; both-ways fixtures.

---

## 12. Reading Library Strategy

- C / OS foundations: man pages + "The Linux Programming Interface" (chapter-precise pointers).
- Concurrency: "The Little Book of Semaphores", cppreference atomics notes.
- Distributed: "Designing Data-Intensive Applications", Raft paper (raft.github.io), "Patterns of
  Distributed Systems", Herlihy & Wing for linearizability.
- Every exercise references exact chapters/sections; auto-graded checkpoint via `quiz`.

---

## 13. Command Cheat-Sheet

```sh
export PATH="$HOME/.local/go/bin:$PATH"          # if PATH lost in new shell
cd /home/salmane/GolandProjects/systems-forge

make setup          # builds bin/forge            (todo)
./bin/forge list    # modules + exercises + progress
./bin/forge show M2-ex03   # subject.md for an exercise
./bin/forge check M2-ex03  # run grader (exit 1 on fail)
./bin/forge score   # progress per module
./bin/forge selftest # grader fixtures            → must be "0 mismatches"
./bin/forge         # TUI
```

Must run from repo root (engine resolves `subjects/`, `tools/fixtures/`, `answers/`, `progress.db`
relative to CWD). `go build ./...` runs inside `forge/`. Root Makefile targets (§16) remove the foot-guns.

---

## 14. Log / Changelog

- **2026-09-02** **Task B — portfolio / "so what" layer** — made the portfolio promise
  concrete at the module level:
  - Authored the three **missing** `module.md` files: **M15** (Transactions & Chaos,
    ex01–ex06), **M16** (Storage Engines, LSM ex01–ex05), **M17** (Observability,
    ex01–ex05) — intro, milestones, deterministic rule, and the so-what section.
  - Added a `## So what? (interview / portfolio)` section to **all 18 module.md** files:
    job-relevance, 4 interview questions, and a named portfolio artifact each (M12-ex06
    raft cluster, M15-ex06 durable txstore, M16-ex05 persistent-kv, M17-ex05 telemetry,
    plus M0–M14 gates).
  - M12 module.md milestones updated to include the new ex06 cluster. Docs only — grading
    untouched; `forge score` still 92/92. See tracker item #31.
- **2026-09-02** **Task A — Micro-App capstones → TOTAL 92/92** — shipped two new `ex06`
  capstone exercises (no rewrite of the working ex05 gates), turning two flagship modules
  into tangible portfolio artifacts:
  - **M12-ex06 `raft cluster`** — a **real 3-node raft cluster** in one binary: peers
    a/b/c each on its own TCP endpoint (`TARGETPORT+1/+2/+3`), RequestVote/AppendEntries
    exchanged **over real sockets**, leader commits on strict majority (2 of 3), node steps
    down on a higher term. Client-driven, no wall clock. `log b`/`log c` matching the
    leader's `read` (`x|y`) proves real replication to followers. → **M12 6/6**.
    Root: `subjects/M12-raft/ex06-cluster/` (`cluster.go` byte-rewritten, deterministic).
  - **M15-ex06 `durable txstore`** — the **durable transactional KV Micro-App** uniting M15
    transactions with M16 WAL durability: every committed `set` is applied *and* appended to
    a write-ahead log; `wal` exposes the ordered durable record; `begin`/`commit`/`rollback`
    stage atomically and a rollback leaves the WAL untouched; state persists across fresh
    connections. → **M15 6/6**. Root: `subjects/M15-tx/ex06-durabletxstore/` (`durable.go`).
  - Both-ways net fixtures: `net/fail/ex06-nodewnly` (entries never replicated to followers →
    `log b`/`log c` diverge) and `net/fail/ex06-walded` (writes never reach the WAL → `wal`
    lies empty). Selftest → **45 fixtures, 0 mismatches**.
  - `forge score` → M12 6/6, M15 6/6, **TOTAL 92/92 (0 remaining)**; no orphaned servers.
    See tracker item #30.
- **2026-09-02** **Gap-closure → TOTAL 90/90** — shipped the four unfinished exercises
  (M0-ex03/04/05, M15-ex05), the last `[F]`s, closing the mission-wording gap ("completeness
  over speed / never half-arsed"). Fixes:
  - M0-ex03: reference `answers/check.py` + `quiz.txt` → PASS (was missing the deliverable).
  - M0-ex04: reference `answers/Makefile` + fixed `answers/micro.c` (write before `free`) +
    `quiz.txt` → PASS; `subjects/` keeps the planted-bug `micro.c` as the starting scaffold.
  - M0-ex05: `answers/mini_forge.sh` + `quiz.txt`; **fixed a real spec bug** — subject
    demanded `set -euo pipefail` under the grader's `sh` (dash has no `pipefail`) → made it
    POSIX `set -eu` + guard each step, corrected the inline Ground-Truth pipe note, and
    repaired `sample/expected.txt` (missing trailing newline: `gate ok` vs `gate ok\n`). PASS.
  - M15-ex05: authored `subject.md` (was absent), rewrote `txstore.go` into a coherent line
    protocol (`set/get/begin/commit/rollback/status/drop/list`; per-connection staged writes;
    shared committed store; `get`→`nil` on absent) and made `exercise.json` self-consistent
    (was contradictory). Added both-ways net fixtures `net/{pass,fail}/ex05-txstore*`
    (17640/17641; fail = uncommitted staged writes leak across connections). PASS.
  Result: **M0 5/5, M15 5/5, TOTAL 90/90 (0 remaining)**; `make selfcheck` OK; selftest
  **43 fixtures, 0 mismatches** (was 41); no orphaned servers. See tracker item #23.
- **2026-09-02** **M16 "Storage Engines" + M17 "Observability" authored + validated** —
  `subjects/M16-storage/ex01–ex05` (memtable, sstable, wal, compact, persistent TCP gateway)
  and `subjects/M17-observability/ex01–ex05` (metrics, tracing, spans, summaries, telemetry
  TCP gateway). Fully Go, all **no-wall-clock** (explicit durations/steps, byte-deterministic),
  harness-shape `build`+`stdout`+`quiz` (ex01–ex04) and `net` (ex05). Both-ways proven: four
  new net fixtures — `net/{pass,fail}/ex05-persistent*` (M16; fail = per-connection fresh
  store) + `net/{pass,fail}/ex05-telemetry*` (M17; fail = per-connection fresh registry,
  no conn_total) → selftest **41 fixtures, 0 mismatches**. `forge score` → M16 5/5, M17 5/5,
  `TOTAL 86/90` (M0-ex03/04/05 + M15-ex05 shipping deferred to 2026-09-02 session); no orphaned
  servers. See §10 M16/M17 detail and the tracker items 21/22.
- **2026-08-30** Toolchain verified; Go 1.27 installed; repo renamed + `git init`; scaffold created.
- **2026-08-30** Full core written (cur/sandbox/store/methods/check/tui/cli) and compile-fixed; several real bugs
  caught at first build (shadowed `normalize`, missing working dirs, args slice panic).
- **2026-08-30** Fixture suite authored; **selftest 21/21 green**.
- **2026-08-30** PLAN.md created to survive session switching.
- **2026-08-30** TUI two-column layout built + verified under tmux. Root causes found during
  debug: (a) `initViewport()` re-created the viewport after content had been set, wiping the
  spec; (b) startup selection sat on the module header, so no spec could render; (c) `setSel()`
  guarded `a.viewPort == nil`, but `viewport.Model` is a value type — that line never compiled;
  (d) the old stacked layout was taller than the terminal, so tmux scrolled and only the bottom
  rows were visible (masked the real bugs). Rebuilt as tree-only left column + viewport right
  column, both height-clamped; sel Highlighter removed in favor of inline `styleActive`.
- **2026-08-30** **M1 "C gate" authored + validated** — `subjects/M1-clib/ex01–ex05` (ft_strlen,
  ft_strlcpy, ft_memset, ft_strcmp, gate: ft_strdup). Each: strict build (`-Wall -Wextra -Werror`
  + ASan/UBSan) + quiz. Reference solutions in `answers/` all PASS; broken impl FAILs
  (determinism re-proven). Two harness traps fixed during bring-up: `"a\x00b"` silently
  becomes a 2-byte string (hex escape greed) → use `"a\x00" "b"`; `strncmp` does not stop at
  NUL, so the strlcpy harness compares only the copied region with `memcmp`.
- **2026-08-30** **M2 "Processes" authored + validated** — `subjects/M2-procs/ex01–ex05`
  (fork/wait, spawn idiom, pipes, signals, gate: pipeline). Whole-program exercises; graded
  `build`+`quiz` with multi-run specs (`ex02` success+failure, `ex05` two words). Reference
  solutions all PASS. Grader proven against real broken solutions: no-handler signal ⇒
  killed → `exit -1`; unclosed write end ⇒ 20s `timed out`. Decision recorded in §10: M2
  stays on proven build/stdout/quiz; `process`/`net` deferred to M6/M8 where timing is
  deterministic.
- **2026-08-30** **M3 "Memory" authored + validated** — `subjects/M3-memory/ex01–ex05`
  (bump arena, copy-on-write via mmap+fork, fixed slot pool, mmap allocator, gate: OOM-safe
  growable buffer). ex01/ex03/ex04/ex05 are harness-shape (header + harness `main.c`,
  student writes impl); ex02 is whole-program. Graded `build`+`quiz`; every allocator is
  static-arena/pool/registry-backed so OOM is a deterministic boundary and nothing relies on
  `malloc`. Reference solutions all PASS; broken solutions proven to FAIL: pool handing out
  the same slot (29 duplicate FAILs), `MAP_SHARED` on the private CoW page, and gate
  ignoring refusal (`r=0` bug). `forge score` → `M3-memory 5/5`, `TOTAL 16/20`.
- **2026-08-30** **M4 "Concurrency" authored + validated** — `subjects/M4-concurrency/ex01–ex05`
  (join & exit, mutex bomb, atomic counter, readers & writers, gate: bounded queue).
  Whole-program except the gate (harness shape: `queue.h` + `main.c`, student writes
  `queue.c`). **All runs under ThreadSanitizer**: `"sanitizers":["thread"]` and every run
  wrapped `setarch x86_64 -R ./test` (stock gcc TSan crashes on this kernel's ASLR;
  `-R` keeps it deterministic). Reference solutions all PASS; six broken classes proven to
  FAIL — wrong exit value (ex01, `sum 15`), omitted mutex / plain `long tally++` (ex02,
  ex03), no locks (ex04), mutex-free queue (ex05) all exit 66 via TSan `data race`; LIFO
  queue (ex05) → `fifo violated`. Bonus: `pthread_cond_wait(…, NULL)` is a build error
  under `-Werror=nonnull`. `forge score` → `M4-concurrency 5/5`, `TOTAL 21/25`.
- **2026-08-30** **M5 "Files & I/O" authored + validated** — `subjects/M5-filesio/ex01–ex05`
  (descriptors, read vs stdio, dup2 redirect, tee, gate: log parser). Whole-program except
  the gate (harness shape: `fix.h` + `main.c`, student writes `log.c`). ex02/ex05 are
  multi-run (fixtures); ex04 adds an `artifact` pinning the tee's 24-byte `out.txt`.
  Reference solutions all PASS; five broken classes proven to FAIL — reused fd number
  (`3 and 3`), buffer-size counting (`read 49` / `mismatch`), missing stdio flush
  (`captured: ` empty), stdout-only tee (`readback 0`), wrong delimiter (`parse failed`).
  **Quiz trap found**: `parseKV` splits at the first `:` — a question containing an
  internal colon (`True or false: …`) breaks the key; ex04 Q2 reworded. Broken-proof
  hygiene: manual runs must keep program stdout separate from the program's own
  `out.txt`. `forge score` → `M5-filesio 5/5`, `TOTAL 26/30`.
- **2026-08-30** **M6 "Networking" authored + validated** — `subjects/M6-networking/ex01–ex05`
  (echo, line chat, HTTP-ish, read timeout, gate: stateful protocol server). Whole-program;
  graded `build`+`net`+`quiz` — **first curriculum use of the `net` grader** (deferred from
  M2; the runner proves deterministic: retry-dial until accept, scripted send/expect with
  deadlines, byte-exact verdicts). **Engine change**: `NetSpec.connections` added — each
  session dials a fresh TCP connection against one server process, so accept-loop survival
  (ex02) and cross-connection state (ex05, `HELLO 2`) are gradeable; the earlier draft used
  two `net` methods, which `normalize()` correctly rejected as duplicate types. New
  fixtures `net/pass/ex02-sessions` + `net/fail/ex02-once` (server that exits after one
  connection), selftest now **23 fixtures, 0 mismatches**. Reference solutions all PASS;
  five broken classes proven to FAIL — hard-coded port (ex01, connect refused), single-
  connection-then-exit (ex02, `conn 2` FAIL), 404 answered with the 200 body (ex03), no
  `SO_RCVTIMEO` (ex04, deadline), no greeting-on-accept (ex05). Two behavioral traps
  surfaced in reference bring-up: `SIGPIPE` must be ignored or a closing client can kill
  the write loop; and HTTP must loop over requests per connection (persistent semantics),
  not close after the first response. `forge score` → `M6-networking 5/5`, `TOTAL 31/35`.
- **2026-08-31** **M10 "Commit & Consensus" authored + validated** —
  `subjects/M10-commit/ex01–ex05`. Fully Go, all **no-wall-clock**, byte-deterministic.
  ex01–ex04 harness-shape `build`+`stdout`+`quiz`: ex01 two-phase coordinator
  (`commit` only when EVERY participant votes commit, else `abort` for all), ex02 vote
  semantics (commit is a binding promise — abort-after-commit is forbidden), ex03 leader
  election (bully: newest logical stamp wins, ties broken by higher id, order-
  independent), ex04 quorum (strict majority `n/2+1`; a 2-group needs both, a 4-group
  needs 3 — `n/2` ties are not quorum). ex05 transaction: whole-program `gate.go`, TCP,
  decides each `tx <id>` from a scaffold `votes.txt` (presence = all participants
  prepared → `COMMIT`; absent → `ABORT`), byte-exact across two connections. Both-ways:
  references PASS; five broken classes proven FAIL — majority-commit coordinator
  (ex01), always-commit voter (ex02), id-major election (ex03, `e1: leader 5`), `n/2`
  quorum (ex04, 4-group needs 2), commit-for-every-tx gateway (ex05); refs restored
  immediately. New net fixtures `pass/ex05-tx` + `fail/ex05-txabort` → selftest
  **29 fixtures, 0 mismatches**. `forge score` → `M10-commit 5/5`, `TOTAL 52/55`; no
  orphaned `gate` processes.
- **2026-08-31** **M11 "Replicated Log & Consistency" authored + validated** —
  `subjects/M11-log/ex01–ex05`. Fully Go, all **no-wall-clock**, byte-deterministic.
  ex01–ex04 harness-shape `build`+`stdout`+`quiz`: ex01 append (`Log` tail push,
  `Append` returns index, `Len`, `All` defensive copy), ex02 committed (`NewLog`,
  `CommitIndex`, `Commit(idx)` pure setter that never lowers starting at -1,
  `Committed()` returns ONLY the durable prefix), ex03 index (pure arithmetic on
  `State{Commit,Len}`: `NextIndex` = `Len`, `Uncommitted` = `Len-Commit-1`), ex04
  snapshot (`Snapshot(upto)` drops the compactable prefix, keeps a `base` offset so
  remaining entries keep their public indexes; renumbering after compaction is the
  bug). ex05 replica: whole-program `replica.go`, TCP on `TARGETPORT`, one shared
  log for the process lifetime; `append/commit/read` where `read` returns only the
  committed prefix joined by `|` (uncommitted tail never leaks across connections).
  Both-ways: references PASS; five broken classes proven FAIL — prepend-not-append
  (ex01), `Committed` leaks uncommitted (ex02), `NextIndex=Len-1` (ex03), renumber-
  after-compaction (ex04), read-leaks-uncommitted (ex05, conn 2 `x|y|z` not `x|y`);
  refs restored immediately. New net fixtures `pass/ex05-replica` +
  `fail/ex05-replicaleak` → selftest **31 fixtures, 0 mismatches**. `forge score` →
  `M11-log 5/5`, `TOTAL 57/60`; no orphaned `replica` processes.
- **2026-08-31** **M12 "Raft — flagship capstone" authored + validated** —
  `subjects/M12-raft/ex01–ex05`. Fully Go, all **no-wall-clock**, byte-deterministic.
  ex01–ex04 harness-shape `build`+`stdout`+`quiz`: ex01 term (monotonic term,
  `SeeTerm` steps down + resets the vote only on a strictly higher term; a
  stale `see 3` at term 6 is ignored), ex02 vote (up-to-date-log rule: higher
  last-term wins, equal term → higher/equal last-index; voting on length alone
  is the bug), ex03 logmatch (matching-prefix AppendEntries check; append only
  when `log[prevIndex].Term == prevTerm`, else REJECT), ex04 election (strict
  majority `n/2+1`; `floor(n/2)` is the bug). ex05 gatenode: whole-program
  `node.go`, TCP on `TARGETPORT`, one raft node state (monotonic term + committed
  log) for the process lifetime; `see/append/read` where append is refused while
  term 0 (`not leader`) and read returns only the committed prefix, byte-exact
  across two connections. Both-ways: references PASS; five broken classes proven
  FAIL — term-regress (ex01), vote-on-length-only (ex02), append-despite-mismatch
  (ex03), `n/2` quorum (ex04), stale-see-regresses (ex05, conn 1 `see 2` got
  `term 2` not `term 3`); refs restored immediately. New net fixtures
  `pass/ex05-gatenode` + `fail/ex05-gateregress` → selftest **33 fixtures,
  0 mismatches**. `forge score` → `M12-raft 5/5`, `TOTAL 62/65`; no orphaned
  `node` processes.
- **2026-08-31** **M13 "Sharding" authored + validated** —
  `subjects/M13-shard/ex01–ex05`. Fully Go, all **no-wall-clock**, byte-deterministic.
  ex01–ex04 harness-shape `build`+`stdout`+`quiz`: ex01 slot (FNV-1a 32 slot
  assignment `% nSlots`; `banana`/`cherry` collide on 0, `apple`/`elder` on 3),
  ex02 range (key-range routing: smallest `i` with `key <= boundaries[i]`, else
  `len`; the on-boundary keys `k`→0 and `u`→2 lock inclusive `<=` semantics),
  ex03 ring (consistent hashing lookup: first node at-or-after, wrapping; exact
  hits are inclusive), ex04 rebalance (move counting by owning node hash; add
  200→1, add 700→1, remove 300→2). ex05 shardgate: whole-program `gate.go`, TCP
  on `TARGETPORT`, three isolated key→value shards with sorted boundaries
  `m`,`t`; `set/get` routes by key range and persists across connections
  (apple→0, mango→1, zebra→2, absent melon→`1:?`), byte-exact. Both-ways:
  references PASS; five broken classes proven FAIL — constant-hash (ex01),
  exclusive-`<` boundary (ex02, `k->1 u->3`), no-wrap ring (ex03, `key 600 ->
  node 2`), count-by-index (ex04, `add 200 -> 4 moved`), wrong-shard-routing
  (ex05, everything → 0); refs restored immediately. New net fixtures
  `pass/ex05-shardgate` + `fail/ex05-shardwrong` → selftest **35 fixtures,
  0 mismatches**. `forge score` → `M13-shard 5/5`, `TOTAL 67/70`; no orphaned
  `gate` processes.
- **2026-08-31** **M14 "Membership" authored + validated** —
  `subjects/M14-membership/ex01–ex05`. Fully Go, all **no-wall-clock** (timeouts
  are modeled as ticks/steps, never wall time), byte-deterministic.
  ex01–ex04 harness-shape `build`+`stdout`+`quiz`: ex01 heartbeat (`Peer`
  countdown, `Beat` resets + revives even a failed peer, `Tick` counts down to
  `failed`), ex02 suspect (`Lifecycle` alive→`suspect`→`failed` across two tick
  thresholds; `Beat` clears suspicion), ex03 gossip (`Merge(a,b)` keeps the
  higher version per member, commutative, sorted), ex04 evict (`EvictStale`
  keeps members with `now-Seq <= staleAfter`; at now=10/staleAfter=3 only
  `d`(4) evicted). ex05 membershipgate: whole-program `watch.go`, TCP on
  `TARGETPORT`, one shared name→state view (alive/suspect/dead) for the process
  lifetime; `list/status/beat/suspect/drop` where `list` shows live+suspect
  only, `drop` evicts, and `beat` revives a suspect — byte-exact across two
  connections. Both-ways: references PASS; five broken classes proven FAIL —
  fail-early + never-revive (ex01), fail-at-suspect-threshold (ex02),
  first-occurrence non-commutative merge (ex03), evict-at-bound (ex04),
  lists-dead + never-revives (ex05); refs restored immediately. New net
  fixtures `pass/ex05-membership` + `fail/ex05-membershipdrop` → selftest
  **37 fixtures, 0 mismatches**. `forge score` → `M14-membership 5/5`,
  `TOTAL 72/75`; no orphaned `watch` processes.
- **2026-08-31** **M9 "Time & Ordering" authored + validated** —
  `subjects/M9-time/ex01–ex05`. Fully Go, all **no-wall-clock** (transcript-driven,
  byte-deterministic). ex01–ex04 harness-shape `build`+`stdout`+`quiz`: ex01 Lamport
  clock (receive rule `c = max(c, s) + 1`), ex02 vector clock (own-component advance
  + elementwise max fold on receive), ex03 causality (`CausallyBefore` strict
  componentwise, `Concurrent` = incomparable — equality is not happened-before),
  ex04 total order (`(lamport, pid)` lexicographic tie-break, `sort.SliceStable`).
  ex05 sequencer: whole-program `srv.go`, TCP on `TARGETPORT`, Lamport receive rule
  with a **persistent** counter across connections, replies `R<seq> <content> <stamp>`
  byte-exact (conn 2 opens at `R12`, never resets). Both-ways: references PASS; five
  broken classes proven FAIL — receive merge ignored (ex01), vector fold ignored
  (ex02), non-strict happened-before (ex03, `e0 -> e3` wrongly true), pid-major
  ordering (ex04, `e1 e0 e2 ...`), no-stamp-merge sequencer (ex05, conn 2 `R2` not
  `R12`); refs restored immediately. New net fixtures `pass/ex05-clock` +
  `fail/ex05-reset` (sequencer that wrongly resets per connection) → selftest
  **27 fixtures, 0 mismatches**. `forge score` → `M9-time 5/5`, `TOTAL 47/50`.
- **2026-08-31** **Bug-fix + determinism-hardening pass; full re-validation** —
  three things closed: **(1) M0-ex02 shell portability**: the reference `solve.sh` uses
  `set -euo pipefail` but the grader ran it under `sh` (Debian `sh` = `dash` →
  "illegal option `pipefail`"). Per user option B, the `exercise.json` invocation is
  now `bash solve.sh` (learner keeps their own shebang); M0-tools 1/5 → 2/5.
  **(2) M5-ex04 stale-`out.txt` inheritance**: a broken tee that never writes
  `out.txt` could inherit a previous run's 24-byte file and PASS the `artifact`
  check (`make fclean` never removes `out.txt`). Closed with a new optional
  `artifact.wipe: true` (check pre-cleans the file before any method runs); both-ways
  re-proof: stdout-identical non-writing variant now FAILs `missing` where it
  previously false-PASSed. **(3) progress.db lost its M5–M8 rows** (db state was
  rebuilt from a M0–M4-only loop); every authored exercise re-checked — all 29 PASS,
  selftest 25 fixtures 0 mismatches. `forge score` → `TOTAL 42/45`.
- **2026-08-30** **M8 "Messaging + the Switch" authored + validated** —
  `subjects/M8-messaging/ex01–ex05` (framing, ordered delivery, relay, fault injection,
  gate). **Fully Go.** ex01/ex02 harness-shape (student writes `frame.go`/`queue.go`),
  graded `build`+`stdout`+`quiz`. ex03–ex05 the **switch**: the student writes
  `switch.go` (whole-program); the scaffold ships a provided counter backend
  (`backend/main.go`, built as `./svc`) the switch spawns on an ephemeral port. **Zero
  grader changes**: `StartEnv` already forwards `TARGETPORT`/`TARGETFAULTS`; faults come
  from `faults.txt` (`dup:N`/`drop:N`/`hold:N`, per-connection frame numbers) and are
  proven **data-level** via the backend's numbered, content-echoing replies — no timing
  anywhere. Both-ways: references PASS; five broken classes proven FAIL — single-`Read`
  codec, arrival-order replay, one shared backend connection (`conn 2 → R4`), hold-as-
  drop, hold released before the following frame. Findings: (a) make binary/source-dir
  collision (`backend:` target vs `backend/` dir) → binary renamed `svc`; (b) exact-
  bytes **union wants** absorb coalesced multi-replies so `netTransaction` can't
  over-read; (c) net spec hygiene — hand-written fault scenarios must recount frame
  indexes per connection (two drafts miscounted and were fixed in the .json). New
  fixtures `net/{pass, fail}/ex05-*` (server spawning a worker subprocess; fail variant
  never relays) → selftest **25 fixtures, 0 mismatches**. `forge score` →
  `M8-messaging 5/5`, `TOTAL 41/45`; no orphaned processes after multi-conn runs.
- **2026-08-30** **M7 "Resilience" authored + validated** — `subjects/M7-resilience/ex01–ex05`
  (backoff, circuit breaker, health checks, graceful shutdown, gate: mini supervisor).
  **First two-language module**: ex01 is Python — the **first `lang: python` exercise and
  first curriculum use of the `stdout` grader method** (`python3 main.py`, normalize strip,
  stderr must stay empty); ex02–ex05 are Go on a shared module Makefile (`PATH` prefixed,
  `go build -o test .`). M7 rule: nothing printed may depend on wall-clock time. ex02 is
  harness-shape (provided `main.go`, student writes `breaker.go`, fake injected clock);
  ex01/03/04/05 whole-program. Reference solutions all PASS; five broken classes proven to
  FAIL — no retry loop (ex01, missing `attempt 2/3` lines), breaker that never trips (ex02,
  no `open`/`half-open` lines), `/down` that never flips state (ex03, conn 3 `/health` still
  200), no graceful path (ex04, missing `received SIGTERM`), flaky task failing every retry
  (ex05, budget exhausted → `giving up on worker b`, exit 1). Bring-up traps: the flaky peer
  closes with unread data → RST, so `recv` raises `ConnectionResetError` (both EOF and RST
  are failed attempts); breaker probes resume at the cooldown boundary **inclusively**
  (`!now.Before(openUntil)`), not `After`. Both-ways hygiene lesson: restore references
  immediately after each broken run — the swap-and-grade loop overwrote answers' references
  and they were rebuilt from spec and re-verified PASS. `forge score` → `M7-resilience 5/5`,
  `TOTAL 36/40`; selftest 23 fixtures, 0 mismatches; no orphaned `./test` after the ex03 net
  run.

---

## 15. Next Steps (ordered)

1. `tools/selfcheck.sh`: build from clean + `forge selftest` gate (README must stay truthful).
2. Root `Makefile`: `setup` (go build), `list/show/check/score/selftest` (run `bin/forge …`), `init`.
3. Verify `forge init` + `score` end-to-end with a 2-exercise subject.
4. Write M0 module (§10) with both-ways fixtures per exercise; add its gate.
5. ~~Smoke-test the TUI~~ **DONE** — two-column tree|spec layout verified under tmux (nav,
   spec render, jump keys, width/height clamps). Remaining: async check + score view touch-up
   in the same session if needed.
6. ~~Resume M1 (C gate)~~ **DONE 2026-08-30** — module authored, validated both-ways, all PASS.
7. ~~M2 (Processes)~~ **DONE 2026-08-30** — module authored, all PASS; broken solutions
   proven to FAIL (signal-death and timeout classes). See §10 M2 detail.
8. ~~M3 (Memory)~~ **DONE 2026-08-30** — module authored, all PASS; broken solutions proven
   to FAIL (slot overlap, MAP_SHARED misuse, ignored refusal). Stayed on build/stdout/quiz
   (no `process` needed — static arenas make OOM deterministic). See §10 M3 detail.
9. ~~M4 (Concurrency)~~ **DONE 2026-08-30** — module authored, all PASS under
   ThreadSanitizer (`setarch x86_64 -R` runs); six broken classes proven to FAIL (TSan
   race exits + wrong-value/LIFO logic fails). `forge score` → M4 5/5, TOTAL 21/25.
   See §10 M4 detail.
10. ~~M5 (Files & I/O)~~ **DONE 2026-08-30** — module authored, all PASS (multi-run
     fixtures, artifact on the tee). Five broken classes proven to FAIL (fd miscount,
     buffer-size counting, missing flush, stdout-only tee, wrong delimiter). Also hardened:
     quiz questions must not contain an internal `:`. See §10 M5 detail.
11. ~~M6 (Networking)~~ **DONE 2026-08-30** — module authored, all PASS on the first
     curriculum use of the `net` grader (echo, line chat, HTTP-ish, read timeout, stateful
     gate). Grader grew `connections` sessions; five broken classes proven to FAIL. See §10
     M6 detail. `forge score` → M6 5/5, TOTAL 31/35.
12. ~~**M7 (Resilience)**~~ **DONE 2026-08-30** — module authored per the blueprint (Go,
     Python): retries/backoff, circuit breaker, health checks, graceful shutdown, watchdog
     gate. First `lang: python` exercise on the `stdout` grader; first shared-Go-Makefile
     module. 5/5 PASS, five broken classes proven to FAIL; `forge score` → M7 5/5,
     TOTAL 36/40. See §10 M7 detail.
13. ~~**M8 (Messaging + the Switch)**~~ **DONE 2026-08-30** — module authored: framing,
     ordered delivery (ex01/ex02, harness-shape Go) + the **fault-injection switch**
     (ex03 relay, ex04 `dup`/`drop`/`hold` from `faults.txt`, ex05 gate). The switch
     spawns its own counter backend, so the runner stays single-process; `StartEnv`
     already carries `TARGETFAULTS`; all fault proofs are data-level. **Zero grader
     changes** — the design questions from the original item resolved as: line-oriented
     fault file via env (no control channel), counter backend for byte-distinct
     replies, exact union wants to absorb coalesced replies. 5/5 PASS, five broken
     classes proven FAIL; new fixtures to 25; `forge score` → M8 5/5, TOTAL 41/45.
     See §10 M8 detail.
14. ~~**M9 (Time & Ordering)**~~ **DONE 2026-08-31** — module authored: logical clocks,
      Lamport, causality/vector clocks, total order, and a TCP sequencer. All no-wall-
      clock, byte-deterministic transcripts; ex01–ex04 harness-shape Go, ex05 `net`.
      5/5 PASS, five broken classes proven FAIL; new fixtures to 27; `forge score` →
      M9 5/5, TOTAL 47/50. See §10 M9 detail.
15. ~~**M10 (Commit & Consensus)**~~ **DONE 2026-08-31** — module authored: coordinator,
      vote, leader election, quorum, and a TCP transaction gateway. All no-wall-clock,
      byte-deterministic; ex01–ex04 harness-shape Go, ex05 `net`. 5/5 PASS, five broken
      classes proven FAIL; new fixtures to 29; `forge score` → M10 5/5, TOTAL 52/55.
      See §10 M10 detail.
16. ~~**M11 (Replicated Log & Consistency)**~~ **DONE 2026-08-31** — module authored: append
      log, commit index, index math, snapshot, and a TCP replica. All no-wall-clock,
      byte-deterministic; ex01–ex04 harness-shape Go, ex05 `net`. 5/5 PASS, five broken
      classes proven FAIL; new fixtures to 31; `forge score` → M11 5/5, TOTAL 57/60.
      See §10 M11 detail.
17. ~~**M12 (Raft — flagship capstone)**~~ **DONE 2026-08-31** — module authored: the core
      safety mechanics of Raft — monotonic term, up-to-date-log vote rule, matching-prefix
      AppendEntries, strict-majority quorum, and a TCP gate node. All no-wall-clock,
      byte-deterministic; ex01–ex04 harness-shape Go, ex05 `net`. 5/5 PASS, five broken
      classes proven FAIL; new fixtures to 33; `forge score` → M12 5/5, TOTAL 62/65.
      See §10 M12 detail.
18. ~~**M13 (Sharding)**~~ **DONE 2026-08-31** — module authored: FNV slot assignment,
      key-range routing, consistent-hash ring lookup, move-count rebalancing, and a TCP
      key-range shard gateway. All no-wall-clock, byte-deterministic; ex01–ex04
      harness-shape Go, ex05 `net`. 5/5 PASS, five broken classes proven FAIL; new
      fixtures to 35; `forge score` → M13 5/5, TOTAL 67/70. See §10 M13 detail.
19. ~~**M14 (Membership)**~~ **DONE 2026-08-31** — module authored: heartbeat countdown,
      alive→suspect→failed lifecycle, commutative gossip merge, staleness eviction, and
      a TCP membership gateway. All no-wall-clock (timeouts tick-modeled);
      ex01–ex04 harness-shape Go, ex05 `net`. 5/5 PASS, five broken classes proven FAIL;
      new fixtures to 37; `forge score` → M14 5/5, TOTAL 72/75. See §10 M14 detail.
20. ~~**M15 (Transactions + chaos)**~~ **DONE 2026-08-31** — module authored: 2PC/txn
      transactional store, write-ahead log with last-writer-wins replay, lost-update
      conflict detection/resolution, chaos node-availability drilling, and a TCP
      transactional KV-store gateway. All no-wall-clock, byte-deterministic; ex01–ex04
      harness-shape Go, ex05 `net`. ex01–ex04 done then; **ex05 gateway actually shipped later**
      (2026-09-02, tracker item #23, alongside the M0 gaps) → M15 5/5 and TOTAL 90/90.
      See §10 M15 detail.
21. ~~**M16 (Storage Engines)**~~ **DONE 2026-09-02** — module authored: LSM-tree
      memtable (sorted buffer), on-disk SSTable with sparse trailing index, append-only
      WAL with replay, newest-wins compaction with tombstones, and a WAL-persistent TCP
      KV gateway. All no-wall-clock, byte-deterministic; ex01–ex04 harness-shape Go,
      ex05 `net`. 5/5 PASS. Both-ways fixtures added (`net/pass|fail/ex05-persistent*`).
      See §10 M16 detail.
22. ~~**M17 (Observability)**~~ **DONE 2026-09-02** — module authored: concurrency-safe
      counter registry, span-tree tracing (explicit durations), spans with attributes +
      FindAll/MaxDur, count/sum/max/p50/p95 summary, and a TCP telemetry gateway.
      All no-wall-clock (explicit durations), byte-deterministic; ex01–ex04 harness-shape
      Go, ex05 `net`. 5/5 PASS. Both-ways fixtures added (`net/pass|fail/ex05-telemetry*`).
      Selftest then **41 fixtures, 0 mismatches**; `forge score` → M16 5/5, M17 5/5,
      TOTAL 86/90 (M0-ex03/04/05 + M15-ex05 still unshipped). See §10 M17 detail.
23. ~~**Curriculum gap-closure → TOTAL 90/90**~~ **DONE 2026-09-02** — shipped the four
      unfinished exercises, closing the last `[F]`s (mission gap: the M0 foundation was 40%
      unshipped; nothing was structurally *missing* elsewhere). Specifics:
      - **M0-ex03** `answer/check.py` (python stdlib) + `quiz.txt` → PASS.
      - **M0-ex04** `answers/Makefile` (all/fclean/re, `??=`CC, `+=`CFLAGS) + fixed `micro.c`
        (write moved before `free`) + `quiz.txt` → PASS. `subjects/` keeps the planted-bug
        `micro.c` as the pedagogical starting point; the reference lives in `answers/`.
      - **M0-ex05** `answers/mini_forge.sh` + `quiz.txt`; fixed a real **spec bug**: the
        subject demanded `set -euo pipefail` while the grader runs `sh` (dash, no `pipefail`)
        — changed to POSIX `set -eu` + guard each step, corrected the Ground-Truth pipe note,
        and fixed the `sample/expected.txt` missing trailing newline (`gate ok\n`). PASS.
      - **M15-ex05** authored `subject.md` (was missing), rewrote `txstore.go` into a coherent
        line protocol (`set/get/begin/commit/rollback/status/drop/list`, per-connection staged
        writes, shared committed store; `get`→`nil` on absent), and made `exercise.json`
        self-consistent (was contradictory: `list`+`status` expectations contradicted the
        implementation). Added both-ways net fixtures `net/pass|fail/ex05-txstore*` (ports
        17640/17641; fail = staged writes leak across connections). PASS.
      Result: **M0 5/5, M15 5/5, `forge score` → TOTAL 90/90 (0 remaining)**;
      `make selfcheck` OK, selftest **43 fixtures, 0 mismatches**; no orphaned servers.
      See §14 2026-09-02 changelog entry.
24. Nice-to-have: add an `artifact` check (`ft_strdup.c` must exist w/ `malloc` call) to
      M1-gate — presently graded by build+quiz only, which is already both-ways deterministic.
25. ~~Nice-to-have (M5 hardening)~~ **PARTIALLY CLOSED 2026-08-31** —
      `answers/` workdirs accumulate run artifacts (`test`, `*.o`, `out.txt`), and
      `make fclean` never removes `out.txt`. Added optional `artifact.wipe: true` —
      the check pre-cleans the artifact before any method runs, closing the stale-
      `out.txt` false-PASS for run-generated files (proven both-ways on M5-ex04).
      Residual: a broken Makefile whose `fclean` fails while sources are unchanged
      could still inherit a stale `test`/`switch` binary. Not observed in any module
      (curriculum Makefiles all define a working `fclean`); if it ever bites, add a
      spec-derived binary pre-clean (`./name` in runs/start/probes → `os.Remove`
      before `make all`).
26. ~~Nice-to-have (M6): orphaned-server check~~ **VERIFIED 2026-08-30** — after
      `forge check M6-ex05`, `pgrep` finds no leftover `./test` process; the net runner's
      `h.Stop(true)` kills the spawned process group cleanly. Closed without code change.
27. **Both-ways hygiene (from M7, standard since M8)**: in one scripted loop, restore the
      reference to `answers/` *immediately* after each broken-variant run — M7's swap-grade
      loop overwrote the references and they had to be rebuilt from spec. Keep a per-module
      `cp` back-to-back within the same command.

28. ~~**Resource Integrity Audit**~~ **DONE 2026-09-02** — deep-read every `subject.md` file
      and verified each of the 18 external URLs (M0–M16) resolves to a live, canonical source
      that genuinely covers the exercise's claim. No broken/TODO/empty links. Content fixes
      applied while reading: removed three leaked self-talk edit artifacts (M1-ex01 "§3.2?
      No — that is M2", M2-ex04 "§2.14 sigsuspend? No", M3-ex01 "§29.4? No"); corrected the
      TLPI chapter mis-citation §9.4→§7.1 "Allocating Memory on the Heap" (M3-ex01/ex03/ex05,
      M3 is Chapter 7, not 9); fixed M0-ex01 Gnu make typo `§9.2Phony` → `§9.2 "Phony Targets"`;
      completed the clipped TLPI §26.1 citation (M2-ex02 → "Waiting on a Child Process");
      reworded the garbled M1-ex04 goal line ("equality test is a sequence of `>` operators")
      to an accurate walk-in-lockstep description. Grading untouched (Readings only);
      `forge selftest` still 43/43. All citations are pointer-supplementary; grader logic
      (`build/stdout/quiz/net`) is fully offline.
29. **Capstone Re-engineering**: Transition `ex05` gates from simple integrators into "Mini-Capstones." Define a cohesive "Micro-App" for each module (e.g., a process-based shell for M2, a transactional KV-store for M15) to ensure each module results in a tangible portfolio project.
30. ~~**Task A — Micro-App capstones (M12 cluster + M15 durable KV)**~~ **DONE 2026-09-02** —
      added one new `ex06` capstone to M12 and M15 (no rewrite of working ex05 gates).
      - **M12-ex06 `raft cluster`** (`subjects/M12-raft/ex06-cluster/`): a **real 3-node
        raft cluster**. One binary boots peers a/b/c, each on its own TCP endpoint (client
        port +1/+2/+3), and they exchange RequestVote/AppendEntries **over real sockets**;
        the leader commits on a strict majority (2 of 3); a node steps down on a higher term.
        Byte-deterministic: the client drives `elect`/`propose`/`read`/`log`/`term` — no wall
        clock. `log b`/`log c` matching the leader's `read` proves true replication. → M12 6/6.
      - **M15-ex06 `durable txstore`** (`subjects/M15-tx/ex06-durabletxstore/`): the **Micro-App
        uniting M15 transactions with M16 WAL durability**. A durable transactional KV gateway
        where every committed `set` is atomically applied *and* appended to a write-ahead log;
        the `wal` command exposes the ordered durable record; `begin`/`commit`/`rollback` stage
        atomically, and a rollback leaves the WAL untouched. Cross-conn persistence proven. → M15 6/6.
      - Both-ways net fixtures added: `net/fail/ex06-nodewnly` (entries never replicated to
        followers → `log b`/`log c` diverge) and `net/fail/ex06-walded` (writes never reach
        the WAL → `wal` lies empty). Selftest → **45 fixtures, 0 mismatches**.
      - `forge score` → **TOTAL 92/92 (0 remaining)**. See §14 2026-09-02 entry.
31. ~~**Task B — portfolio / "so what" layer across all modules**~~ **DONE 2026-09-02** —
      closed the "portfolio" mission promise ("each module results in a tangible portfolio
      project") at the module-doc level:
      - Authored the three **missing** `module.md` files: M15 (Transactions & Chaos),
        M16 (Storage Engines), M17 (Observability) — each with its milestones and the
        deterministic rule, laying out the ex01→ex06 (M15) / ex01→ex05 arc.
      - Added a **`## So what? (interview / portfolio)`** section to **all 18** module.md
        files: a 2–3 line "why it matters on the job", 4 bullet **interview questions**,
        and a concrete **portfolio artifact** named per module (e.g. M12-ex06 raft cluster,
        M15-ex06 durable txstore, M16-ex05 persistent-kv, M17-ex05 telemetry).
      - Bumped M12 module.md milestones to include the new ex06 cluster. Docs only — grading
        untouched; `forge score` still 92/92. See §14 2026-09-02 entry.

- Update §3 snapshot + §14 log after every session.
- Keep README (user contract) and PLAN.md (internal tracker) consistent: README = what it does,
  PLAN = how it's built + where it's going.