# M7-ex04 · Graceful Shutdown

A resilient program also exits well. This exercise runs a **worker**, gets five
units of work done, then shuts down **gracefully on SIGTERM** instead of
dying on the spot — sequencing completions through a channel so the transcript
is ordered by the worker, not by the wall clock. It is the drain phase your
supervisor (ex05) and every long-running service you write from here on will
perform.

## Shape

Write a Go program, whole-program (`main.go` + `Makefile`). The graded
transcript is exact:

```text
task 1 done
task 2 done
task 3 done
task 4 done
task 5 done
received SIGTERM
shutdown complete
```

Rules:

- A **worker goroutine** does the work and reports each finished task through
  a channel; `main` prints `task N done` only when that completion is received
  (so the transcript is ordered by the worker, not by wall clock).
- After every unit of work is done and acknowledged, the program **raises
  SIGTERM to itself** (`syscall.Kill(syscall.Getpid(), syscall.SIGTERM)`), and
  must not die from it — this is the point. It waits for the signal, prints
  `received SIGTERM`, then `shutdown complete`, and exits **0**.
- A correct program exits 0; a program that never installs a handler is killed
  by SIGTERM and exits non-zero with a truncated transcript.
- Small `time.Sleep` for task simulation is fine; never print durations or
  timestamps.
- Exit code 0, nothing on stderr.

## Acceptance

Graded `build` + `quiz`: compile + run `./test`, stdout diffed against
`expected.txt`, exit 0, empty stderr:

- The seven graded lines, in order, exit 0.
- SIGTERM is handled, not fatal — a program that aborts on SIGTERM, or prints
  completion lines without the shutdown phase, fails.
- Worker completion is sequenced through a channel.

`quiz.txt` is complete (see Quiz).

## Readings

- `man 7 signal` — the default action of SIGTERM vs how a handler changes it.
- TLPI §20 (signals), §21 (signal handlers) — the mechanics behind "received SIGTERM,
  draining".
- Go `os/signal` docs — `signal.Notify` semantics for graceful-shutdown goroutines.

## Quiz

1. Which signal does this exercise handle for graceful shutdown?
2. What should a graceful shutdown do to in-flight work?
3. What is the exit code on the clean shutdown path?