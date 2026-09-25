# M2-ex04 · signals

The pipe proves data can synchronize processes; a signal proves control can arrive without
another process. This exercise keeps the handler deliberately tiny and proves delivery with
a flag after `raise` returns. You deliver `main.c`, and the async-signal-safety boundary
that later process supervisors depend on.

## Shape

Write `main.c` so `make all` produces `./test`, a program that installs a handler for
`SIGUSR1` and *proves* it ran:

1. install a handler via `signal(SIGUSR1, on_sigusr1)` — check `SIG_ERR`;
2. print `installing`;
3. deliver the signal to yourself with `raise(SIGUSR1)` — check the return;
4. after `raise` returns, the handler has already completed (a synchronous delivery), so
   inspect the flag the handler set:

```text
installing
caught SIGUSR1
```

The handler's only job is `got = 1` on a `volatile sig_atomic_t`. If you had not installed
the handler, the default `SIGUSR1` action would kill the process — that is half the test.

The handler body is one assignment to a `volatile sig_atomic_t` flag, then return. **No
`printf`, no `sleep`, no allocation inside the handler**: `printf` is not
async-signal-safe, and putting it in a handler is exactly the kind of bug the module
teaches you to avoid. Use `signal()` or `sigaction()`; either is fine, and check every
return while treating `SIG_ERR`/`-1` as failure paths. Print `installing` once before
`raise`, and print `caught SIGUSR1` (or the failure text) only after `raise` returns. No
`CFLAGS` redefines. Empty stderr, exit 0 on success.

## Acceptance

`./test` must print exactly:

```text
installing
caught SIGUSR1
```

and exit 0.

- The flag is read after `raise` returns, so the handler demonstrably ran.
- With no handler installed, the default disposition kills the process and the run fails;
  test your program both ways.
- Stderr stays empty, and `quiz.txt` is complete.

Graded `build` + `quiz`: `build` is a strict compile, stdout diff, exit 0, and empty
stderr. Async-signal-safety is part of the discipline even though `-Werror` will not catch
it; the readings teach it and the test double-check (`./test` without installing) proves
the disposition. `quiz.txt` answers must match.

## Readings

- **Cold-start glossary** — *signal*: a software interrupt delivered to a process; *disposition*:
  the process-defined action for a signal (default / ignore / handler); *default action*: what
  happens with no handler installed — for `SIGUSR1` that is **terminate the process**;
  *delivery*: the moment the disposition runs; *synchronous*: self-inflicted via `raise()`, as
  opposed to arriving from outside the process. A handler must only do **async-signal-safe**
  work (set a flag) — calling `printf` inside it is exactly the class of bug this module
  teaches you to avoid; that is why `volatile sig_atomic_t` exists. With these six words you
  can open TLPI below without flailing.
- The Linux Programming Interface, §20.1 "The Concept of Signals", §20.2 "Types of Standard
  Signals" (default dispositions), §20.4 "Sending Signals: kill() and raise()" — the
  fundamentals that §21 assumes; read this chapter first.
- The Linux Programming Interface, §21.1 "Designing Signal Handlers", §21.1.2 "Reentrant and
  Async-Signal-Safe Functions", §21.1.3 "Global Variables and the sig_atomic_t Data Type".
- `man 2 signal`, `man 2 raise`, `man 7 signal-safety`.
- `sig_atomic_t` in the C standard: https://en.cppreference.com/w/c/program/sig_atomic_t
- **Depth returns later:** terminal-signal handling and `SIGCHLD`-based reaping come back at
  M7-ex04 (graceful shutdown) and M7-ex05 (mini supervisor).

## Quiz

1. Which integer type is read and written atomically inside a signal handler?
2. Which function sends a signal to the current process?
