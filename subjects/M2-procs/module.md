# M2 · Processes

You learned to write functions; now you learn that programs are not the only thing running.
Unix is built on processes: units of execution that `fork` into new ones, `exec` new images,
`wait` for each other, and talk through pipes. Everything from the init system to a web
server is these three verbs composed.

**Rules of the module**

- You write the whole program — `main.c` — and your own `Makefile` (same contract as M1:
  `all`→`./test`, `fclean`, `re`; never redefine `CFLAGS`/`LDFLAGS`).
- `forge` compiles with `-std=gnu11 -Wall -Wextra -Werror -fsanitize=address,undefined` and
  runs `./test` zero, one, or more times (with arguments), diffing each invocation against an
  expected file (`normalize: strip`), and requiring a clean exit and empty stderr.
- All exercises must be **deterministic**: no printing raw PIDs, no races, no sleeps-for-luck.
  Design so the output order is forced by the code (wait for the child, then print).
- Never modify the `C` discipline learned in M1: every resource you open is closed, every
  return checked, every child reaped.

| ex | topic | artifact you deliver |
|----|-------|----------------------|
| ex01 | fork + wait | a `fork1`-style `main.c`: fork, child exits 7, parent reaps and verifies status |
| ex02 | the spawn idiom | `main.c`: fork → `execvp` → wait, with a correct exec-failure path (exit 127) |
| ex03 | pipes | `main.c`: child writes through a pipe, parent reads till EOF and reports |
| ex04 | signals | `main.c`: install a `SIGUSR1` handler, `raise`, prove it ran (flag, not printf) |
| ex05 | **Gate: pipeline** | `main.c`: two children, one pipe, `dup2` + `echo`→`tr`, parent reaps both |