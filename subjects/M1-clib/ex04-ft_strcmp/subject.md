# M1-ex04 · ft_strcmp

`ft_memset` made the byte representation explicit; this exercise compares those bytes
without letting the platform's signed `char` change the order. The result is the
lexicographic contract the gate's heap copy will preserve. You deliver `ft_strcmp.c` and
the `Makefile` that the provided harness links into `./test`.

## Shape

Implement `ft_strcmp` in `ft_strcmp.c`:

```c
int ft_strcmp(const char *s1, const char *s2);
```

Perform a lexicographic byte comparison: return `< 0` if `s1` comes first, `0` if equal,
and `> 0` if `s1` comes after. Walk both strings in lockstep while the bytes are equal and
non-NUL; a NUL byte is smaller than any other byte, so `"abc"` is less than `"abcd"`.
`forge` runs `make fclean`, `make all`, `./test`.

`ft_strcmp.c` and `Makefile` are yours to write; do not touch `main.c` or
`ft_strcmp.h`. The **signed-byte trap** is part of the contract: compare every byte as
`unsigned char`. `char` may be signed on your platform, so `s1[i] - s2[i]` on bytes above
0x7f is wrong. Only the sign matters; the exact magnitude is unspecified. Stop at the
first differing byte or at a NUL, and never read past a NUL except to confirm it differs.
No calls to `strcmp`/`memcmp`/`strncmp`, and no `CFLAGS` redefines.

## Acceptance

The strict build and harness stdout diff must pass. The harness checks signs and the
exact zero case, including the high-bit range where signed-versus-unsigned bugs surface.

- Equal strings, `("hello", "hello")`, return exactly 0.
- `("abc", "abd")` returns negative; `("abd", "abc")` returns positive.
- `("abc", "abcd")` returns negative because a prefix is smaller.
- `("A", "a")` returns negative, so uppercase sorts before lowercase.
- `("\xff", "\x00")` returns **positive** because bytes are compared as unsigned.
- Empty versus non-empty works in both directions.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile and harness stdout diff; `quiz.txt`
answers must match.

## Readings

- `man 3 strcmp` — read the NOTES about return-value sign and `strcasecmp`/`strcoll`
  neighbours.
- C standard note on `strcmp` byte comparison: cppreference,
  https://en.cppreference.com/w/c/string/byte/strcmp
- K&R §5.5 "Character Pointers and Functions" — the classic `strcmp` example and why the
  characters are compared "as unsigned char" in the standard.
- Beej's Guide to C, "Strings"/comparisons section (the three-way compare idiom).

## Quiz

1. strcmp must interpret bytes as which type to be correct on negative bytes?
2. When s1 sorts before s2, strcmp returns a value that is what?
