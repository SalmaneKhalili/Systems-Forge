# M2 · Processes

You learned to write functions; now you learn that programs are not the only thing running.
Unix is built on processes: units of execution that `fork` into new ones, `exec` new images,
`wait` for each other, and talk through pipes. The module moves from one child and one
parent to a composed `echo`→`tr` pipeline, with every lifecycle edge made deterministic.

## The build

- **ex01 · fork + wait** — deliver a `fork1`-style `main.c`: fork, child exits 7, parent reaps and verifies status.
- **ex02 · the spawn idiom** — deliver `main.c`: `fork` → `execvp` → wait, with a correct exec-failure path (exit 127).
- **ex03 · pipes** — deliver `main.c`: child writes through a pipe, parent reads till EOF and reports.
- **ex04 · signals** — deliver `main.c`: install a `SIGUSR1` handler, `raise`, and prove it ran with a flag, not `printf`.
- **ex05 · Gate: pipeline** — deliver `main.c`: two children, one pipe, `dup2` + `echo`→`tr`, with the parent reaping both.

## Rules

You write the whole program — `main.c` — and your own `Makefile` (same contract as M1:
`all`→`./test`, `fclean`, `re`; never redefine `CFLAGS`/`LDFLAGS`). `forge` compiles with
`-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined` and runs `./test` zero, one,
or more times, with arguments, diffing each invocation against an expected file
(`normalize: strip`), and requiring a clean exit and empty stderr.

All exercises must be **deterministic**: no printing raw PIDs, no races, no sleeps-for-luck.
Design so the output order is forced by the code (wait for the child, then print). Never
modify the `C` discipline learned in M1: every resource you open is closed, every return
is checked, and every child is reaped.

## Prerequisites

Before starting M2, you should be comfortable with everything from M0 and M1, plus:

- Write a `main(int argc, char **argv)` function and loop over `argv[1]` to `argv[argc - 1]`.
- Explain the difference between `exit(1)` and `return 1` from `main`.
- Call `printf` and `fprintf(stderr, ...)` and explain why `stderr` is separate from `stdout`.
- Explain what `perror` prints and how it relates to `errno`.

You do NOT need to know: signals, pipes, `dup2`, `fork`, or process groups. You will learn
those here.

## So what? (interview / portfolio)

Every OS and every container runtime is a composition of `fork`, `exec`, `wait` and pipes —
so this is the module that makes "process lifecycle" a claim you can speak to precisely.
A shell, an init system, a supervisor, and a CI runner are all this module's verbs wired
differently; the pipeline gate is a miniature of exactly that.

**Interview questions this module arms you for:**
- What actually happens when you type a command in a shell (`fork`+`execvp`+`wait`)?
- How do you prevent a zombie process, and what is the missing `wait` doing?
- What is the exec-failure convention (exit 127) and why is that a contract?
- How does a pipe transfer data, and why read till EOF?

**Portfolio artifact:** M2-ex05 `pipeline` — a fork/exec/pipe graph that runs a real
`echo`→`tr` redirection with every child reaped.
