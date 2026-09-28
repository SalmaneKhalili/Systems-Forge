# M5-ex03 · redirect

## Goal

Write `main.c` so `make all` produces `./test` — the smallest possible **redirection**:
a child whose stdout is repointed at a file, a parent that proves it.

```
captured: redirected
ALL PASS
```

- `open("out.txt", O_CREAT|O_TRUNC|O_WRONLY, 0644)`;
- `fork()`; in the **child**: `dup2(fd, STDOUT_FILENO)` — now fd 1 *is* the file —
  `close(fd)` (the least interesting fd leak is one you just duplicated away), `printf`
  one line, **flush**, then leave via `_exit(0)`;
- in the **parent**: `close(fd)` too, `waitpid` the child, reopen `out.txt`, `read` it
  back, and print `captured: <actual>`.

If what came back is `redirected\n` the program ends `ALL PASS`; otherwise it prints the
captured bytes and exits 1.

## Constraints

- The trap: `printf` reaches the **stdio buffer**, not the file, until it flushes — and
  `_exit` **skips** `atexit` handlers *and* the stdio flush. A child that prints and
  immediately `_exit(0)` loses its own output. Flush explicitly.
- `dup2` must happen *before* anything is printed (order matters).
- Reap the child (`waitpid`) before reading the file — "redirected" must be in the file
  before the parent looks, and that only holds once the child is done.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan. Never redefine `CFLAGS`/`LDFLAGS`.

## Acceptance criteria

- [ ] `captured: redirected` + `ALL PASS`, exit 0, empty stderr
- [ ] child is waited on; `out.txt` re-read from the start; fds closed both sides
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which function repoints one file descriptor to another open file?: <answer>
What does _exit() skip that exit() runs?: <answer>
```

## Readings

- `man 2 dup2` — "after a successful return, newfd and oldfd refer to the same open file
  description".
- `man 3 fflush`, `man 2 _exit` — why the flush-and-`_exit` combo is the child discipline.
- TLPI §5.5 "duplicating file descriptors" figure 5-3; §5.4 (the *why*: every program
  already redirects — the shell did it before your `main`).
- The shell's own redirection (`cmd > out.txt`) is exactly this exercise, done before
  `execve`. You are building the part `bash` abstracts away.

## How you are graded

- `build`: strict compile + `./test`, stdout diff, exit 0, empty stderr. Forget the flush
  (or print before `dup2`, or never `dup2` at all): `out.txt` comes back empty, the
  program prints the wrong `captured:` line (or nothing captured) → FAIL.
- `quiz`: `quiz.txt` answers must match.