# M2-ex05 · Gate: pipeline

The gate of M2. Everything so far was one process handshaking with another. Now compose:
two children, one pipe, two `exec`s, and a parent who waits for all of it — the exact
skeleton of a shell running `echo WORD | tr a-z A-Z`.

## Goal

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

```
pipeline done (2 children)
```

The uppercase text lands on the parent's stdout because child B's stdout was never
redirected — it inherits the terminal/pipe the parent already has.

## Constraints

- The uppercase mapping must come from **executing `/usr/bin/tr`** through the pipe. Writing
  your own case-mapping loop is off-topic (and the grader can't fully see it — the point of
  this module is not writing extra code, it is making the pipeline behave).
- `dup2` is the only allowed redirect; close the originals right after. "close both ends in
  the parent" is not optional: child A must get `SIGPIPE`-clean EOF behavior and child B
  must not inherit your pipes.
- `exec` failure paths print `echo failed` / `tr failed` and return 127 as specified.
- Every `fork`/`pipe`/`dup2`/`execlp` return is checked. No `system(3)`, no `popen(3)`,
  no `sleep` — the pipe, not timing, sequences everything.
- Deterministic order enforced by code: child B can only finish after child A closes the pipe
  (EOF); the parent prints nothing until both are reaped.
- No `CFLAGS` redefines. Empty stderr, exit 0 on the happy path.

## Acceptance criteria

- [ ] `./test hello` prints exactly:
      ```
      HELLO
      pipeline done (2 children)
      ```
      and exits 0
- [ ] `./test "42 school"` prints exactly:
      ```
      42 SCHOOL
      pipeline done (2 children)
      ```
      and exits 0
- [ ] both children are exec'd (not shelled), both reaped, pipe ends correctly closed
- [ ] removing the parent's `close(fds[1])` or the child's close produces a hang or a wrong
      byte stream — test by deleting a close and rerunning
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which syscall redirects a file descriptor to another?: <answer>
A pipe delivers written bytes in what order?: <answer>
```

## Readings

- The Linux Programming Interface, §44.4 "Using Pipes to Connect Filters" — study the
  figure of two children joined by a pipe (TLPI Figure 44-6).
- §5.5 "Duplicating File Descriptors" (dup/dup2) if you haven't read it.
- `man 2 dup2`, `man 2 pipe`, `man 3 execlp`, `man 2 waitpid` again.

## How you are graded

- `build` with two runs against two expected files (the two words). The success of both runs
  plus exit 0 plus empty stderr is the whole gate: any leak, hang, or wrong redirect shows up
  as a diff, a timeout, or a sanitizer report.
- `quiz`: `quiz.txt` answers must match.