# M0-ex02 · Shell hygiene

After the Makefile fixes the build entry points, shell hygiene makes the commands around
that build safe to compose. This exercise hands you a small integer-sum script and makes
its failure policy part of the artifact. You deliver `solve.sh` with strict shell settings,
bounded dependencies, and one deterministic line of output.

## Shape

Write `solve.sh` that sums its positional arguments (all integers) and prints the sum on a
single line. The grader invokes it as:

```text
sh solve.sh 1 2 -3 7      # must print:  7
```

The script must begin with a _working_ shebang and enable `set -euo pipefail`. We grep
the file to prove it, and the quiz checks you know what each flag does. Use only POSIX
`sh` built-ins, `printf`, and arithmetic: no `awk`, no `bc`, no `python`. Arguments are
integers, possibly negative, with at least one argument; input is well-formed.

The ground-truth implementation pattern is small and explicit. `$@` represents all
positional arguments, so iterate over them with:

```bash
for arg in "$@"; do
    # process $arg
done
```

POSIX shell arithmetic uses `$(( expression ))`; inside it variables do not need the `$`
prefix and standard operators such as `+` and `-` behave as expected. Accumulate with:

```bash
sum=$(( sum + arg ))
```

Initialize the variable before use, especially under `set -u`, which fails on unbound
variables:

```bash
sum=0
```

## Acceptance

`sh solve.sh 1 2 -3 7` must print exactly `7`. The grader also calls the script with an
empty string argument (`sh solve.sh "" 5`) under the unbound-variable trap and expects it
not to crash: empty is a defined value, while an _unset_ variable is the trap.

- The first line is `#!/usr/bin/env bash` or `#!/bin/bash`.
- The four flags dominate the whole script, with no `set +e` shenanigans.
- `quiz.txt` is complete.

Graded `stdout` + `artifact` + `quiz`: stdout must match the expected line, `solve.sh`
must literally contain `set -euo pipefail`, and the answers from `quiz.txt` must match.

## Readings

- `man bash`, sections about the `set` builtin: `-e`, `-u`, `-o pipefail`, `-x`.
- **Ground Truth: Shell Control Flow & Arithmetic** (Inline Reference below).
- TLPI Chapter 1 "History and Standards" (§1.4 "Summary") for orientation; the
  "everything is a file" idea is the universal I/O model of §5.1–5.2 (deep-dive in M5;
  process basics come in M2).

## Quiz

1. Which flag makes a script exit on any failing command?
2. What does pipefail do?
