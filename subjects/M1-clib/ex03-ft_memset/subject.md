# M1-ex03 · ft_memset

The bounded copy taught you not to write past a destination; `ft_memset` makes the byte
range itself the contract. The signature exposes C's type boundary, and the sentinel
checks make every untouched byte observable. You deliver `ft_memset.c` and the
`Makefile` that the provided harness links into `./test`.

## Shape

Implement `ft_memset` in `ft_memset.c`:

```c
void *ft_memset(void *s, int c, size_t n);
```

Fill the first `n` bytes at `s` with `c` converted to `unsigned char`, and return `s`.
This is not about knowing `memset`; it is about being honest about types: the second
parameter is `int`, the byte written is `unsigned char`, and the pointer travels as
`void *`. `forge` runs `make fclean`, `make all`, `./test`.

`ft_memset.c` and `Makefile` are yours to write; do not touch `main.c` or
`ft_memset.h`. Handle `n == 0` without touching anything, as the harness verifies with
sentinel bytes. Work for every `c` value, including `0x00`, `0xff`, and `0x42`, over
byte-exact ranges. The conversion rule is `(unsigned char)c` — never cast the pointer
differently. Return the original pointer, not `s+n`. No `memset`/`bzero` calls and no
`CFLAGS` redefines.

## Acceptance

The strict build and harness stdout diff must pass. The harness paints sentinel bytes and
checks both the filled range and that everything else is unchanged, including the return
value.

- `ft_memset(buf, 0xff, 6)` writes six `0xff` bytes and nothing else.
- `ft_memset(buf, 0, 4)` clears exactly four bytes.
- `ft_memset(buf + 8, 1, 3)` leaves bytes 0–7 and 11+ untouched.
- `n == 0` returns the start pointer with every byte intact.
- The return value is always the pointer passed in.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile and harness stdout diff; `quiz.txt`
answers must match.

## Readings

- `man 3 memset` and `man 3 bzero` (note the deprecation wording).
- "The Linux Programming Interface", §3.6.1 "Feature Test Macros" is the *only* TLPI
  section you need this time — plus the C99/C11 and-cast rules:
- cppreference "C data types" → the conversion rules for integral values:
  https://en.cppreference.com/w/c/language/conversion
- Why `void *` exists: K&R §5.11 "Pointers to Functions" § is next module; today just read
  K&R §5.6 "Pointer Arrays; Pointers to Pointers" intro.

## Quiz

1. Which header declares memset?
2. What address does memset return?
