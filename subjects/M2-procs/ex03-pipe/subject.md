# M2-ex03 · pipes

With `fork` and `wait` in place, a pipe gives the child and parent a data path as well as
a lifecycle path. EOF is the handoff that proves every write end is gone, so the parent
must close its copy and read until the kernel reports zero bytes. You deliver `main.c`
and the exact two-line report.

## Shape

Write `main.c` so `make all` produces `./test`, a program that moves exactly **19 bytes**
from a child to its parent:

1. create a pipe with `pipe(fds[2])`;
2. `fork()`;
3. in the child: close the read end, `write(fds[1], "hello over the pipe", 19)`, close the
   write end, exit;
4. in the parent: close the write end, `read` in a loop until **EOF** (read returns 0),
   close the read end, then print exactly:

```text
read 19 bytes: hello over the pipe
pipe ok
```

The loop is the point: EOF happens only when *every* write end of the pipe is closed,
which means the parent must close its own copy too. Use `pipe`, `fork`, `read`, `write`,
and `waitpid`; no `sleep`, no `usleep`, and no spin-waiting. The data flow itself must
synchronize parent and child.

The parent must read in a loop and check each `read` return (`< 0` is an error),
accumulating into a buffer until `read` returns 0. You may not print the raw buffer without
validating it: the final line must reflect a real check (`memcmp` against the expected 19
bytes and a byte-count check). A program that prints `read 19 bytes` hard-coded without
reading is not acceptable, although the grader can only see behavior, so keep the
discipline. No `CFLAGS` redefines. Empty stderr.

## Acceptance

`./test` must print exactly:

```text
read 19 bytes: hello over the pipe
pipe ok
```

and exit 0.

- The child writes through the pipe, the parent reads through it, and the 19 bytes match.
- The parent sees EOF and reaps the child with `waitpid`.
- Forgetting to close the parent's write end hangs forever; the `read` loop must end, and
  a timeout identifies exactly that bug.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile, `./test` stdout diff, exit 0, and
empty stderr. The notorious parent-keeps-write-end failure times out instead of failing
cleanly, and that is the exercise's feedback. `quiz.txt` answers must match.

## Readings

- **Reading ladder** — the EOF-on-read rule is the whole exercise; get the syscall contract
  first, then the pipe-plus-fork walkthrough in the book.
- `man 2 pipe`, `man 2 fork`, `man 2 read`, `man 2 write`, `man 2 waitpid`.
- The Linux Programming Interface, §44.1–44.4 "Pipes and FIFOs" — read the pipe-plus-fork
  walkthrough carefully (Figure 44-2).
- The EOF rule is in the "Count the write ends" paragraph of TLPI §44.2 — finding it is part
  of the exercise.

## Quiz

1. Which syscall creates a pipe?
2. When parent and child share a pipe, what must the parent close to see EOF on read?
