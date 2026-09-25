# M5 · Files & I/O

M5 works below the language runtime: every byte travels through a small table of integers,
the file descriptors that `open`, `read`, `write`, and `close` manage. ex01 makes that table
visible, ex02 compares system-call and stdio transfers, ex03 redirects a child's stdout with
`dup2`, ex04 builds a real `tee`, and ex05 parses structured text from one descriptor. The
same descriptor model carries the rest of the curriculum's I/O.

## The build

- **ex01 · descriptors** — the fd table: numbers, reuse, and the lowest available descriptor.
- **ex02 · read vs stdio** — `read()`/`write()` versus stdio, byte for byte.
- **ex03 · redirect** — repoint stdout with `dup2` across a `fork`.
- **ex04 · tee** — a real `tee`: stdin → stdout + file.
- **ex05 · Gate: log parser** — a log parser over an fd (harness shape).

## Rules

All five exercises are graded `build` + `quiz` under `-std=gnu11 -Wall -Wextra -Werror` with ASan/UBSan. ex04's copied file is additionally asserted by an `artifact` check.

Determinism is the point: verdicts are exact byte counts and exact strings, printed only after every fd is closed or every child is reaped — no PIDs, no addresses, and no races.

## Prerequisites

Before starting M5, you should be comfortable with everything from M0–M4, plus:

- Call `open(path, O_RDONLY)` and check the return value is ≥ 0 (or -1 on error).
- Explain what a file descriptor is: a small non-negative integer that the kernel uses to
  refer to an open file. `0` = stdin, `1` = stdout, `2` = stderr.
- Call `read(fd, buf, count)` and handle three cases: full read, partial read (return < count),
  and EOF (return == 0).
- Call `write(fd, buf, count)` and handle short writes (return < count).
- Call `close(fd)` and explain why you must close file descriptors (leak the descriptor table).
- Explain `errno` and `perror` — if `open` returns -1, `perror("open")` prints a human-readable
  message.
- Know what `fork` does (M2) — you will combine `fork` + `dup2` in ex03.

You do NOT need to know: `dup2`, pipes, `lseek`, or `select`/`poll`. You will learn them here.

## So what? (interview / portfolio)

The file-descriptor table is the mental model behind I/O of every kind: the same integer
interface explains redirects, sockets, and pipes. Knowing the difference between `read`/`write`
and `fread`/`fwrite`, and why `dup2` is all a redirect needs, is low-level fluency that
transfers directly to systems work.

**Interview questions this module arms you for:**
- When do file descriptors get reused, and what is "the lowest available fd"?
- What is the real difference between `read`/`write` and `fread`/`fwrite` (buffering)?
- How is `>` in a shell implemented at the API level? (`dup2` across a fork)
- How would you write a `tee` that mirrors stdin to both a file and stdout?

**Portfolio artifact:** M5-ex04 `tee` — a real `tee` (stdin → stdout + file), verified by
an exact byte-count artifact check.