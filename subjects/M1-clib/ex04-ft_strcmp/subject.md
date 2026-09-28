# M1-ex04 · ft_strcmp

## Goal

Implement `ft_strcmp` in `ft_strcmp.c`:

```c
int ft_strcmp(const char *s1, const char *s2);
```

Lexicographic byte comparison: returns `< 0` if `s1` comes first, `0` if equal, `> 0` if `s1`
comes after. Walk both strings in lockstep while the bytes are equal and non-NUL; a NUL byte
is smaller than any other byte, so `"abc"` is less than `"abcd"`. `forge` runs `make fclean`, `make all`, `./test`.

## Constraints

- Yours to write: `ft_strcmp.c` and `Makefile`. Do not touch `main.c` / `ft_strcmp.h`.
- **Signed-byte trap:** compare every byte as `unsigned char`. `char` may be signed on your
  platform, so `s1[i] - s2[i]` on bytes above 0x7f is wrong. You only care about the *sign*;
  the exact magnitude is unspecified.
- Compare stops at the first differing byte or at a NUL — never read past a NUL except to
  confirm it differs.
- No calls to `strcmp`/`memcmp`/`strncmp`. No `CFLAGS` redefines.

## Acceptance criteria

- [ ] equal strings return exactly 0: `("hello", "hello")`
- [ ] `("abc", "abd")` returns negative; `("abd", "abc")` returns positive
- [ ] `("abc", "abcd")` returns negative (prefix is smaller)
- [ ] `("A", "a")` returns negative (uppercase sorts before lowercase)
- [ ] `("\xff", "\x00")` returns **positive** — bytes compared as unsigned
- [ ] empty vs non-empty works in both directions
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
strcmp must interpret bytes as which type to be correct on negative bytes?: <answer>
When s1 sorts before s2, strcmp returns a value that is what?: <answer>
```

## Readings

- `man 3 strcmp` — read the NOTES about return-value sign and `strcasecmp`/`strcoll`
  neighbours.
- C standard note on `strcmp` byte comparison: cppreference,
  https://en.cppreference.com/w/c/string/byte/strcmp
- K&R §5.5 "Character Pointers and Functions" — the classic `strcmp` example and why the
  characters are compared "as unsigned char" in the standard.
- Beej's Guide to C, "Strings"/comparisons section (the three-way compare idiom).

## How you are graded

- `build`: strict compile + harness stdout diff. The harness checks signs and the exact zero
  case, including the high-bit range where signed-vs-unsigned bugs surface.
- `quiz`: `quiz.txt` answers must match.