# M4-ex02 · mutex bomb

## Goal

Write `main.c` so `make all` produces `./test`: two threads, one shared counter, 400,000
increments, exactly zero lost. The counter is `volatile long tally`; each thread does
200,000 `tally++` inside `pthread_mutex_lock`/`unlock`. When both threads are joined the
total must be **exactly 400,000** — every single increment counted.

```
total 400000
locks ok
ALL PASS
```

## Constraints

- `tally` is `volatile`: a stray unsynchronized version must not be silently optimized away
  from being a data race (a race here is caught fatally, but `volatile` makes it
  unavoidable, not "depends on optimization").
- The only synchronization is `pthread_mutex_t` with `pthread_mutex_lock`/`unlock` around
  the increment. No atomics, no `sleep`, no spinning.
- One mutex, initialized statically (`PTHREAD_MUTEX_INITIALIZER`), shared by both threads.
- `main` prints only after both `pthread_join`s — nothing printed from inside a thread.
- `-pthread` on the cc lines only (never redefine `CFLAGS`/`LDFLAGS`).

## Acceptance criteria

- [ ] `total 400000` exactly, printed once, from `main`
- [ ] exit 0, empty stderr
- [ ] ThreadSanitizer silent (locks serialize every increment by design)
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which function blocks until the caller owns a mutex?: <answer>
What happens to other threads when one thread owns a mutex?: <answer>
```

## Readings

- The Linux Programming Interface, Chapter 30 "Threads: Thread Synchronization" (§30.1
  "Protecting Accesses to Shared Variables: Mutexes"):
  why plain increments lose updates, the lock/unlock discipline, and what else a mutex
  protects.
- `man pthread_mutex_lock` — the blocking semantics ("returns only after the caller owns
  the mutex") is the whole exercise.
- TLPI's classic figure of two threads and one `long` with lost updates — the exact bug
  this exercise eliminates.

## How you are graded

- `build`: strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) + `./test`
  stdout diff, exit 0, empty stderr, under ThreadSanitizer. Drop the lock: the increments
  race, TSan aborts (exit 66) → FAIL — even if the total happens to come out right that
  run, the race itself is fatal.
- `quiz`: `quiz.txt` answers must match.