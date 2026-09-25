# M5-ex02 · read vs stdio

With the descriptor table established, ex02 uses the same fd through two APIs. A system-call
reader and an stdio reader consume each fixture in different chunk sizes, then compare
independent totals and checksums. Write `main.c` so `make all` produces `./test FILE` and
the two provided runs agree byte for byte.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test FILE`. Run it once with `input.txt` and once with `second.txt`; both fixtures are
provided.

- With the **system call**, call `open(argv[1])`, read in chunks of **7** bytes, and keep
  the running total and byte-checksum.
- With **stdio**, call `fopen(argv[1], "r")`, use `fread()` in chunks of **11** bytes, and
  keep the same total and checksum.
- If both totals and both checksums match, print `match 1`.

The checksum is the sum of all bytes (`unsigned char`), cast wide so it cannot overflow. Both
paths must handle short reads (`read` returning fewer bytes than requested is normal), EOF
(`read`/`fread` return 0), and errors (`return 1` on a negative read). The two chunk sizes
are deliberately different: buffer reuse is fine, but totals and checksums must come from
actual counts, never the buffer size.

Compile with `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan, without redefining
`CFLAGS`/`LDFLAGS`. Print nothing to stderr on happy paths.

## Acceptance

`make all`, then `./test input.txt` and `./test second.txt` must each exit 0 and print exactly:

```text
input.txt read 44 checksum 4045
input.txt fread 44 checksum 4045
match 1
ALL PASS
```

```text
second.txt read 11 checksum 715
second.txt fread 11 checksum 715
match 1
ALL PASS
```

- For each fixture, line 1 and line 2 agree on the `read` and `fread` totals and checksums, and the program prints `match 1` and `ALL PASS`.
- Short reads and EOF are handled; a buffer size is never mistaken for bytes transferred.
- Both runs exit 0 with empty stderr and `quiz.txt` is complete (see below).

The `build` grade is a strict compile plus `./test input.txt` and `./test second.txt`, each stdout diffed against its own expected file, with exit 0 and empty stderr. Counting `read` bytes as if the buffer were always full inflates totals and makes the checksums disagree → FAIL. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- `man 2 read`, `man 2 open`, `man 2 close` — the syscall contract: return value is "bytes
  transferred", 0 means EOF, -1 an error; partial reads are possible.
- `man 3 fread`, `man 3 feof` — stdio buffers internally; `fread` fills a caller buffer on
  top of that.
- TLPI Chapter 4 "File I/O: The Universal I/O Model" (§4.2 "Universality of I/O") — the
  same `read`/`write` pair that moves bytes for files, devices, pipes, and sockets.
- TLPI §13.2 "Buffering in the stdio Library" if you want the syscall-vs-stdio buffering
  split explained.

## Quiz

1. Which system call returns bytes to a buffer you provide?
2. Which stdio function reads fixed-size blocks into a buffer?