# M2-ex04 · signals

## Goal

Write `main.c` so `make all` produces `./test`, a program that installs a handler for
`SIGUSR1` and *proves* it ran:

1. install a handler via `signal(SIGUSR1, on_sigusr1)` — check `SIG_ERR`;
2. print `installing`;
3. deliver the signal to yourself with `raise(SIGUSR1)` — check the return;
4. after `raise` returns, the handler has already completed (a synchronous delivery), so
   inspect the flag the handler set:

```
installing
caught SIGUSR1
```

The handler's only job is `got = 1` on a `volatile sig_atomic_t`. If you had not installed
the handler, the default `SIGUSR1` action would kill the process — that is half the test.

## Constraints

- Handler body: one assignment to a `volatile sig_atomic_t` flag, then return. **No
  `printf`, no `sleep`, no allocation inside the handler.** (printf is not async-signal-safe:
  putting it in a handler is exactly the kind of bug the module teaches you to avoid.)
- Use `signal()` or `sigaction()` — either is fine; check every return and treat
  `SIG_ERR`/`-1` as failure paths.
- Print `installing` once, before `raise`; print `caught SIGUSR1` (or the failure text) only
  after `raise` returns.
- No `CFLAGS` redefines. Empty stderr, exit 0 on success.

## Acceptance criteria

- [ ] `./test` prints exactly the two lines above and exits 0
- [ ] the flag is read *after* `raise` returns — the handler demonstrably ran
- [ ] the default disposition is not enough: with no handler installed the process dies and
      the run fails, so test your program both ways
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which integer type is read and written atomically inside a signal handler?: <answer>
Which function sends a signal to the current process?: <answer>
```

## Readings

- The Linux Programming Interface, §21.1 "Establishing a Handler", §21.2
  "Reentrant and Async-Signal-Safe Functions" (why no printf).
- `man 2 signal`, `man 2 raise`, `man 7 signal-safety`.
- C standard note on `sig_atomic_t`: cppreference
  https://en.cppreference.com/w/c/atomic/sig_atomic_t

## How you are graded

- `build`: strict compile; stdout diff, exit 0, empty stderr. Async-signal-safety of your
  handler is part of the discipline — `-Werror` won't catch it, but the readings will teach
  it and your test double-check (`./test` without installing) will prove the disposition.
- `quiz`: `quiz.txt` answers must match.