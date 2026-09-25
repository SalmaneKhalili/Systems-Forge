# M4-ex05 · Gate: Bounded Queue

The gate combines M4's lifecycle and locking rules in a bounded producer–consumer queue.
A producer and consumer hammer the same capacity-16 queue from separate threads, so the
implementation must preserve every value and FIFO order while blocking, not spinning, on
full and empty states. Implement the provided contract in `queue.c`; the harness supplies
`queue.h` and `main.c`.

## Shape

The gate is a harness shape: the provided `main.c` (the harness) and `queue.h` are not
editable. You implement `queue.c` with:

```c
Queue *queue_new(size_t cap);   /* allocate+init; NULL on failure */
void   queue_free(Queue *q);    /* release everything */
int    queue_push(Queue *q, int v);  /* 0, blocking while full */
int    queue_pop(Queue *q, int *out); /* 0, blocking while empty */
```

The producer pushes 0…999 and the consumer pops them through a capacity-16 queue. The
queue must deliver every value exactly once and in FIFO order, with no data race under
ThreadSanitizer. One `pthread_mutex_t` guards `count`/`head`/`tail`/buffer. Two
`pthread_cond_t` track the two wait situations: `not_full` (producers wait here) and
`not_empty` (consumers wait here).

Use the block/wait pattern: `pthread_mutex_lock`, then
`while (condition) pthread_cond_wait(&cv, &mutex);` — the `while` matters because a
spurious wakeup must re-check the condition. When a producer adds a value it signals
`not_empty`; when a consumer takes one it signals `not_full`. A forgotten signal leaves a
thread waiting forever. The `Queue` struct is yours and lives in `queue.c` because the
header hides it. `queue_cap` is 16 in the harness and `ITEMS` is 1000, so the queue is
exercised full and empty many times over.

There are no `sleep`, busy-wait loops for the wait side, or atomics-only shortcuts. Add
`-pthread` on the cc lines only, never redefining `CFLAGS`/`LDFLAGS`. `forge` compiles
`main.c` (provided harness) + `queue.c`, runs `./test`, and diffs stdout.

## Acceptance

A clean run must have `queue_new(16)` succeed and start with an empty queue, then produce and
consume 1000 values with exit 0 and empty stderr. `fifo ok` means value `i` comes out before
value `i+1`; the values are batched into `checksum 499500`. ThreadSanitizer must be completely
silent. The full and empty wait paths must actually block and wake, with no spin, so the run
finishes quickly instead of deadlocking; complete `quiz.txt` (see below).

The `build` grade is a strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) plus the `./test` stdout diff, exit 0, and empty stderr, under ThreadSanitizer. A queue with no mutex lets the producer and consumer race on `count`/buffer, so TSan aborts → FAIL. A stack posing as a queue makes the first pop 999 rather than 0, producing `fifo violated` → FAIL. A condvar that is never signalled leaves the consumer waiting forever, so the run times out (20s) → FAIL. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- **Reading ladder** — start with `man pthread_cond_wait` and
  `man pthread_cond_signal`, then TLPI §30.2 for the wait loop and
  spurious wakeups.
- The Linux Programming Interface, §30.2 "Signaling Changes of State: Condition Variables"
  (esp. §30.2.2 "Signaling and Waiting on Condition Variables") — the
  lock→wait→resignal loop in figure 30-2, the `while` loop contract, and the
  spurious-wakeup rationale, plus the wait-mechanism overview.
- `man pthread_cond_wait`, `man pthread_cond_signal` — "the mutex is released and remains
  unlocked until another thread signals".
- Producer-consumer with bounded buffer is TLPI §30.2.3's classic example figure 30-2.

## Quiz

1. Which function waits on a condition variable while releasing its mutex?
2. Which function wakes one thread waiting on a condition variable?