# M0-ex02 · Shell hygiene

## Goal

Write `solve.sh` that sums its positional arguments (all integers) and prints the sum on a
single line. It is invoked by the grader as:

```
sh solve.sh 1 2 -3 7      # must print:  7
```

## Constraints

- The script must begin with a _working_ shebang and enable `set -euo pipefail`.
  We grep the file to prove it, and the quiz checks you know what each flag does.
- No external tools beyond POSIX `sh`: use built-ins, `printf`, arithmetic.
  No `awk`, no `bc`, no `python`.
- Arguments are integers, possibly negative, at least one. Input is well-formed.

## Acceptance criteria

- [x] `sh solve.sh 1 2 -3 7` prints exactly `7`
- [x] running with `solve.sh` unbound-variable traps (`set -u`) — try calling it with
      an empty string arg (`sh solve.sh "" 5`) and explain why it does not crash
      (empty is a defined value; an _unset_ variable is the trap)
- [x] first line is `#!/usr/bin/env bash` or `#!/bin/bash`
- [x] the four flags dominate the whole script (no `set +e` shenanigans)
- [x] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which flag makes a script exit on any failing command?: <answer>
What does pipefail do?: <answer>
```

## Readings

- `man bash`, sections about the `set` builtin: `-e`, `-u`, `-o pipefail`, `-x`.
- **Ground Truth: Shell Control Flow & Arithmetic** (Inline Reference below).
- TLPI Chapter 1 "History and Standards" (§1.4 "Summary") for orientation; the
  "everything is a file" idea is the universal I/O model of §5.1–5.2 (deep-dive in M5;
  process basics come in M2).

### Ground Truth: Shell Control Flow & Arithmetic

To sum a list of arguments in a POSIX-compliant shell script:

1.  **Iterating over Arguments:**
    The special parameter `$@` represents all positional arguments passed to the script. You can loop over them using:
    ```bash
    for arg in "$@"; do
        # process $arg
    done
    ```

2.  **Arithmetic Expansion:**
    POSIX shell arithmetic uses the `$(( expression ))` syntax. Inside this structure, variables do not need the `$` prefix, and standard operators (`+`, `-`, etc.) behave as expected:
    ```bash
    sum=$(( sum + arg ))
    ```

3.  **Initializing Variables:**
    In shell scripting, make sure variables are initialized before use, especially when `set -u` (fail on unbound variables) is active:
    ```bash
    sum=0
    ```

## How you are graded

- `stdout`: your script's output must match the expected line.
- `artifact`: `solve.sh` must literally contain `set -euo pipefail`.
- `quiz`: answers from `quiz.txt`.
