# M2-ex03 · pipes

## Goal

Write `main.c` so `make all` produces `./test`, a program that moves exactly **19 bytes**
from a child to its parent:

1. create a pipe with `pipe(fds[2])`;
2. `fork()`;
3. in the child: close the read end, `write(fds[1], "hello over the pipe", 19)`, close the
   write end, exit;
4. in the parent: close the write end, `read` in a loop until **EOF** (read returns 0),
   close the read end, then print exactly:

```
read 19 bytes: hello over the pipe
pipe ok
```

The loop is the point: EOF happens only when *every* write end of the pipe is closed — which
means the parent must close its own copy too.

## Constraints

- Use `pipe`, `fork`, `read`, `write`, `waitpid`. No `sleep`, no `usleep`, no spin-waiting:
  the data flow itself must synchronize parent and child.
- The parent must read in a loop and check each `read` return (`< 0` is an error); accumulate
  into a buffer until `read` returns 0.
- You may not print the raw buffer without validating it: the final line must reflect a real
  check (`memcmp` against the expected 19 bytes and a byte count check). A program that
  prints `read 19 bytes` hard-coded without reading is not acceptable — but the grader can
  only see behavior, so keep the discipline.
- No `CFLAGS` redefines. Empty stderr.

## Acceptance criteria

- [ ] `./test` prints exactly the two lines above and exits 0
- [ ] the child wrote through the pipe, the parent read through it, and the bytes match
- [ ] the parent saw EOF **and** the child was reaped (`waitpid`)
- [ ] forgetting to close the parent's write-end will hang forever → your `read` loop must
      end; if you ran it and it timed out, that is the very bug this exercise exists for
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which syscall creates a pipe?: <answer>
When parent and child share a pipe, what must the parent close to see EOF on read?: <answer>
```

## Readings

- The Linux Programming Interface, §44.1–44.4 "Pipes and FIFOs" — read the pipe-plus-fork
  walkthrough carefully (Figure 44-2).
- `man 2 pipe`, `man 2 read`, `man 2 write`.
- The EOF rule is in the "Count the write ends" paragraph of TLPI §44.2 — finding it is part
  of the exercise.

## How you are graded

- `build`: strict compile; `./test` stdout diff, exit 0, empty stderr. The notorious failure
  (parent keeping its write end) times out instead of failing cleanly — that *is* your
  feedback.
- `quiz`: `quiz.txt` answers must match.