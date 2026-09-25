# M2-ex01 · fork and wait

M2 starts by splitting one program into two processes and joining them with a status
handshake. The parent must not guess that the child finished; it must reap the exact PID
and decode the exit status. You deliver `main.c`, and the ordering contract that makes the
two-line result deterministic.

## Shape

Write `main.c` so that `make all` produces `./test`, a program that:

1. calls `fork()`;
2. in the child: prints `child: hello from the child` and exits with status **7**;
3. in the parent: calls `waitpid(pid, &st, 0)`, and — only if it reaped that exact child,
   `WIFEXITED(st)` is true, and `WEXITSTATUS(st) == 7` — prints
   `parent: reaped the child, exit status 7`.

`forge` runs `make fclean`, `make all`, `./test` and diffs stdout (whitespace-insensitive).
Use `fork`, `waitpid`, `WIFEXITED`, and `WEXITSTATUS`; no `sleep` and no busy-waiting. The
child's message must be printed **before** the parent's line, because the parent prints
only after `waitpid` returns. Never print a PID or any other value that differs between
runs. Check every return: `fork() < 0` is a failure path, and `waitpid`'s return must equal
`pid`.

`perror` output (stderr) is forbidden here because the grader requires empty stderr. Handle
failure paths with `return 1` after a short stdout message, or simply exit. No `CFLAGS`
redefines.

## Acceptance

`./test` must print exactly the following two lines, in this order, and exit 0:

```text
child: hello from the child
parent: reaped the child, exit status 7
```

- The parent actually waits on the forked PID and verifies the returned child; a blind
  `wait(NULL)` that skips the verify does not satisfy the `waitpid` return-value check.
- Exit status 7 crosses the process boundary through `WIFEXITED` and `WEXITSTATUS`.
- Stderr is empty, there is no sanitizer finding, and output is deterministic across runs.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile plus `./test` stdout diff against
`expected.txt` (whitespace-insensitive), exit 0, and empty stderr; a program that exits 7,
prints one line, or races fails the diff or exit check. `quiz.txt` answers must match.

## Readings

- **Reading ladder** — fork is one call, three outcomes; learn its return-value contract first,
  then the waiting discipline, then where book detail lives.
- `man 2 fork`, `man 2 waitpid`, `man 2 _exit` (why a child returns via `_exit` while the
  parent keeps running).
- The Linux Programming Interface, Chapter 24 "Process Creation" (§24.2 "Creating a New
  Process: fork()"), §24.4 "Shared File Descriptors after fork()", §26.1 "Waiting on a Child
  Process", §26.2 "Orphans and Zombies" (wait/waitpid, status macros).

## Quiz

1. Which syscall waits for a specific child process?
2. Which macro tests whether a child exited normally?
