# M7-ex04 · Graceful shutdown

## Goal

Write a Go program, whole-program (`main.go` + `Makefile`), that runs a **worker**, gets
five units of work done, then shuts down **gracefully on SIGTERM** instead of dying on the
spot.

The graded transcript is exact:

```
task 1 done
task 2 done
task 3 done
task 4 done
task 5 done
received SIGTERM
shutdown complete
```

## Constraints

- A **worker goroutine** does the work and reports each finished task through a channel;
  `main` prints `task N done` only when that completion is received (so the transcript is
  ordered by the worker, not by wall clock).
- After every unit of work is done and acknowledged, the program **raises SIGTERM to
  itself** (`syscall.Kill(syscall.Getpid(), syscall.SIGTERM)`), and must not die from it —
  this is the point. It waits for the signal, prints `received SIGTERM`, then
  `shutdown complete`, and exits **0**.
- A correct program exits 0; a program that never installs a handler is killed by SIGTERM
  and exits non-zero with a truncated transcript.
- Small `time.Sleep` for task simulation is fine; never print durations or timestamps.
- Exit code 0, nothing on stderr.

## Acceptance criteria

- [ ] the seven graded lines, in order, exit 0
- [ ] SIGTERM is handled, not fatal
- [ ] worker completion is sequenced through a channel
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which signal does this exercise handle for graceful shutdown?: <answer>
What should a graceful shutdown do to in-flight work?: <answer>
What is the exit code on the clean shutdown path?: <answer>
```

## Readings

- `man 7 signal` — the default action of SIGTERM vs how a handler changes it.
- TLPI §20 (signals), §22 (handlers) — the mechanics behind "received SIGTERM, draining".
- Go `os/signal` docs — `signal.Notify` semantics for graceful-shutdown goroutines.

## How you are graded

- `build` compiles + runs `./test`, syslout diffed against `expected.txt` (exit 0,
  empty stderr) + `quiz`. A program that aborts on SIGTERM, or prints completion lines
  without the shutdown phase → FAIL.