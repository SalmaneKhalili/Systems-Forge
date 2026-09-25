# M4-ex02 · mutex bomb

With thread lifetime established, ex02 puts a shared mutable value in front of two
threads. Each thread performs 200,000 increments of one `volatile long tally`, and the
mutex must make the joined total exact rather than merely plausible. Write the whole
`main.c` so `make all` produces `./test`; `main` reports the result only after both joins.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test`. The shared counter is `volatile long tally`; each of two threads does 200,000
`tally++` inside `pthread_mutex_lock`/`unlock`. After both threads are joined, the total
must be **exactly 400,000** — every increment counted.

`tally` is `volatile`: a stray unsynchronized version must not be silently optimized away
from being a data race. A race here is caught fatally, but `volatile` makes it unavoidable,
not dependent on optimization. The only synchronization is one `pthread_mutex_t`, initialized
statically with `PTHREAD_MUTEX_INITIALIZER` and shared by both threads. There are no atomics,
no `sleep`, and no spinning. `main` prints only after both `pthread_join`s; nothing is
printed from inside a thread.

Add `-pthread` on the cc lines only, never redefining `CFLAGS`/`LDFLAGS`.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
total 400000
locks ok
ALL PASS
```

- `total 400000` is printed once, from `main`, after both joins; every increment is counted.
- `locks ok` and `ALL PASS` complete the transcript without output from either thread.
- Exit is 0 and stderr is empty. ThreadSanitizer stays silent because the mutex serializes every increment by design, and `quiz.txt` is complete (see below).

The `build` grade is a strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) plus the `./test` stdout diff, exit 0, and empty stderr, under ThreadSanitizer. Drop the lock and the increments race, so TSan aborts (exit 66) → FAIL, even if that run's total happens to come out right. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- The Linux Programming Interface, Chapter 30 "Threads: Thread Synchronization" (§30.1
  "Protecting Accesses to Shared Variables: Mutexes"):
  why plain increments lose updates, the lock/unlock discipline, and what else a mutex
  protects.
- `man pthread_mutex_lock` — the blocking semantics ("returns only after the caller owns
  the mutex") is the whole exercise.
- TLPI's classic figure of two threads and one `long` with lost updates — the exact bug
  this exercise eliminates.

## Quiz

1. Which function blocks until the caller owns a mutex?
2. What happens to other threads when one thread owns a mutex?