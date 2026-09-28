# M3-ex02 · copy-on-write

## Goal

Write `main.c` so `make all` produces `./test`, a program that answers one question with an
experiment: *do `fork` and memory sharing actually do what the manual pages claim?*

1. `mmap` two anonymous pages of `PROT_READ|PROT_WRITE`: one `MAP_PRIVATE`, one `MAP_SHARED`.
2. Fill both with the byte `'A'`.
3. `fork`; the child sets the first byte of **both** pages to `'B'`, then `_exit(0)` (no
   `printf`, no `return` from `main`).
4. The parent `waitpid`s the child, then reads the first byte of each page:

```
private page kept A after child write (copy-on-write)
shared page shows B (MAP_SHARED)
```

Print the lines exactly when the checks pass; on a mismatch print `FAIL: …` and exit 1.

## Constraints

- Both mappings are created **before** `fork` (that is the point — mappings are inherited).
- `mmap`/`fork` return values all checked; `MAP_FAILED` and `pid < 0` are failure paths.
- The child is allowed only `write`-class work and `_exit`. No `printf` in the child, no
  `sleep` anywhere, no `usleep`.
- `munmap` both pages at the end (leak-free discipline even though exit would reclaim them).
- No `CFLAGS` redefines. Empty stderr, deterministic output.

## Acceptance criteria

- [ ] the private page still reads `A` after the child wrote `B` (copy-on-write)
- [ ] the shared page reads `B` after the child wrote `B` (shared mapping semantics)
- [ ] the child was reaped with exit status 0
- [ ] both pages `munmap`'d; no leak report from ASan/LeakSanitizer
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Anonymous private mappings use copy-on-write across fork by default — true or false?: <answer>
Which flag makes a forked child's writes visible to the parent?: <answer>
```

## Readings

- **Reading ladder** — start with `man 2 mmap` (MAP_PRIVATE / MAP_SHARED /
  MAP_ANONYMOUS), then TLPI §49.1–49.3 below.
- The Linux Programming Interface, §49.1–49.3 "Memory Mappings" — the private/shared
  contrast and Figure 49-1 (shared vs private file mapping). Also §24.2 "Sharing of file
  offsets" sidebar about what fork does to memory.
- `man 2 mmap` — read the `MAP_PRIVATE` and `MAP_SHARED` paragraphs, and `MAP_ANONYMOUS`.
- The "copy-on-write" paragraph in TLPI §24.2, "Semantics of fork()".

## How you are graded

- `build`: strict compile + `./test` stdout diff, exit 0, empty stderr. The two checks are
  independent — flipping a flag in `mmap` flips one line; both verified.
- `quiz`: `quiz.txt` answers must match.