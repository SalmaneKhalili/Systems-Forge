# M1-ex03 · ft_memset

## Goal

Implement `ft_memset` in `ft_memset.c`:

```c
void *ft_memset(void *s, int c, size_t n);
```

Fills the first `n` bytes at `s` with `c` converted to `unsigned char`, and returns `s`.
Nothing here is about knowing `memset` — it is about being honest about types: the second
parameter is `int`, the byte written is `unsigned char`, and the pointer travels as `void *`.
`forge` runs `make fclean`, `make all`, `./test`.

## Constraints

- Yours to write: `ft_memset.c` and `Makefile`. Do not touch `main.c` / `ft_memset.h`.
- Must handle `n == 0` without touching anything (harness verifies via sentinel bytes).
- Must work for every `c` value, including `0x00`, `0xff`, and `0x42`, over byte-exact ranges.
- Conversion rule: `(unsigned char)c` — never cast the pointer differently.
- Returns the original pointer (not `s+n`). No `memset`/`bzero` calls. No `CFLAGS` redefines.

## Acceptance criteria

- [ ] `ft_memset(buf, 0xff, 6)` writes six `0xff` bytes and nothing else
- [ ] `ft_memset(buf, 0, 4)` clears exactly four bytes
- [ ] `ft_memset(buf + 8, 1, 3)` leaves bytes 0–7 and 11+ untouched
- [ ] `n == 0` returns the start pointer with every byte intact
- [ ] return value is always the pointer you passed in
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which header declares memset?: <answer>
What address does memset return?: <answer>
```

## Readings

- `man 3 memset` and `man 3 bzero` (note the deprecation wording).
- "The Linux Programming Interface", §2.5 "Standardized Versions of C" is the *only* TLPI
  section you need this time — plus the C99/C11 and-cast rules:
- cppreference "C data types" → the conversion rules for integral values:
  https://en.cppreference.com/w/c/language/conversion
- Why `void *` exists: K&R §5.11 "Pointers to Functions" § is next module; today just read
  K&R §5.6 "Pointer Arrays; Pointers to Pointers" intro.

## How you are graded

- `build`: strict compile + harness stdout diff. The harness paints sentinel bytes and checks
  the filled range AND that everything else is unchanged, and that the return value matches.
- `quiz`: `quiz.txt` answers must match.