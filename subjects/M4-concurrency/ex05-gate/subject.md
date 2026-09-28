# M4-ex05 gate · bounded queue

## Goal

Implement `queue.c` backing the provided `queue.h`: a **thread-safe bounded FIFO** that two
threads hammer simultaneously — a producer pushing 0…999 and a consumer popping them, both
through a capacity-16 queue. It must deliver:

- every value, exactly once, **in FIFO order**;
- with **no data race** (ThreadSanitizer is watching);
- without busy-spinning — when the queue is full, `queue_push` *waits*; when it is empty,
  `queue_pop` *waits* — a condition variable does the waiting.

```c
Queue *queue_new(size_t cap);   /* allocate+init; NULL on failure */
void   queue_free(Queue *q);    /* release everything */
int    queue_push(Queue *q, int v);  /* 0, blocking while full */
int    queue_pop(Queue *q, int *out); /* 0, blocking while empty */
```

`forge` compiles `main.c` (provided harness) + `queue.c`, runs `./test`, diffs stdout.

## Constraints

- One `pthread_mutex_t` guards `count`/`head`/`tail`/buffer. Two `pthread_cond_t` track
  the two wait situations: `not_full` (producers wait here) and `not_empty` (consumers wait
  here).
- The block/wait pattern is the textbook one: `pthread_mutex_lock`, then
  `while (condition) pthread_cond_wait(&cv, &mutex);` — the `while` matters, a spurious
  wakeup must re-check. When a producer adds a value it must signal `not_empty`; when a
  consumer takes one it must signal `not_full`. A forgotten signal = a thread that waits
  forever.
- The `Queue` struct is yours — it lives in `queue.c` (the header hides it).
- `queue_cap` is 16 in the harness; `ITEMS` is 1000, so the queue is exercised full *and*
  empty many times over.
- No `sleep`, no busy-wait loops for the wait side, no atomics-only shortcut.
- `-pthread` on the cc lines only (never redefine `CFLAGS`/`LDFLAGS`).

## Acceptance criteria

- [ ] `queue_new(16)` succeeds; empty queue starts correct
- [ ] 1000 values produced and 1000 consumed, exit 0, empty stderr
- [ ] `fifo ok` — value `i` comes out before value `i+1`, batched into `checksum 499500`
- [ ] ThreadSanitizer completely silent
- [ ] the full/empty wait paths actually block and wake (no spin), so the run finishes
      quickly, not in a deadlock you mistake for correctness
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which function waits on a condition variable while releasing its mutex?: <answer>
Which function wakes one thread waiting on a condition variable?: <answer>
```

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

## How you are graded

- `build`: strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) + `./test`
  stdout diff, exit 0, empty stderr, under ThreadSanitizer. A queue with no mutex: producer
  and consumer race on `count`/buffer → TSan aborts → FAIL. A stack posing as a queue: the
  first pop is 999, not 0 → `fifo violated` → FAIL. A condvar that is never signalled:
  consumer waits forever → run times out (20s) → FAIL.
- `quiz`: `quiz.txt` answers must match.