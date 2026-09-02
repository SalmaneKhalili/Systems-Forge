# M1 · The C gate

The piscine's first instrument was the toolkit; now the instrument is the language. This
module is a small, strict re-implementation bank in the 42 `libft` style: you hand-write a
handful of `string.h`/`stdlib.h` functions from scratch, under a compile pipeline that treats
every warning as an error and runs AddressSanitizer + UndefinedBehaviorSanitizer on every run.

**Rules of the module**

- Every function lives in its own `.c` file next to `main.c` (the provided test harness) and a
  header you may not modify. `forge` compiles them together with
  `-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined`.
- You write the `Makefile` yourself: targets `all` (link `test`), `fclean`, `re`. Do **not**
  redefine `CFLAGS`/`LDFLAGS` — `forge` injects the strict flags via the environment; a
  Makefile that stomps them is a fixable mistake and a quick FAIL.
- Do not edit `main.c` or the `.h`. If you think you need to, you are testing yourself, not
  the harness.
- Read the assigned pages before coding; the quiz is open-notes but the answers are only
  meaningful if you actually met the material.
- At the end of the module your `ft_strdup` (gate) must allocate, copy, terminate, and leak
  nothing — all of it your own code.

| ex | function | artifact you deliver |
|----|----------|----------------------|
| ex01 | `ft_strlen` | `ft_strlen.c` + `Makefile`; harness verifies lengths incl. embedded NUL |
| ex02 | `ft_strlcpy` | `ft_strlcpy.c` + `Makefile`; bounded copy returning the untruncated length |
| ex03 | `ft_memset` | `ft_memset.c` + `Makefile`; byte fills, untouched regions stay intact |
| ex04 | `ft_strcmp` | `ft_strcmp.c` + `Makefile`; unsigned-byte lexicographic compare |
| ex05 | **Gate: `ft_strdup`** | `ft_strdup.c` + `Makefile`; own `malloc`, leak-free under ASan |