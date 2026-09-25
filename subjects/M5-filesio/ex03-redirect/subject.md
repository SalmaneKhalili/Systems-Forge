# M5-ex03 · redirect

Once a descriptor is a number in the table, `dup2` can change what that number means. ex03
uses that operation across `fork`: the child writes through fd 1 into `out.txt`, and the
parent verifies the bytes after reaping it. Write `main.c` so `make all` produces `./test`.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test`. It builds the smallest possible **redirection**. Call
`open("out.txt", O_CREAT|O_TRUNC|O_WRONLY, 0644)`, then `fork()`.

In the **child**, call `dup2(fd, STDOUT_FILENO)` so fd 1 *is* the file, close `fd` (the
least interesting fd leak is one you just duplicated away), `printf` one line, flush it,
and leave via `_exit(0)`. In the **parent**, close `fd` too, `waitpid` the child, reopen
`out.txt`, `read` it back from the start, and print `captured: <actual>`.

`dup2` must happen before anything is printed. The child must flush explicitly: `printf`
reaches the **stdio buffer**, not the file, until it flushes, and `_exit` **skips** `atexit`
handlers *and* the stdio flush. A child that prints and immediately `_exit(0)` loses its own
output. Reap the child with `waitpid` before reading the file so `redirected` is in the file
before the parent looks.

Compile with `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan, without redefining
`CFLAGS`/`LDFLAGS`.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
captured: redirected
ALL PASS
```

- The parent reads back `redirected\n` from the start of `out.txt`, then prints `ALL PASS`; otherwise it prints the captured bytes and exits 1.
- The child is waited on, fds are closed on both sides, and the successful run has empty stderr and exit 0.
- Complete `quiz.txt` (see below).

The `build` grade is a strict compile plus `./test`, stdout diff, exit 0, and empty stderr. Forgetting the flush, printing before `dup2`, or never calling `dup2` makes `out.txt` come back empty, so the program prints the wrong `captured:` line or captures nothing → FAIL. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- `man 2 dup2` — "after a successful return, newfd and oldfd refer to the same open file
  description".
- `man 3 fflush`, `man 2 _exit` — why the flush-and-`_exit` combo is the child discipline.
- TLPI §5.5 "duplicating file descriptors" figure 5-3; §5.4 (the *why*: every program
  already redirects — the shell did it before your `main`).
- The shell's own redirection (`cmd > out.txt`) is exactly this exercise, done before
  `execve`. You are building the part `bash` abstracts away.

## Quiz

1. Which function repoints one file descriptor to another open file?
2. What does _exit() skip that exit() runs?