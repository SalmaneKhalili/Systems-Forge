# M0-ex01 · Makefile discipline

M0 starts with the build contract every later C exercise inherits: a small `Makefile` must
turn the provided source into the exact binary the runner executes. This exercise makes
that contract explicit before shell scripts, checkers, and sanitizers add their own
tooling. You deliver a `Makefile` whose `all`, `fclean`, and `re` targets work under the
strict build.

## Shape

The exercise provides `main.c`; you write `Makefile` and build the provided `main.c` into
a program called exactly `hello`. The runner executes:

```text
make fclean      # tolerated if it fails, expected to exist
make all         # must succeed under -Wall -Wextra -Werror (+ sanitizers here)
./hello          # must print:  hello from the forge
```

`all` must build; `fclean` must remove build artifacts; `re` must do fclean then all. The
program name must be exactly `hello`, because the runner execs `./hello`. A single
`main.c` compiles to a single binary out of the box; if you think you need more files,
you are overcomplicating it.

The flags above are injected via `CFLAGS`. Do not use `-w` or other flag-suppression
tricks, and do not redefine `CFLAGS` destructively in your Makefile.

## Acceptance

`make all` must produce `./hello`, `make fclean` must remove it, and `make re` must rebuild
it from scratch. The runner then executes `./hello`, which must exit 0 and print exactly:

```text
hello from the forge
```

- `make all` produces the exact `./hello` binary required by the runner.
- `make fclean` removes `./hello`; `make re` repeats `fclean` and then `all`.
- The fixed flags remain intact, and a clean exit (0) is part of the grade.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile plus an output diff against
`expected.txt`; `quiz` reads answers from `quiz.txt` and matches the canonical answers.

## Readings

- **Reading ladder** — two short reads: make's manual for the rules, then the C-dialect note.
- GNU make manual: §2 "An Introduction to Makefiles", §9.2 "Phony Targets" — `man make` or
  https://www.gnu.org/software/make/manual/
- The Linux Programming Interface, §3.6.1 "Feature Test Macros" (why `-std=gnu11` selects
  a C standard / feature-set variant).

## Quiz

1. What target removes build artifacts?
2. What does `make re` do?
