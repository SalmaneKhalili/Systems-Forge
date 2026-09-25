# M1 · The C gate

The toolkit is ready; now the language becomes the instrument. This module is a small,
strict re-implementation bank in the 42 `libft` style: you hand-write a handful of
`string.h`/`stdlib.h` functions from scratch under a compile pipeline that treats every
warning as an error and runs AddressSanitizer + UndefinedBehaviorSanitizer on every run.
The exercises move from byte and string contracts to the gate's owned heap allocation.

## The build

- **ex01 · `ft_strlen`** — deliver `ft_strlen.c` and a `Makefile`; the harness verifies lengths including embedded NUL.
- **ex02 · `ft_strlcpy`** — deliver `ft_strlcpy.c` and a `Makefile` for a bounded copy that returns the untruncated length.
- **ex03 · `ft_memset`** — deliver `ft_memset.c` and a `Makefile`; byte fills leave untouched regions intact.
- **ex04 · `ft_strcmp`** — deliver `ft_strcmp.c` and a `Makefile` for an unsigned-byte lexicographic compare.
- **ex05 · Gate: `ft_strdup`** — deliver `ft_strdup.c` and a `Makefile`; use your own `malloc`, remain leak-free under ASan, and hand the module its memory-ownership gate.

## Rules

Every function lives in its own `.c` file next to `main.c` (the provided test harness) and
a header you may not modify. `forge` compiles them together with
`-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined`.

You write the `Makefile` yourself: targets `all` (link `test`), `fclean`, and `re`. Do
**not** redefine `CFLAGS`/`LDFLAGS`; `forge` injects the strict flags via the environment,
and a Makefile that stomps them is a fixable mistake and a quick FAIL. Do not edit
`main.c` or the `.h`; if you think you need to, you are testing yourself, not the
harness. Read the assigned pages before coding; the quiz is open-notes but meaningful
only if you actually met the material.

At the end of the module your `ft_strdup` gate must allocate, copy, terminate, and leak
nothing — all of it your own code.

## Prerequisites

Before starting M1, you should be able to do these specific things in C:

- Write a `for` loop that counts from 0 to N and accesses `array[i]`.
- Write an `if`/`else` branch that returns different values based on a condition.
- Write a function that takes `const char *s` and returns `size_t`.
- Explain what a NUL-terminated string is and why `"hello"` has length 5, not 6.
- Explain the difference between `'A'` (a char literal, integer value 65) and `"A"` (a string
  of length 1 containing `'\0'`).
- Use `sizeof` on a variable vs `sizeof` on a type.
- Explain what `unsigned char` is and why reading bytes as `unsigned char` avoids sign-extension.
- Write a simple `Makefile` with targets `all`, `clean`, `re` that compiles one `.c` file
  into a binary (M0-ex01 teaches this; you must have done it).

You do NOT need to know: pointers-to-pointers, dynamic memory allocation, structs, or the
C standard library beyond `strlen`/`strcmp`/`memcpy` at the "I can call it" level.

## So what? (interview / portfolio)

`libft`-style re-implementations are the classic entry test for C/embedded roles precisely
because they strip away the standard library and prove you understand memory ownership,
bounded writes, and undefined behavior — not that you memorised signatures. The discipline
here (own `malloc` + ASan leak check) is the exact thing interviewers probe with "walk me
through a memory bug you've fixed."

**Interview questions this module arms you for:**
- What does `strlcpy` return and why is that subtle? (the untruncated source length)
- How do you guarantee a custom `strdup` leaks nothing under AddressSanitizer?
- What is the difference between `strncpy` and a properly bounded copy?
- Why is an unsigned-byte compare the correct model for `strcmp`?

**Portfolio artifact:** M1-ex05 `ft_strdup` — a hand-written, leak-free `malloc` copy under
ASan + Werror.
