# M0-ex01 · Makefile discipline

## Goal

Write a `Makefile` in this exercise directory that builds the provided `main.c` into a
program called `hello`, using targets **`all`**, **`fclean`**, and **`re`**.

`forge` will run:

```
make fclean      # tolerated if it fails, expected to exist
make all         # must succeed under -Wall -Wextra -Werror (+ sanitizers here)
./hello          # must print:  hello from the forge
```

## Constraints

- Program name **must** be exactly `hello` — the runner execs `./hello`.
- `all` must build; `fclean` must remove build artifacts; `re` must do fclean then all.
- No `-w` / flag-suppression tricks: the flags above are injected via `CFLAGS` and you must
  not override them destructively in your Makefile (no redefining `CFLAGS`).
- A single `main.c` compiles to a single binary "out of the box". If you think you need more
  files, you are overcomplicating it.

## Acceptance criteria

- [x] `make all` produces `./hello`
- [x] `make fclean` removes `./hello`
- [x] `make re` rebuilds from scratch
- [x] `./hello` prints exactly `hello from the forge`
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
What target removes build artifacts?: <answer>
What does `make re` do?: <answer>
```

## Readings

- **Reading ladder** — two short reads: make's manual for the rules, then the C-dialect note.
- GNU make manual: §2 "An Introduction to Makefiles", §9.2 "Phony Targets" — `man make` or
  https://www.gnu.org/software/make/manual/
- The Linux Programming Interface, §3.6.1 "Feature Test Macros" (why `-std=gnu11` selects
  a C standard / feature-set variant).

## How you are graded

- `build`: strict compile + output diff against `expected.txt`.
- `quiz`: answers read from `quiz.txt` match the canonical answers.
