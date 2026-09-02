# M7-ex05 gate · Mini supervisor

## Goal

Build the module's **gate**: a mini supervisor that runs three workers, runs a **watchdog**
that catches a worker crashing mid-shift, **restarts it after a bounded number of attempts
with real backoff**, and finally performs a **graceful shutdown**. It is the smallest
honest version of what runs your containers.

Workers (fixed, deterministic):

- worker `a`: 3 tasks, none flaky
- worker `b`: 3 tasks, **task 2 crashes on its first attempt** (the flake)
- worker `c`: 2 tasks, none flaky

The supervisor steps the workers in round-robin order (`a`, then `b`, then `c`) until every
worker has completed; a crashed task is retried. The graded transcript:

```
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

## Constraints

- **Event-driven, not time-driven**: the flake is a *programmed* failure ("task 2 fails
  until retried"), not a timer. A retried task must be labelled `retry ok` in the
  transcript.
- The watchdog **preserves completed tasks**: restarting `b` does not replay `b1`. A
  supervisor that resets a worker's progress duplicates `b1` → different transcript → FAIL.
- Backoff is real (a small `time.Sleep` between a crash and the retry) but never printed.
- If a worker fails more than **2 consecutive attempts** at the same task, print
  `giving up on worker <name>` and exit **1** (a thrashing supervisor must stop). Nothing
  in the graded transcript triggers it — that exit path is the safety valve.
- Whole-program `main.go` + `Makefile` (the module's shared Makefile). Exit 0 on the clean
  path; nothing on stderr.

## Acceptance criteria

- [ ] the 11 graded lines, in order, exit 0
- [ ] watchdog restarts exactly the crashed worker, preserving its completed progress
- [ ] retried task labelled `retry ok`, appears exactly once in the final transcript
- [ ] graceful shutdown lines after every worker completes
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
What unit does the watchdog restart after one of its tasks fails?: <answer>
How many consecutive failures trigger the supervisor to give up on a worker?: <answer>
What is the final phase after every worker finishes its tasks?: <answer>
```

## Readings

- The supervisor pattern in *Release It!* (monitoring, restart, backoff) — the watchdog
  concept this program operationalizes.
- TLPI §62 — the failures a watchdog exists to catch (a peer that simply stops making
  progress).
- Kubernetes *Restart Policies* — what a real restart policy decides (Never/Always/
  OnFailure) and why bounded restarts matter.

## How you are graded

- `build` compiles + runs `./test`, stdout diffed against `expected.txt` (exit 0, empty
  stderr) + `quiz`. Watches without restarting, restarts without backoff bounds, or resets
  progress on restart → different transcript → FAIL.