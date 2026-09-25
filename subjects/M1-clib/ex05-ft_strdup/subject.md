# M1-ex05 · Gate: ft_strdup

The earlier functions cover loops, bounded writes, fills, and comparisons. The gate adds
ownership: a heap copy must be correctly sized, independent from its source, and clean
under the sanitizers. You deliver `ft_strdup.c` and the `Makefile` that make that
allocation contract real.

## Shape

Implement `ft_strdup` in `ft_strdup.c`:

```c
char *ft_strdup(const char *s);
```

Return a fresh heap allocation containing an exact copy of `s`, including its terminating
NUL, or `NULL` if allocation fails. You may reuse your own `ft_strlen` / `ft_memcpy`-style
loops by copying them into this directory, or write everything inline. Either way it
must be **your** code.

`forge` compiles with `-fsanitize=address,undefined`: a leak, an off-by-one heap overflow,
or any UB fails you with a sanitizer report. `ft_strdup.c` and `Makefile` are yours to
write; do not touch `main.c` or `ft_strdup.h`. You may not call `strdup`, `strcpy`,
`strlen`, `memcpy`, `strcat`, or similar functions. `malloc` (and only `malloc`) is
allowed.

`malloc(len)` with a `strlen`-derived `len` is wrong: make room for the NUL. If `malloc`
returns `NULL`, return `NULL` and touch nothing. The copy must be independent, so
mutating it must not change the source. Every `ft_strdup` allocation the harness makes
is freed by the harness; anything that leaks leaked inside *your* function. No `CFLAGS`
redefines.

## Acceptance

The strict build and harness stdout diff must pass, with ASan/UBSan providing the
allocation check. The run must exit 0 with no leak, overflow, or UB report.

- `ft_strdup("")` yields an empty, NUL-terminated heap string.
- `"hello"` and a 33-character string come back byte-identical.
- The returned pointer has a different address from the source.
- Mutating the copy does not touch the source.
- The harness frees each allocation, and no sanitizer finding reaches stderr.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile and harness stdout diff with implicit
ASan/UBSan enforcement; any sanitizer finding makes the process exit nonzero with a report
on stderr and fails the run. `quiz.txt` answers must match.

## Readings

- `man 3 malloc` — read the RETURN VALUE sentence about `malloc(0)` and the ERRORS section;
  note glibc's `strdup` and the `_GNU_SOURCE`/`_POSIX_C_SOURCE` feature-macro wrinkle.
- "The Linux Programming Interface", §7.1 "Allocating Memory on the Heap".
- What ASan/LeakSanitizer do (and why they run as part of this exercise):
  https://clang.llvm.org/docs/AddressSanitizer.html
- K&R §5.5 + §5.6 (pointers, arrays, and pointer arithmetic — your copy loop lives here).

## Quiz

1. What does ft_strdup return if malloc fails?
2. Which sanitizer reports a leak on exit?
