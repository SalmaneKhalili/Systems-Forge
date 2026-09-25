# M2-ex02 · the spawn idiom

The first exercise proved that a parent can reap a child; this one replaces the child's
image with a requested command. The interesting path is the boundary: `execvp` only returns
when it fails, and the parent must turn that failure into the conventional status 127. You
deliver the complete `main.c` and preserve both the success and failure transcripts.

## Shape

Write `main.c` so `make all` produces `./test`, a program that takes a command and its
arguments on the command line and spawns it the Unix way:

1. `fork()` a child;
2. in the child, `execvp(cmd, argv + 1)` — `argv[1]` is the program, `argv[1..]` is its
   argument vector;
3. exec theoretically never returns: the only way you reach the next statement is a
   failure — then print `exec failed for <cmd>` and `return 127`;
4. the parent `waitpid`s the child; when the child exited normally with status 127 (that is,
   an exec failure was *reported*), it prints `parent: exec failed and was reported (127)`;
   when the child exited normally with status 0, it prints `parent: exec ok`. Then the parent
   **exits with the child's exit status** — 0 on success, 127 on exec failure. This is how
   a shell propagates a child's status outward.

Use `fork`, `execvp`, `waitpid`, and `WIFEXITED`/`WEXITSTATUS`. No `system(3)` and no
`popen(3)`: the point is the exec boundary. Forward the program name as `argv[0]` of the
exec'd image by forwarding `&argv[1]`; do not hard-code the argument list. On the failure
path, after `execvp` returns, `cmd` is `argv[1]` of the *current* process — print it as
knowledge, then return 127. No `perror` (stderr must stay empty). All output is on stdout;
exit statuses are 0 on successful spawn and 127 on exec failure, and the parent must
report the 127 it observed. No `CFLAGS` redefines.

## Acceptance

`forge` runs `./test /bin/echo hello`; it must exit 0 and print:

```text
hello
parent: exec ok
```

It then runs the bogus command:

```text
./test /definitely/not/a/command
```

That run must print exactly:

```text
exec failed for /definitely/not/a/command
parent: exec failed and was reported (127)
```

and exit with status **127**. Both runs are diffed against their expected files.

- The successful run proves `execvp` replaced the child image and the parent observed 0.
- The failure run proves the exec-failure path ran, was reported, and propagated as 127.
- Neither run writes stderr or leaves a zombie; both must wait, and output is deterministic.
- `quiz.txt` is complete.

Graded `build` + `quiz`: `build` has two runs, success path (diff + exit 0) and failure
path (diff + exit 127), each separately. The failure run cannot merge its output into the
success output. `quiz.txt` answers must match.

## Readings

- **Reading ladder** — the two syscalls first (that is the whole idiom), then the book's
  chapter, then the reaping note.
- `man 2 fork`, `man 2 execve`, `man 3 execvp` (note: `execvp` searches `PATH`), `man 2 _exit`.
- The Linux Programming Interface, §27.1 "Executing a New Program: execve()", §27.2 "The exec()
  Library Functions" (execve and friends), §28.3 "Speed of Process Creation" (why fork/exec is
  fast) — skim; §26.1 "Waiting on a Child Process" for conflict-free reaping.
- exec(3) family reference: https://man7.org/linux/man-pages/man3/exec.3.html

## Quiz

1. On success, what does exec never do?
2. Which exec family function searches PATH for the program?
