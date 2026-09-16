# M4-ex03 · atomic counter

## Goal

Write `main.c` so `make all` produces `./test`: the same two-thread, 400,000-increment
race-free total as ex02 — but with **no lock at all**, because the counter is atomic.

```c
static atomic_ulong tally;              /* from <stdatomic.h> */
...
atomic_fetch_add(&tally, 1UL);          /* one atomic read-modify-write, no mutex */
```

Two threads of 200,000 `atomic_fetch_add` each. Joined total must be exactly 400,000.

```
total 400000
atomic ok
ALL PASS
```

## Constraints

- `tally` is `_Atomic` (use `atomic_ulong` or `_Atomic unsigned long`); the increment is
  `atomic_fetch_add` — never `tally++` on a plain `long`, never a lock from ex02.
- No `pthread_mutex_*`, no `pthread_rwlock_*`, no `sleep`, no spinning.
- `main` prints only after both joins.
- `-pthread` on the cc lines only (never redefine `CFLAGS`/`LDFLAGS`).

## Acceptance criteria

- [ ] `total 400000` exactly, printed once, from `main`
- [ ] exit 0, empty stderr
- [ ] ThreadSanitizer silent — C11 atomics are recognized as proper synchronization, so a
      real atomic build passes clean even though it is lock-free
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which C11 storage qualifier makes a variable updated atomically without locks?: <answer>
Which function atomically adds a value and returns the old value?: <answer>
```

## Readings

- The Linux Programming Interface, Chapter 30 "Threads: Thread Synchronization" —
  the "Atomics" material: the cases where locks are the wrong tool and C11 atomic
  operations take over (and their memory-ordering caveats).
- cppreference *C atomics*: https://en.cppreference.com/w/c/atomic — `atomic_fetch_add`,
  `atomic_ulong`, and the seq_cst default.
- `man 7 pthreads` "handling shared data" cross-reference — atomics are the lock-free
  alternative it hints at (but never legislates).

## How you are graded

- `build`: strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) + `./test`
  stdout diff, exit 0, empty stderr, under ThreadSanitizer. Substitute a plain `long` with
  `tally++`: the increments race, TSan aborts → FAIL. Note: TSan does not flag an `_Atomic`
  access — that is exactly what makes the correct solution pass.
- `quiz`: `quiz.txt` answers must match.