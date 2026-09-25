# M4-ex04 · readers and writers

After the counter exercises, ex04 gives readers and writers different access to one
invariant. Four readers may check the record concurrently while two writers update it, so
`pthread_rwlock_t` provides the shared-state boundary that ex05 later combines with
condition variables. Write the whole `main.c` so `make all` produces `./test`, with all
aggregate lines emitted by `main` after every join.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test`. Four reader threads and two writer threads operate on a shared record containing
`long a`, `long b`, and `int version`, with the invariant **`b == 2 * a`**.

- Writer *w* (2 threads × 10 rounds): `wrlock`, `a++`, `b = 2*a`, `version++`, unlock.
- Reader *r* (4 threads × 50 rounds): `rdlock`, check the invariant and count it, unlock.

The counting variables (`reads`, `consistent`) are guarded by a small plain
`pthread_mutex_t`, never by the rwlock: all four readers write them while holding only the
*read* lock. The rwlock's purpose is to let readers check the invariant in parallel. There
are no `sleep`, spinning, or timing games; correctness comes from synchronisation, not
luck. `a`, `b`, and `version` are touched only under the rwlock (readers: `rdlock`,
writers: `wrlock`). `main` prints only after every thread has been joined.

Add `-pthread` on the cc lines only, never redefining `CFLAGS`/`LDFLAGS`.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
read rounds 200
write rounds 20
consistent 200
version 20
ALL PASS
```

- `read rounds 200` and `write rounds 20` report the fixed 4×50 and 2×10 workloads exactly; `consistent 200` and `version 20` prove the invariant and writer progress.
- Readers and writers use the rwlock, while `reads` and `consistent` use their own mutex; the whole program exits 0 with empty stderr and ThreadSanitizer silent.
- Complete `quiz.txt` (see below).

The `build` grade is a strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) plus the `./test` stdout diff, exit 0, and empty stderr, under ThreadSanitizer. Skipping the locks makes the invariant checks race with the writers, so TSan aborts and `consistent` stops equalling 200 → FAIL. Guarding the counters with the rwlock instead of their own mutex makes the four readers race while holding the read lock, so TSan aborts → FAIL. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- The Linux Programming Interface, Chapter 30 "Threads: Thread Synchronization" (§30.1
  mutex discipline underpinning the read-write lock) —
  the read-write lock material: rdlock/wrlock semantics and the "readers are allowed to run
  concurrently" contract.
- `man pthread_rwlock_rdlock`, `man pthread_rwlock_unlock`.
- TLPI Chapter 30 recap: why the rwlock exists — a plain mutex would make every reader wait for
  every other reader, which is precisely the cost this lock removes.

## Quiz

1. Which call acquires a read-write lock for reading?
2. Can two reader threads hold a read-write lock at the same time?