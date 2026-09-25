# M7-ex05 · Gate: Mini Supervisor

The gate of M7 is the smallest honest version of what runs your containers:
a **mini supervisor** that runs three workers, a **watchdog** that catches a
worker crashing mid-shift, **restarts it after a bounded number of attempts
with real backoff**, and finally performs a **graceful shutdown**. It pulls
together the module's primitives — the backoff of ex01, the bounded-retry
discipline of ex02, and the drain of ex04 — into one event-driven transcript.

## Shape

Whole-program `main.go` + `Makefile` (the module's shared Makefile). The
workers are fixed and deterministic:

- worker `a`: 3 tasks, none flaky
- worker `b`: 3 tasks, **task 2 crashes on its first attempt** (the flake)
- worker `c`: 2 tasks, none flaky

The supervisor steps the workers in round-robin order (`a`, then `b`, then
`c`) until every worker has completed; a crashed task is retried. The graded
transcript:

```text
task a1 ok
task b1 ok
task c1 ok
task a2 ok
task b2 failed: watchdog restarting worker b
task c2 ok
task a3 ok
task b2 retry ok
task b3 ok
all workers complete
shutdown complete
```

Rules:

- **Event-driven, not time-driven**: the flake is a *programmed* failure
  ("task 2 fails until retried"), not a timer. A retried task must be labelled
  `retry ok` in the transcript.
- The watchdog **preserves completed tasks**: restarting `b` does not replay
  `b1`. A supervisor that resets a worker's progress duplicates `b1` →
  different transcript → FAIL.
- Backoff is real (a small `time.Sleep` between a crash and the retry) but
  never printed.
- If a worker fails more than **2 consecutive attempts** at the same task,
  print `giving up on worker <name>` and exit **1** (a thrashing supervisor
  must stop). Nothing in the graded transcript triggers it — that exit path is
  the safety valve.
- Exit 0 on the clean path; nothing on stderr.

## Acceptance

Graded `build` + `quiz`: compile + run `./test`, stdout diffed against
`expected.txt` (exit 0, empty stderr):

- The 11 graded lines, in order, exit 0.
- The watchdog restarts exactly the crashed worker, preserving its completed
  progress.
- The retried task is labelled `retry ok` and appears exactly once in the
  final transcript.
- Graceful shutdown lines come after every worker completes.

A watchdog without restart, a restart without backoff bounds, or a restart
that resets progress → different transcript → FAIL. `quiz.txt` is complete
(see Quiz).

## Readings

- The supervisor pattern in *Release It!* (monitoring, restart, backoff) — the watchdog
  concept this program operationalizes.
- TLPI §61.1 "Partial Reads and Writes on Stream Sockets" — the failures a watchdog
  exists to catch (a peer that simply stops making progress).
- Kubernetes *Restart Policies* — what a real restart policy decides (Never/Always/
  OnFailure) and why bounded restarts matter.

## Quiz

1. What unit does the watchdog restart after one of its tasks fails?
2. How many consecutive failures trigger the supervisor to give up on a worker?
3. What is the final phase after every worker finishes its tasks?