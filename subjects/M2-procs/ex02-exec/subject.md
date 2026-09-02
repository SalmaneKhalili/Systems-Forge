# M2-ex02 · the spawn idiom

## Goal

Write `main.c` so `make all` produces `./test`, a program that takes a command and its
arguments on the command line and spawns it the Unix way:

1. `fork()` a child;
2. in the child, `execvp(cmd, argv + 1)` — `argv[1]` is the program, `argv[1..]` is its
   argument vector;
3. exec *theoretically never returns*: the only way you reach the next statement is a
   failure — then print `exec failed for <cmd>` and `return 127`;
4. the parent `waitpid`s the child; when the child exited normally with status 127 (that is,
   an exec failure was *reported*), it prints `parent: exec failed and was reported (127)`;
   when the child exited normally with status 0, it prints `parent: exec ok`. Then the parent
   **exits with the child's exit status** — 0 on success, 127 on exec failure. This is how a
   shell propagates a child's status outward.

`forge` runs `./test /bin/echo hello` (must produce `hello` followed by `parent: exec ok`,
exit 0) and runs it a second time with a bogus command:

```
./test /definitely/not/a/command
```

which must print `exec failed for /definitely/not/a/command`, then
`parent: exec failed and was reported (127)`, and exit with status **127**. Both runs are
diffed against their expected files.

## Constraints

- Use `fork`, `execvp`, `waitpid`, `WIFEXITED`/`WEXITSTATUS`. **No `system(3)`, no
  `popen(3)`** — the point is the exec boundary.
- cmd argument vector: pass the program name as `argv[0]` of the exec'd image by forwarding
  `&argv[1]`. Do not hard-code the argument list.
- exec failure path: after `execvp` returns, that `cmd` is `argv[1]` of the *current* process
  — print it as knowledge, then return 127. No `perror` (stderr must stay empty).
- All output on stdout; exit statuses: 0 on the successful spawn, 127 on exec failure, and
  the parent must report the 127 it observed.
- No `CFLAGS` redefines.

## Acceptance criteria

- [ ] `./test /bin/echo hello` prints exactly:
      ```
      hello
      parent: exec ok
      ```
      and exits 0.
- [ ] `./test /definitely/not/a/command` prints `exec failed for /definitely/not/a/command`,
      then `parent: exec failed and was reported (127)`, and **the program exits 127**
      (the parent's exit status mirrors the child's)
- [ ] no stderr, deterministic, no zombies (both must wait)
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
On success, what does exec never do?: <answer>
Which exec family function searches PATH for the program?: <answer>
```

## Readings

- The Linux Programming Interface, §27.2 "Executing a Program" (execve and friends),
  §24.5 "Optimizing fork()" — skim; §26.1 "Waiting on a Child Process" for conflict-free
  reaping.
- `man 2 execve`, `man 3 execvp` (note: `execvp` searches `PATH`), `man 2 _exit` vs `exit`.
- "The exec() family": https://man7.org/linux/man-pages/man3/exec.3.html

## How you are graded

- `build` with two runs: success path (diff + exit 0) and failure path (diff + exit 127).
  Each run separately. The failure run proves you actually code the exec-failure path — you
  cannot merge the two outputs into one.
- `quiz`: `quiz.txt` answers must match.