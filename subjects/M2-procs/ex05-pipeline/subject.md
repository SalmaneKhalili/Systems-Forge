# M2-ex05 · Gate: pipeline

The M2 gate composes the module's process verbs: two children, one pipe, two `exec`s, and
a parent that waits for all of it. The result is the skeleton of a shell running
`echo WORD | tr a-z A-Z`, with descriptor closure and reaping doing the sequencing. You
deliver `main.c` and the exact pipeline transcript.

## Shape

Write `main.c` so `make all` produces `./test`, a program that takes a word on the command
line (default `supply chain`) and produces its uppercase form by piping two exec'd programs:

1. create a pipe with `pipe(fds[2])`;
2. `fork()` child A: close read end, `dup2(fds[1], STDOUT_FILENO)`, close the original write
   end, then `execlp("echo", "echo", word, NULL)` — if it returns, print `echo failed` and
   return 127;
3. `fork()` child B: close write end, `dup2(fds[0], STDIN_FILENO)`, close the original read
   end, then `execlp("tr", "tr", "a-z", "A-Z", NULL)` — if it returns, print `tr failed` and
   return 127;
4. in the parent: close **both** pipe ends, `waitpid` child A and child B, and print:

```text
pipeline done (2 children)
```

The uppercase text lands on the parent's stdout because child B's stdout was never
redirected — it inherits the terminal/pipe the parent already has.

The uppercase mapping must come from **executing `/usr/bin/tr`** through the pipe. Writing
your own case-mapping loop is off-topic; the point of this module is not extra code, it is
making the pipeline behave. `dup2` is the only allowed redirect; close the originals right
after. Closing both ends in the parent is not optional: child A must get `SIGPIPE`-clean EOF
behavior and child B must not inherit your pipes.

Exec failure paths print `echo failed` / `tr failed` and return 127 as specified. Every
`fork`/`pipe`/`dup2`/`execlp` return is checked. No `system(3)`, no `popen(3)`, and no
`sleep`; the pipe, not timing, sequences everything. Child B can only finish after child A
closes the pipe (EOF), and the parent prints nothing until both children are reaped. No
`CFLAGS` redefines. Empty stderr, exit 0 on the happy path.

## Acceptance

`./test hello` must print exactly:

```text
HELLO
pipeline done (2 children)
```

and exit 0. `./test "42 school"` must print exactly:

```text
42 SCHOOL
pipeline done (2 children)
```

and exit 0.

- Both children are exec'd, not shelled, both are reaped, and the pipe ends are closed
  correctly.
- Removing the parent's `close(fds[1])` or a child's close produces a hang or a wrong byte
  stream; delete a close and rerun to test the descriptor contract.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` has two runs against two expected files (the two words).
Both runs passing, exit 0, and empty stderr are the whole gate; a leak, hang, or wrong
redirect appears as a diff, timeout, or sanitizer report. `quiz.txt` answers must match.

## Readings

- **Reading ladder** — the gate composes the three verbs you already trained (fork, exec,
  pipes); refresh the syscalls first, then the composed figure, then the dup2 detail.
- `man 2 pipe`, `man 2 fork`, `man 2 dup2`, `man 3 execlp`, `man 2 waitpid`.
- The Linux Programming Interface, §44.4 "Using Pipes to Connect Filters" — study the figure
  of two children joined by a pipe (TLPI Figure 44-6).
- §5.5 "Duplicating File Descriptors" (dup/dup2) if you haven't read it.
- Reread your M2-ex02 (spawn idiom) and M2-ex03 (EOF on read) — the gate is their composition.

## Quiz

1. Which syscall redirects a file descriptor to another?
2. A pipe delivers written bytes in what order?
