# M4-ex03 · atomic counter

ex02 showed what a mutex buys; ex03 keeps the same two-thread total but removes the lock.
The shared counter is a C11 atomic, so each increment is one indivisible read-modify-write
while the program still finishes with exactly 400,000. Write the whole `main.c` so
`make all` produces `./test`, and print only after both joins.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test`. Two threads perform 200,000 increments each, and their joined total must be
exactly 400,000. Use the C11 atomic form:

```c
static atomic_ulong tally;              /* from <stdatomic.h> */
...
atomic_fetch_add(&tally, 1UL);          /* one atomic read-modify-write, no mutex */
```

`tally` is `_Atomic` (use `atomic_ulong` or `_Atomic unsigned long`); the increment is
`atomic_fetch_add`, never `tally++` on a plain `long`, and there is no lock from ex02.
There are no `pthread_mutex_*`, `pthread_rwlock_*`, `sleep`, or spinning operations.
`main` prints only after both joins. Add `-pthread` on the cc lines only, never redefining
`CFLAGS`/`LDFLAGS`.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
total 400000
atomic ok
ALL PASS
```

- `total 400000` is printed once, from `main`, after both joins.
- The counter uses `atomic_fetch_add` rather than a plain increment or a mutex, and the run exits 0 with empty stderr.
- ThreadSanitizer stays silent: C11 atomics are recognized as proper synchronization, so a real atomic build passes clean even though it is lock-free. Complete `quiz.txt` (see below).

The `build` grade is a strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) plus the `./test` stdout diff, exit 0, and empty stderr, under ThreadSanitizer. Substituting a plain `long` with `tally++` makes the increments race, so TSan aborts → FAIL. TSan does not flag an `_Atomic` access; that is exactly what makes the correct solution pass. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- The Linux Programming Interface, Chapter 30 "Threads: Thread Synchronization" —
  the "Atomics" material: the cases where locks are the wrong tool and C11 atomic
  operations take over (and their memory-ordering caveats).
- cppreference *C atomics*: https://en.cppreference.com/w/c/atomic — `atomic_fetch_add`,
  `atomic_ulong`, and the seq_cst default.
- `man 7 pthreads` "handling shared data" cross-reference — atomics are the lock-free
  alternative it hints at (but never legislates).

## Quiz

1. Which C11 storage qualifier makes a variable updated atomically without locks?
2. Which function atomically adds a value and returns the old value?