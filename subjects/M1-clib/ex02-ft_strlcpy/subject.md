# M1-ex02 · ft_strlcpy

`ft_strlen` established where a string ends; this exercise adds the destination boundary.
The copy is deliberately length-aware: a bounded write and the full source length travel
together, so callers can detect truncation without another scan. You deliver
`ft_strlcpy.c` and the `Makefile` while the provided harness checks every untouched byte.

## Shape

Implement `ft_strlcpy` in `ft_strlcpy.c`:

```c
size_t ft_strlcpy(char *dst, const char *src, size_t size);
```

This is the bounded, length-aware sibling of `strcpy`: it copies *at most* `size - 1`
bytes of `src` into `dst` and always NUL-terminates when `size > 0`. It returns the full
length of `src`; copying gets truncated, but the return value does not. Callers can detect
truncation by comparing the return value against `size`. `forge` runs `make fclean`,
`make all`, `./test` as in ex01.

`ft_strlcpy.c` and `Makefile` are yours to write; do not touch `main.c` or
`ft_strlcpy.h`. When `size == 0`, `dst` may be garbage (even NULL) and must not be
written at all — the harness verifies the buffer is byte-for-byte untouched. The return
value is the length of `src` **before** truncation. No calls to library string functions,
and no redefining `CFLAGS`.

## Acceptance

The strict build and harness must pass, with the copied bytes, return values, terminating
NUL, and untouched tail bytes all correct.

- `"hello"` into size 10 copies `hello\0` and returns 5.
- `"hello"` into size 3 copies `he\0` and returns 5; nothing past index 2 is touched.
- `"hello world"` (11) into size 6 copies `hello\0` and still returns **11**.
- `"hello"` into size 0 writes nothing at all and returns 5.
- Size 1 yields an empty string (`dst[0] == '\0'`).
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile and harness stdout diff
(whitespace-insensitive); the harness checks return values, copied bytes, the terminating
NUL, and that untouched tail bytes are intact. `quiz.txt` answers must match.

## Readings

- `man 3 strcpy`, `man 3 strncpy` — read the BUGS section of `strncpy` carefully.
- The BSD `strlcpy` design rationale (man page or online):
  - https://man.openbsd.org/strlcpy.3 and the classic "strlcpy and strlcat – consistent,
    safe, string copy and concatenation" paper by Miller & de Raadt.
- K&R §5.5 again, this time for what "bounded copy" does to a `char *`.

## Quiz

1. What does strlcpy return when the source does not fit?
2. At most how many bytes of the source can a call with size n copy?
