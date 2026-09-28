# M0-ex03 · A five-minute checker

## Goal

Write `check.py` — Python 3, standard library only — that audits the `sample/` tree in this
directory and reports which *required* files are present.

`forge` invokes it as:

```
python3 check.py sample/
```

Rules for what "good" means:

- Required files (relative paths): `src/main.c` and `src/util.h`.
- For each required file that exists, print exactly `ok: src/main.c` etc., one per line,
  in the order given above.
- For each required file that is missing, print `missing: src/main.c`.
- Exit 0 either way: the report is the answer.
- `sample/` may contain other files (`extra.txt`); ignore them.

The grader's `sample/` currently contains both required files, so the run must print:

```
ok: src/main.c
ok: src/util.h
```

## Constraints

- `os`, `sys`, `pathlib` (stdlib) are enough. No third-party packages.
- The script must start with `#!/usr/bin/env python3` (we grep for it).
- Handle the tree path as `sys.argv[1]`; do not hard-code `sample`.
- Cross-platform-friendly path joining (use `os.path` / `pathlib`).

## Acceptance criteria

- [ ] `python3 check.py sample/` prints the two `ok:` lines in order
- [ ] `python3 check.py sample/` with `src/main.c` removed prints `missing: src/main.c`
- [ ] shebang present
- [ ] no imports outside the standard library
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
What is the standard location for the interpreter on unix?: <answer>
Which module is best for parsing command-line arguments?: <answer>
```

## Readings

- `man 1 python3` exists — skim the "Environment variables" section.
- The pathlib docs tutorial: https://docs.python.org/3/library/pathlib.html
  (the page you need more than any other this week).

## How you are graded

- `stdout`: diff of `check.py sample/` against `expected.txt`.
- `artifact`: `check.py` contains the shebang.
- `quiz`: answers from `quiz.txt`.