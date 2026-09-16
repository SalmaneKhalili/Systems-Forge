# M5-ex04 · tee

## Goal

Write `main.c` so `make all` produces `./test` — a working **`tee`**: every byte that
comes in on stdin goes *two* places (stdout and a file), exactly once, in order.

```
hello world
second line
wrote 24
readback 24
ALL PASS
```

Given the inline input `hello world\nsecond line\n` (24 bytes), `./test` must:

1. copy stdin → stdout **and** — in the same read-write loop — → `out.txt` (created with
   `O_TRUNC`, mode 0644) using `read`/`write` only;
2. close `out.txt`, then reopen it read-only and count its bytes;
3. print `wrote <input bytes>` and `readback <file bytes>`; matching → `ALL PASS`.

## Constraints

- `read` in small chunks (e.g. 17) — the loop is the exercise; `write` may return fewer
  bytes than asked — handle that (short writes happen; retry or fail loudly).
- Copy **before** the file is closed: you may not reopen a half-written file.
- Bytes must hit stdout before the verdict lines, in input order — no buffering games.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan. Never redefine `CFLAGS`/`LDFLAGS`.
- Print nothing to stderr on happy paths.

## Acceptance criteria

- [ ] stdout is exactly the 24 input bytes followed by the three verdict lines
- [ ] `wrote 24`, `readback 24`, `ALL PASS`, exit 0, empty stderr
- [ ] `out.txt` exists and holds exactly 24 bytes (an `artifact` check confirms it)
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which system call copies a set of bytes to an open file descriptor?: <answer>
Read() and write() return the number of bytes actually transferred.: <answer>
```

## Readings

- `man 2 read`, `man 2 write`, `man 2 open` — short reads AND short writes are in the
  contract; a real `tee` of a pipe must re-loop.
- TLPI Chapter 4 "File I/O: The Universal I/O Model" — `write` on a pipe/socket may
  partially transfer; that's why `w` must be compared to `n`.
- `man 1 tee` (if the manpage is present): the real tool you are cloning.

## How you are graded

- `build`: strict compile + `./test` with the inline input; stdout diff; the run's
  `artifact` method then asserts `out.txt` exists with exactly 24 bytes. forget writing the
  file (stdout-only copy): `readback 0` → `mismatch` → FAIL, and the `artifact` check finds
  no `out.txt`. Count only what fits a fixed buffer instead of what `read` returned: totals
  inflate → FAIL.
- `quiz`: `quiz.txt` answers must match.