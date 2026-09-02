# M5-ex02 · read vs stdio

## Goal

Write `main.c` so `make all` produces `./test FILE` **twice**, once for `input.txt`, once
for `second.txt` (both provided). Read and count the same file two ways, and prove they
agree:

- with the **system call**: `open(argv[1])`, `read()` in chunks of **7** bytes, keep the
  running total and a byte-checksum;
- with **stdio**: `fopen(argv[1], "r")`, `fread()` in chunks of **11** bytes, same total
  and checksum;
- if both totals and both checksums match, `match 1`.

```
input.txt read 44 checksum 4045
input.txt fread 44 checksum 4045
match 1
ALL PASS
```

```
second.txt read 11 checksum 715
second.txt fread 11 checksum 715
match 1
ALL PASS
```

The checksum is the sum of all bytes (`unsigned char`), cast wide so it cannot overflow.

## Constraints

- Both paths must handle **short reads** (`read` returning fewer bytes than requested is
  normal), EOF (`read`/`fread` return 0), and errors (`return 1` on a negative read).
- The two chunk sizes are deliberately different — reuse of a buffer is fine, but totals
  and checksums must come out of *actual* counts, never the buffer size.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan. Never redefine `CFLAGS`/`LDFLAGS`.
- Print nothing to stderr on happy paths.

## Acceptance criteria

- [ ] both runs: line 1 and line 2 agree (`read` `fread` totals + equal checksums)
- [ ] `match 1` + `ALL PASS`, exit 0, empty stderr on both fixtures
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which system call returns bytes to a buffer you provide?: <answer>
Which stdio function reads fixed-size blocks into a buffer?: <answer>
```

## Readings

- `man 2 read`, `man 2 open`, `man 2 close` — the syscall contract: return value is "bytes
  transferred", 0 means EOF, -1 an error; partial reads are possible.
- `man 3 fread`, `man 3 feof` — stdio buffers internally; `fread` fills a caller buffer on
  top of that.
- TLPI §5.1–5.2: "the universal I/O model" — the same `read`/`write` pair that moves
  bytes for files, devices, pipes, and sockets.
- TLPI §12.2 if you want the syscall-vs-stdio buffering split explained.

## How you are graded

- `build`: strict compile + `./test input.txt` and `./test second.txt`, each stdout diffed
  against its own expected file, exit 0, empty stderr. Count `read` bytes as if the buffer
  were always full (the classic bug): totals inflate and the checksums disagree → FAIL.
- `quiz`: `quiz.txt` answers must match.