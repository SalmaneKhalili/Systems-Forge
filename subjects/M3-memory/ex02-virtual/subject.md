# M3-ex02 · Copy-on-Write

ex01 built an allocator; this exercise answers one question with an experiment:
*do `fork` and memory sharing actually do what the manual pages claim?* You
write a whole program that `mmap`s two anonymous pages — one `MAP_PRIVATE`,
one `MAP_SHARED` — and proves that a child's writes behave differently on
each. The private/shared contrast you demonstrate is the same mechanism behind
M2's `fork`, made visible.

## Shape

Write `main.c` so `make all` produces `./test`. The program:

1. `mmap`s two anonymous pages of `PROT_READ|PROT_WRITE`: one `MAP_PRIVATE`,
   one `MAP_SHARED`.
2. Fills both with the byte `'A'`.
3. `fork`s; the child sets the first byte of **both** pages to `'B'`, then
   `_exit(0)` (no `printf`, no `return` from `main`).
4. The parent `waitpid`s the child, then reads the first byte of each page
   and prints:

```
private page kept A after child write (copy-on-write)
shared page shows B (MAP_SHARED)
```

Print those two lines exactly when the checks pass; on a mismatch print
`FAIL: …` and exit 1. Rules:

- Both mappings are created **before** `fork` — that is the point; mappings
  are inherited.
- `mmap`/`fork` return values are all checked; `MAP_FAILED` and `pid < 0` are
  failure paths.
- The child does only `write`-class work plus `_exit`. No `printf` in the
  child, no `sleep` anywhere, no `usleep`.
- `munmap` both pages at the end (leak-free discipline even though exit would
  reclaim them).
- No `CFLAGS` redefines. Empty stderr, deterministic output.

## Acceptance

Graded `build`: strict compile + `./test` stdout diff, exit 0, empty stderr.
The two checks are independent — flipping a flag in `mmap` flips one line:

- The private page still reads `A` after the child wrote `B` (copy-on-write).
- The shared page reads `B` after the child wrote `B` (shared mapping
  semantics).
- The child was reaped with exit status 0.
- Both pages are `munmap`'d; no leak report from ASan/LeakSanitizer.

`quiz.txt` is complete (see Quiz).

## Readings

- **Reading ladder** — start with `man 2 mmap` (MAP_PRIVATE / MAP_SHARED /
  MAP_ANONYMOUS), then TLPI §49.1–49.3 below.
- The Linux Programming Interface, §49.1–49.3 "Memory Mappings" — the private/shared
  contrast and Figure 49-1 (shared vs private file mapping). Also §24.2 "Sharing of file
  offsets" sidebar about what fork does to memory.
- `man 2 mmap` — read the `MAP_PRIVATE` and `MAP_SHARED` paragraphs, and `MAP_ANONYMOUS`.
- The "copy-on-write" paragraph in TLPI §24.2, "Semantics of fork()".

## Quiz

1. Anonymous private mappings use copy-on-write across fork by default — true or false?
2. Which flag makes a forked child's writes visible to the parent?