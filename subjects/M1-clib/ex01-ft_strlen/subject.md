# M1-ex01 · ft_strlen

## Goal

Implement your own `ft_strlen`, in `ft_strlen.c`, with this exact signature:

```c
size_t ft_strlen(const char *s);   /* declared in ft_strlen.h */
```

It must return the number of characters before the terminating **NUL** byte — nothing else.
You also write the `Makefile`. `forge` runs:

```
make fclean   # tolerated if it fails, expected to exist
make all      # strict flags + sanitizers injected via CFLAGS
./test        # provided harness compares you against strlen(3)
```

## Constraints

- Only `ft_strlen.c` and `Makefile` are yours to write. Do not touch `main.c` or
  `ft_strlen.h`. Do not modify `NUL` handling: the string ends **at** the NUL byte.
- No calls to `strlen`, `strnlen`, `memchr`, `strchr`, or `sizeof`-style shortcuts of other
  people's code (you may use `sizeof` in your implementation only in the sense of your own
  logic; the point is the loop).
- Must be clean under `-std=gnu11 -Wall -Wextra -Werror` + ASan/UBSan; `char` is not
  guaranteed unsigned — read bytes as-is, you are only counting them.
- The `Makefile` must produce `./test` and provide `fclean`/`re`; never redefine `CFLAGS`.

## Acceptance criteria

- [ ] `ft_strlen.c` compiles into `./test` with the provided `main.c`
- [ ] returns 0 for `""`, 5 for `"hello"`, 48 for a 48-char string
- [ ] stops at the first NUL: `"a\0b"` must count as 1
- [ ] counts every byte: a 2-byte + 3-byte UTF-8 string counts 5, not 2
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which character terminates a C string?: <answer>
Which unsigned type do strlen-style functions use for sizes?: <answer>
```

## Readings

Read these before you type a loop:

- `man 3 strlen` and `man 3 string` (the whole string-function family overview).
- K&R §5.5 "Character Pointers and Functions" (what `char *` really is).
- "The Linux Programming Interface", §2.5 "Standardized Versions of C" (why `gnu11`).
- Beej's Guide to C: "Strings and String Functions".
  https://beej.us/guide/bgc/html/split/strings.html

## How you are graded

- `build`: your `Makefile` compiles `main.c` + `ft_strlen.c`; the harness's stdout must equal
  `expected.txt` exactly (whitespace-insensitive). Warnings or sanitizer hits fail you.
- `quiz`: answers read from `quiz.txt` match the canonical answers.