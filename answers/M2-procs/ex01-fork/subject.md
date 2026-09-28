# M2-ex01 · fork and wait

## Goal

Write `main.c` so that `make all` produces `./test`, a program that:

1. calls `fork()`,
2. in the child: prints `child: hello from the child` and exits with status **7**,
3. in the parent: calls `waitpid(pid, &st, 0)`, and — only if it reaped that exact child,
   `WIFEXITED(st)` is true, and `WEXITSTATUS(st) == 7` — prints
   `parent: reaped the child, exit status 7`.

`forge` runs `make fclean`, `make all`, `./test` and diffs stdout (whitespace-insensitive).

Expected output, in this exact order:

```
child: hello from the child
parent: reaped the child, exit status 7
```

## Constraints

- Use `fork`, `waitpid`, `WIFEXITED`, `WEXITSTATUS`. No `sleep`, no busy-waiting.
- The child's message must be printed **before** the parent's line (the parent prints only
  after `waitpid` returns — that is what makes the order deterministic).
- Never print a PID or any other value that differs between runs.
- Check every return: `fork() < 0` is a failure path, `waitpid`'s return must equal `pid`.
- `perror` output (stderr) is forbidden here: the grader requires empty stderr. Handle the
  failure paths with `return 1` after a short stdout message, or simply exit.
- No `CFLAGS` redefines.

## Acceptance criteria

- [ ] `./test` prints exactly the two lines above and exits 0
- [ ] the parent actually waited on the forked pid (a blind `wait(NULL)` that skips the
      verify will not match `waitpid`'s return-value check the way you think — test it)
- [ ] exit status 7 crosses the boundary intact (`WIFEXITED` + `WEXITSTATUS`)
- [ ] no stderr, no sanitizer finding, deterministic across runs
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which syscall waits for a specific child process?: <answer>
Which macro tests whether a child exited normally?: <answer>
```

## Readings

- **Reading ladder** — fork is one call, three outcomes; learn its return-value contract first,
  then the waiting discipline, then where book detail lives.
- `man 2 fork`, `man 2 waitpid`, `man 2 _exit` (why a child returns via `_exit` while the
  parent keeps running).
- The Linux Programming Interface, Chapter 24 "Process Creation" (§24.2 "Creating a New
  Process: fork()"), §24.4 "Shared File Descriptors after fork()", §26.1 "Waiting on a Child
  Process", §26.2 "Orphans and Zombies" (wait/waitpid, status macros).

## How you are graded

- `build`: strict compile; `./test` stdout must equal `expected.txt` (whitespace-insensitive)
  and the process must exit 0 with empty stderr. A program that exits 7, prints one line, or
  races fails the diff or the exit check.
- `quiz`: `quiz.txt` answers must match.