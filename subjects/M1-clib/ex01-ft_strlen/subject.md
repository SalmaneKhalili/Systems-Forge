# M1-ex01 · ft_strlen

M1 begins with the smallest useful C contract: measure a NUL-terminated string without
handing the work to the standard library. That loop establishes the byte-level discipline
used by the bounded copy, fill, comparison, and allocation exercises that follow. You
deliver `ft_strlen.c` and the `Makefile` that links the provided harness into `./test`.

## Shape

Implement your own `ft_strlen`, in `ft_strlen.c`, with this exact signature:

```c
size_t ft_strlen(const char *s);   /* declared in ft_strlen.h */
```

It must return the number of characters before the terminating **NUL** byte — nothing
else. You also write the `Makefile`. `forge` runs:

```text
make fclean   # tolerated if it fails, expected to exist
make all      # strict flags + sanitizers injected via CFLAGS
./test        # provided harness compares you against strlen(3)
```

Only `ft_strlen.c` and `Makefile` are yours to write. Do not touch `main.c` or
`ft_strlen.h`. Do not modify `NUL` handling: the string ends **at** the NUL byte. No calls
to `strlen`, `strnlen`, `memchr`, `strchr`, or `sizeof`-style shortcuts of other people's
code; you may use `sizeof` in your implementation only in the sense of your own logic,
because the point is the loop.

The code must be clean under `-std=gnu11 -Wall -Wextra -Werror` plus ASan/UBSan. `char` is
not guaranteed unsigned — read bytes as-is, because you are only counting them. The
`Makefile` must produce `./test` and provide `fclean`/`re`; never redefine `CFLAGS`.

## Acceptance

`make all` must compile `ft_strlen.c` with the provided `main.c` into `./test`, and the
harness's stdout must equal `expected.txt` exactly (whitespace-insensitive). Warnings or
sanitizer hits fail the run.

- `""` returns 0, `"hello"` returns 5, and a 48-character string returns 48.
- `"a\0b"` stops at the first NUL and returns 1.
- A 2-byte plus 3-byte UTF-8 string counts 5 bytes, not 2 characters.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is the strict compile and harness stdout diff; `quiz`
answers read from `quiz.txt` must match the canonical answers.

## Readings

Read these before you type a loop:

- `man 3 strlen` and `man 3 string` (the whole string-function family overview).
- K&R §5.5 "Character Pointers and Functions" (what `char *` really is).
- "The Linux Programming Interface", §3.6.1 "Feature Test Macros" (why `gnu11`).
- Beej's Guide to C: "Strings and String Functions".
  https://beej.us/guide/bgc/html/split/strings.html

## Quiz

1. Which character terminates a C string?
2. Which unsigned type do strlen-style functions use for sizes?
