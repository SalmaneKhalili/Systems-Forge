# M0-ex05 · Gate: mini-forge

This is the first *gate* of the curriculum. It closes M0 and proves you can operate the
whole toolchain end to end. You are going to write a **miniature forge** — the same thing
that grades your work, in 20 lines of shell.

## Goal

Write `mini_forge.sh` that grades the exercise under `sample/` in this directory:

1. `make -C sample fclean`
2. `make -C sample all`   (strict warnings matter only in the real forge; here: build must succeed)
3. run `./sample/hello` and compare its output to `sample/expected.txt`
4. if every step passes, print exactly `PASS` and `exit 0`; otherwise exit nonzero
   (with `FAIL` printed)

The grader runs `sh mini_forge.sh` and requires exit 0.

## Constraints

- `mini_forge.sh` must run under `set -eu` (strict mode; `pipefail` is a bash extension —
  `/bin/sh` is dash, so stay POSIX and guard each step explicitly).
- Only POSIX sh built-ins plus `make` and `cmp` (or `diff`).
- The comparison must be real: a byte-compare (`cmp -s`) or `diff`, not eyeballs.
- You are **not allowed** to modify anything under `sample/` — that is the subject matter.
  If you feel the urge, the exercise is wrong, not the sample.

## Acceptance criteria

- [ ] untouched `./sample/hello` builds and its output matches `sample/expected.txt`
- [ ] `sh mini_forge.sh` prints `PASS` and exits 0
- [ ] sabotage the sample (break `main.c`, rebuild) and confirm your script exits nonzero
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which variable holds the exit status of the last command?: <answer>
What does `make -C dir` do?: <answer>
```

## Readings

- `man 1 cmp`, `man 1 sh`.
- **Ground Truth: Redirection & Pipes** (Inline Reference below).
- Re-read GNU make §2 if `-C` is not obvious.

### Ground Truth: Redirection & Pipes

1.  **Redirection:**
    - `>` redirects standard output (stdout) to a file, overwriting its contents.
    - `2>` redirects standard error (stderr) to a file.
    - `&>` or `> file 2>&1` redirects both stdout and stderr.
    - `<` redirects standard input (stdin) from a file.

2.  **Pipes:**
    - `cmd1 | cmd2` links stdout of `cmd1` to stdin of `cmd2`.
    - `pipefail` (bash/ksh) makes a pipeline's exit status reflect the first failing command;
      POSIX `sh` (dash) has no such option, so guard each step with an explicit `if ! … ; then`
      check instead.

## How you are graded

- `scenario`: `sh mini_forge.sh` must exit 0 (a real gateway check — the same contract
  every future scenario-grader uses).
- `quiz`: answers from `quiz.txt`.