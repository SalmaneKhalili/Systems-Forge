# M0-ex03 · A five-minute checker

Shell hygiene gives commands a reliable execution context; this exercise turns that
context into an inspection tool. The checker is the first small program whose output is
its contract: required files are reported in a fixed order and extra files are ignored.
You deliver `check.py`, using only Python's standard library, to audit a supplied tree.

## Shape

Write `check.py` — Python 3, standard library only — that audits the `sample/` tree in
this directory and reports which _required_ files are present. `forge` invokes it as:

```text
python3 check.py sample/
```

The required relative paths are `src/main.c` and `src/util.h`. For each required file that
exists, print exactly `ok: src/main.c` or `ok: src/util.h`, one per line, in that order.
For each missing required file, print `missing: src/main.c` (or the corresponding path).
Exit 0 either way because the report is the answer. `sample/` may contain other files such
as `extra.txt`; ignore them.

`os`, `sys`, and `pathlib` from the standard library are enough. The script must start
with `#!/usr/bin/env python3` because the grader greps for it. Handle the tree path as
`sys.argv[1]`, do not hard-code `sample`, and use cross-platform-friendly path joining with
`os.path` or `pathlib`.

## Acceptance

With the grader's current `sample/`, `python3 check.py sample/` must exit 0 and print
exactly:

```text
ok: src/main.c
ok: src/util.h
```

- The two `ok:` lines prove that both required paths were found and emitted in order.
- Removing `src/main.c` must produce `missing: src/main.c` rather than a traceback or a
  reordered report.
- The shebang is present and no import comes from outside the standard library.
- `quiz.txt` is complete.

Graded `stdout` + `artifact` + `quiz`: stdout is the diff of `check.py sample/` against
`expected.txt`; the artifact check requires the shebang, and quiz answers come from
`quiz.txt`.

## Readings

- `man 1 python3` exists — skim the "Environment variables" section.
- The pathlib docs tutorial: https://docs.python.org/3/library/pathlib.html
  (the page you need more than any other this week).

## Quiz

1. What is the standard location for the interpreter on unix?
2. Which module is best for parsing command-line arguments?
