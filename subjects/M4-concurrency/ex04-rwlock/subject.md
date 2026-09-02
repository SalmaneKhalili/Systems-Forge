# M4-ex04 · readers and writers

## Goal

Write `main.c` so `make all` produces `./test`: 4 reader threads + 2 writer threads over a
shared two-variable record. This is the "80% reads, 20% writes" shape that a plain mutex
serializes unfairly — so the record is guarded by `pthread_rwlock_t`.

- Shared record: `long a`, `long b` with the invariant **`b == 2 * a`**, plus an
  `int version`.
- Writer *w* (2 threads × 10 rounds): `wrlock`, `a++`, `b = 2*a`, `version++`, unlock.
- Reader *r* (4 threads × 50 rounds): `rdlock`, check the invariant and count it, unlock.
- The counting variables (`reads`, `consistent`) are themselves guarded by a small plain
  `pthread_mutex_t` — never by the rwlock (readers hold that *concurrently*).

Aggregate lines, printed by `main` after every join:

```
read rounds 200
write rounds 20
consistent 200
version 20
ALL PASS
```

## Constraints

- The rwlock's whole purpose: readers check the invariant in parallel. No `sleep`, no
  spinning, no timing games — correctness is by synchronisation, not luck.
- `a`,`b`,`version` are touched only under the rwlock (readers: `rdlock`, writers:
  `wrlock`). The counters use their own mutex because all four readers write them while
  holding only the *read* lock.
- `main` prints only after every thread has been joined.
- `-pthread` on the cc lines only (never redefine `CFLAGS`/`LDFLAGS`).

## Acceptance criteria

- [ ] `read rounds 200`, `write rounds 20`, `consistent 200`, `version 20` — all exact
- [ ] exit 0, empty stderr
- [ ] ThreadSanitizer silent on the whole program — readers and writers, and the counters
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which call acquires a read-write lock for reading?: <answer>
Can two reader threads hold a read-write lock at the same time?: <answer>
```

## Readings

- The Linux Programming Interface, Chapter 30 "Threads: Synchronization and Common Mistakes" —
  the read-write lock material: rdlock/wrlock semantics and the "readers are allowed to run
  concurrently" contract.
- `man pthread_rwlock_rdlock`, `man pthread_rwlock_unlock`.
- TLPI Chapter 30 recap: why the rwlock exists — a plain mutex would make every reader wait for
  every other reader, which is precisely the cost this lock removes.

## How you are graded

- `build`: strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) + `./test`
  stdout diff, exit 0, empty stderr, under ThreadSanitizer. Skip the locks entirely: the
  invariant checks race with the writers → TSan aborts (and `consistent` stops equalling
  200) → FAIL. Guard the counters with the rwlock instead of their own mutex: the four
  readers race while holding the read lock → TSan aborts → FAIL.
- `quiz`: `quiz.txt` answers must match.