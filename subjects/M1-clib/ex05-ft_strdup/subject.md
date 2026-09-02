# M1-ex05 · Gate: ft_strdup

The gate of M1. You have written loops, bounded copies, fills, and comparisons. Now you must
own memory: produce a heap copy of a string, with every allocated byte accounted for.

## Goal

Implement `ft_strdup` in `ft_strdup.c`:

```c
char *ft_strdup(const char *s);
```

Returns a fresh allocation on the heap containing an exact copy of `s` (including its
terminating NUL), or `NULL` if allocation fails. You may reuse your own `ft_strlen` /
`ft_memcpy`-style loops by copying them into this directory — or write it all inline. Either
way it must be **your** code.

`forge` compiles with `-fsanitize=address,undefined`: a leak, an off-by-one heap overflow, or
any UB fails you with a sanitizer report.

## Constraints

- Yours to write: `ft_strdup.c` and `Makefile`. Do not touch `main.c` / `ft_strdup.h`.
- You may **not** call `strdup`, `strcpy`, `strlen`, `memcpy`, `strcat`, etc. `malloc` (and
  only `malloc`) is allowed.
- `malloc(len)` with `strlen`-derived `len` is wrong — you must make room for the NUL.
  The harness's ASan build will find you.
- On `NULL` from `malloc`, return `NULL` and touch nothing.
- The copy must be independent: mutating it must not change the source.
- Every `ft_strdup` allocation the harness makes is freed by the harness — if anything leaks,
  it leaked inside *your* function. No `CFLAGS` redefines.

## Acceptance criteria

- [ ] `ft_strdup("")` yields an empty, NUL-terminated heap string
- [ ] `"hello"` and a 33-char string come back byte-identical
- [ ] the returned pointer is a different address than the source
- [ ] mutating the copy does not touch the source
- [ ] no leak, no overflow, no UB under ASan/UBSan (the run must exit 0)
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
What does ft_strdup return if malloc fails?: <answer>
Which sanitizer reports a leak on exit?: <answer>
```

## Readings

- `man 3 malloc` — read the RETURN VALUE sentence about `malloc(0)` and the ERRORS section;
  note glibc's `strdup` and the `_GNU_SOURCE`/`_POSIX_C_SOURCE` feature-macro wrinkle.
- "The Linux Programming Interface", §7.1 "Allocating Memory on the Heap".
- What ASan/LeakSanitizer do (and why they run as part of this exercise):
  https://clang.llvm.org/docs/AddressSanitizer.html
- K&R §5.5 + §5.6 (pointers, arrays, and pointer arithmetic — your copy loop lives here).

## How you are graded

- `build`: strict compile + harness stdout diff, plus implicit ASan/UBSan enforcement:
  any sanitizer finding makes the process exit nonzero with a report on stderr, which fails
  your run.
- `quiz`: `quiz.txt` answers must match.