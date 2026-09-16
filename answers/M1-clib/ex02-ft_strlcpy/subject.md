# M1-ex02 · ft_strlcpy

## Goal

Implement `ft_strlcpy` in `ft_strlcpy.c`:

```c
size_t ft_strlcpy(char *dst, const char *src, size_t size);
```

The bounded, length-aware sibling of `strcpy`: it copies *at most* `size - 1` bytes of `src`
into `dst`, always NUL-terminates — when `size > 0` — and **returns the full length of `src`**
(copying gets truncated, the return value does not). This is the property that makes callers
safe: they can detect truncation by comparing the return value against `size`. `forge` runs
`make fclean`, `make all`, `./test` as in ex01.

## Constraints

- Yours to write: `ft_strlcpy.c` and `Makefile`. Do not touch `main.c` / `ft_strlcpy.h`.
- When `size == 0`, `dst` may be garbage (even NULL) and must not be written at all —
  the harness verifies the buffer is byte-for-byte untouched.
- Return value is the length of `src` **before** truncation.
- No calls to library string functions. No redefining `CFLAGS`.

## Acceptance criteria

- [x] `"hello"` into size 10 copies `hello\0`, returns 5
- [x] `"hello"` into size 3 copies `he\0`, returns 5, and nothing past index 2 is touched
- [x] `"hello world"` (11) into size 6 copies `hello\0`, still returns **11**
- [x] `"hello"` into size 0 writes nothing at all, returns 5
- [x] size 1 yields an empty string (`dst[0] == '\0'`)
- [x] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
What does strlcpy return when the source does not fit?: <answer>
At most how many bytes of the source can a call with size n copy?: <answer>
```

## Readings

- `man 3 strcpy`, `man 3 strncpy` — read the BUGS section of `strncpy` carefully.
- The BSD `strlcpy` design rationale (man page or online):
  - https://man.openbsd.org/strlcpy.3 and the classic "strlcpy and strlcat – consistent,
    safe, string copy and concatenation" paper by Miller & de Raadt.
- K&R §5.5 again, this time for what "bounded copy" does to a `char *`.

## How you are graded

- `build`: strict compile + harness stdout diff (whitespace-insensitive). The harness checks
  return values, copied bytes, the terminating NUL, and that untouched tail bytes are intact.
- `quiz`: `quiz.txt` answers must match.
