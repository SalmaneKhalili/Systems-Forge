# M0-ex05 · Gate: mini-forge

The M0 gate combines the Makefile, shell discipline, tree inspection, and sanitizer
workflow into one executable contract. You now write a miniature `forge`: it cleans,
builds, runs, and byte-compares a sample exercise. The result is the first end-to-end
grader you own, and the artifact you can carry into later modules.

## Shape

Write `mini_forge.sh` to grade the exercise under `sample/` in this directory:

1. `make -C sample fclean`
2. `make -C sample all` (strict warnings matter only in the real forge; here: build must succeed)
3. run `./sample/hello` and compare its output to `sample/expected.txt`
4. if every step passes, print exactly `PASS` and `exit 0`; otherwise exit nonzero
   (with `FAIL` printed)

The grader runs `sh mini_forge.sh` and requires exit 0. `mini_forge.sh` must run under
`set -eu`: `pipefail` is a bash extension, `/bin/sh` is dash, so stay POSIX and guard
each step explicitly. Use only POSIX `sh` built-ins plus `make` and `cmp` (or `diff`).

The comparison must be real: use a byte-compare (`cmp -s`) or `diff`, not eyeballs. You
are **not allowed** to modify anything under `sample/`; that is the subject matter. If you
feel the urge, the exercise is wrong, not the sample.

The ground truth for redirection and pipes is part of the script's contract. `>`
redirects standard output (stdout) to a file, overwriting its contents. `2>` redirects
standard error (stderr) to a file. `&>` or `> file 2>&1` redirects both stdout and stderr.
`<` redirects standard input (stdin) from a file. `cmd1 | cmd2` links stdout of `cmd1` to
stdin of `cmd2`. `pipefail` (bash/ksh) makes a pipeline's exit status reflect the first
failing command; POSIX `sh` (dash) has no such option, so guard each step with an
explicit `if ! … ; then` check instead.

## Acceptance

`sh mini_forge.sh` must exit 0 and print exactly:

```text
PASS
```

- The untouched `./sample/hello` must build and its output must match
  `sample/expected.txt` byte-for-byte.
- `make -C sample fclean`, `make -C sample all`, the program run, and the real comparison
  must all be guarded so failure prints `FAIL` and exits nonzero.
- Sabotaging the sample (breaking `main.c`, then rebuilding) must make the script exit
  nonzero; it cannot report success without executing the comparison.
- `quiz.txt` is complete.

Graded `scenario` + `quiz`: `scenario` runs `sh mini_forge.sh` and requires exit 0 under
the same gateway contract every future scenario-grader uses; `quiz` reads `quiz.txt`.

## Readings

- `man 1 cmp`, `man 1 sh`.
- **Ground Truth: Redirection & Pipes** (Inline Reference below).
- Re-read GNU make §2 if `-C` is not obvious.

## Quiz

1. Which variable holds the exit status of the last command?
2. What does `make -C dir` do?
