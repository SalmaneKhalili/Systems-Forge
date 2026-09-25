# M5-ex04 · tee

The descriptor operations now combine into a tool: ex04 sends every stdin byte to stdout
and `out.txt` in one loop, then reopens the file to verify the copy. The result is a real
`tee`, not a stdout-only wrapper, and the artifact check makes the second destination
observable. Write `main.c` so `make all` produces `./test`.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test`. Given the inline input `hello world\nsecond line\n` (24 bytes), `./test` must:

1. copy stdin → stdout **and**, in the same read-write loop, → `out.txt` (created with
   `O_TRUNC`, mode 0644) using `read`/`write` only;
2. close `out.txt`, then reopen it read-only and count its bytes;
3. print `wrote <input bytes>` and `readback <file bytes>`; matching → `ALL PASS`.

Read in small chunks (e.g. 17): the loop is the exercise. `write` may return fewer bytes
than asked, so handle that short write by retrying or failing loudly. Copy before the file
is closed; you may not reopen a half-written file. Bytes must hit stdout before the verdict
lines, in input order, with no buffering games.

Compile with `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan, without redefining
`CFLAGS`/`LDFLAGS`. Print nothing to stderr on happy paths.

## Acceptance

`make all`, then `./test` with the inline input must exit 0 and print exactly:

```text
hello world
second line
wrote 24
readback 24
ALL PASS
```

- Stdout contains the 24 input bytes followed by `wrote 24`, `readback 24`, and `ALL PASS`, in that order.
- `out.txt` exists and holds exactly 24 bytes; the `artifact` check confirms it.
- The run exits 0 with empty stderr and `quiz.txt` is complete (see below).

The `build` grade is a strict compile plus `./test` with the inline input, stdout diff, and
the run's `artifact` method asserting that `out.txt` exists with exactly 24 bytes. Forgetting
to write the file produces `readback 0` → `mismatch` → FAIL, and the `artifact` check finds
no `out.txt`. Counting only what fits a fixed buffer instead of what `read` returned inflates
the totals → FAIL. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- `man 2 read`, `man 2 write`, `man 2 open` — short reads AND short writes are in the
  contract; a real `tee` of a pipe must re-loop.
- TLPI Chapter 4 "File I/O: The Universal I/O Model" — `write` on a pipe/socket may
  partially transfer; that's why `w` must be compared to `n`.
- `man 1 tee` (if the manpage is present): the real tool you are cloning.

## Quiz

1. Which system call copies a set of bytes to an open file descriptor?
2. Read() and write() return the number of bytes actually transferred.