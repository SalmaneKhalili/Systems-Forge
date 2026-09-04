# systems-forge

A complete, self-graded curriculum for systems and distributed-systems engineering,
built in the 42 piscine tradition: every exercise is a real, working artifact —
an actual server, allocator, protocol, or consensus implementation — never a toy.

Everything runs locally. No daemons, no servers, no accounts.

## Quick start

```sh
make setup        # build the forge binary (requires go 1.27+, gcc/clang, make, python3)
./bin/forge init  # create your personal workspace under answers/
./bin/forge       # launch the TUI
./bin/forge list  # or drive everything from the command line
```

`make` also works as shorthand: `make check M11-ex03`, `make show M0-ex01`, `make score`,
`make selftest`.

## Where things live

```
subjects/     the curriculum — read-only canon: spec (subject.md), scaffold, harness, and
              expected outputs. Never a solved deliverable.
answers/      your personal workspace (gitignored) — WRITE your solutions here.
solutions/    reference solutions (gitignored) — compare AFTER you solve, never before.
tools/anki/   spaced-repetition review decks, generated from every exercise's Q&A review
              metadata (see "Review with Anki" below).
tools/        selfcheck, method fixtures, assistant graders.
forge/        the platform source (TUI + CLI + grader engine, Go).
caps/         (reserved) capstone specifications — currently unused, empty.
switch/       (reserved) the partitionable transport — currently unused, empty.
progress.db   local progress store (gitignored, SQLite) written by forge.
```

## The solve loop (how you actually use it)

1. **Init once**: `./bin/forge init` mirrors `subjects/` into `answers/`. If you run it
   again later it copies over new scaffold files but **never overwrites a file you've
   already written**.
2. **Pick an exercise**: `./bin/forge list` shows every module and exercise id (e.g.
   `M11-ex03`), plus what you've already passed.
3. **Read the spec**: `./bin/forge show M11-ex03` (or open `answers/M11-log/ex03-index/`
   in an editor). The spec names the exact deliverable you must write — e.g. `log.go`,
   `server.go`, `cluster.go`, `check.py`, `solve.sh`.
4. **Write it** in the matching file under `answers/…` (the scaffold already names it).
5. **Grade it**: `./bin/forge check M11-ex03` runs the autograder and records the result.
   A bare workspace auto-FAILs every exercise — you must write the deliverable to pass.
6. **Compare (after you pass)**: open the matching reference in `solutions/` to see an
   expert implementation. Never read solutions before you've passed the exercise.

> Tip: `./bin/forge selftest` validates the grader itself (`make selfcheck` too). If the
> grader ever seems wrong, run this first — it's the platform's own check on itself.

## Grading model

Every exercise declares one or more grading methods (see `subjects/*/exercise.json`):

| method     | what it verifies                                                           |
| ---------- | -------------------------------------------------------------------------- |
| `build`    | `make` under `-Wall -Wextra -Werror` (+ sanitizers) and runs test binaries |
| `stdout`   | normalized diff of a program's output against expected                     |
| `artifact` | properties of produced files (type, symbols, sizes, contents)              |
| `process`  | spawns a program, drives it with signals/probes, asserts state             |
| `net`      | drives a scripted protocol over a socket and asserts responses             |
| `report`   | parses structured output and compares to a reference run                   |
| `scenario` | runs a fault-injection driver against a multi-node system                  |
| `fault`    | asserts a _property_ under injected failures (partition, crash)            |
| `lincheck` | checks operation histories for linearizability                             |

C exercises are compiled with `-std=gnu11 -Wall -Wextra -Werror` and, where the spec
sets sanitizers, `-fsanitize=address,undefined` (and `thread` for the concurrency
module). This discipline is stated in the module/subject docs — warnings are errors by
design, so keep the code warning-clean.

## Review with Anki (spaced repetition)

Every exercise carries Q&A review metadata (`subjects/*/ex*/exercise.json` → `answers[]`).
`tools/anki/export.py` turns those into **tab-separated Anki decks**, one per module, with
front cards tagged `[<module> · <exNN> · <title>]` so you can batch-study per week's
modules. Quizzing happens **only** in Anki — it is not part of exercise grading.

- **Where:** `tools/anki/decks/<MODULE>.txt` (regenerated; e.g. `M12-raft.txt`).
- **Import:** Anki → File → Import → pick a deck file (Fields separated by **Tab**).
  Optionally rename the imported deck to `systems-forge::<Module>` and tag notes with
  `systems-forge`.
- **Regenerate:** `python3 tools/anki/export.py` rewrites all 18 decks from the canon.
  If you add/change review questions, regen and re-import (delete old notes first to avoid
  duplicating drifted cards).
- **Details & caveats:** `tools/anki/README.md`.

## Self-check (for you and for the platform)

`make selfcheck` (or `./bin/forge selftest`) runs a pass/fail fixture pair for every
grading method and asserts each scores as expected — the platform validating itself.