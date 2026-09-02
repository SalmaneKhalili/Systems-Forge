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

The curriculum lives under `subjects/` (read-only canon: specs, scaffolds, harnesses,
expected outputs, quiz questions — never a solved deliverable). Your solutions live under
`answers/` (gitignored). `forge check M11-ex03` runs the autograder for that exercise
inside your answers workspace and records the result in the local SQLite progress DB.

`subjects/` ships the *un*solved exercise (you write the deliverable — e.g. `log.go`,
`server.go`, `cluster.go`, `check.py`, `solve.sh` — from the spec in `subject.md`). A
bare workspace therefore auto-FAILs every exercise. Complete reference solutions live
under `solutions/` (gitignored) so you can compare after you solve one — never before.

## Grading model

Every exercise declares one or more grading methods (see `subjects/*/exercise.json`):

| method     | what it verifies                                                           |
| ---------- | -------------------------------------------------------------------------- |
| `build`    | `make` under `-Wall -Wextra -Werror` (+ sanitizers) and runs test binaries |
| `stdout`   | normalized diff of a program's output against expected                     |
| `artifact` | properties of produced files (type, symbols, sizes, contents)              |
| `process`  | spawns a program, drives it with signals/probes, asserts state             |
| `net`      | drives a scripted protocol over a socket and asserts responses             |
| `quiz`     | checks a Q&A checkpoint file against the exercise's answers                |
| `report`   | parses structured output and compares to a reference run                   |
| `scenario` | runs a fault-injection driver against a multi-node system                  |
| `fault`    | asserts a _property_ under injected failures (partition, crash)            |
| `lincheck` | checks operation histories for linearizability                             |

`make selfcheck` validates the grader itself: it runs a pass/fail fixture pair for
every method and asserts each scores as expected.

## Layout

```
forge/        the platform: TUI + CLI + grader engine (Go)
subjects/     the curriculum, one dir per module, exercises as exNN-<id> (unsolved)
caps/         capstone specifications
switch/       the partitionable transport (built during the curriculum)
answers/      your personal workspace (gitignored) — you write your solutions here
solutions/    reference solutions (gitignored) — compare after you solve
tools/        selfcheck, fixtures, helper graders
progress.db   local progress store (gitignored, SQLite)
```
