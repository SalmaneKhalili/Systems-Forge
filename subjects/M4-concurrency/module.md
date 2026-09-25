# M4 · Concurrency

M4 turns one process into several cooperating threads and makes their shared state deterministic. ex01 establishes thread creation, exit, and join; ex02 protects a counter with a mutex; ex03 uses an atomic for the same total; ex04 uses a read-write lock for concurrent readers; and ex05 combines those disciplines in a bounded producer–consumer queue. The gate runs under ThreadSanitizer, so a data race is an immediate failure rather than a timing-dependent result.

## The build

- **ex01 · join and exit** — whole `main.c`: spawn five threads, finish them with `i * 2`, exercise `pthread_exit`/`return`, join them, and verify the values.
- **ex02 · mutex bomb** — whole `main.c`: two threads, one `volatile` counter, exactly-total increments under `pthread_mutex_t`.
- **ex03 · atomic counter** — whole `main.c`: the same race-free total with `_Atomic` + `atomic_fetch_add`, no lock.
- **ex04 · readers and writers** — whole `main.c`: `pthread_rwlock_t`, invariant-checking readers, version-bumping writers.
- **ex05 · Gate: bounded queue** — `queue.c`: thread-safe FIFO with mutex + condvar, blocking push/pop, cap 16.

## Rules

- There are two shapes. Four exercises are full programs of yours (one `main.c` + `Makefile`); the gate gives you a `queue.h` + a harness `main.c`, and you implement the queue in your own `queue.c` (never edit the harness or header).
- You write every `Makefile` (`all` → `./test`, `fclean`, `re`; never redefine `CFLAGS`/`LDFLAGS`). Thread programs must link Posix threads: add `-pthread` to the compile and link commands, never to `CFLAGS`.
- `forge` compiles with `-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread` and runs the program under **ThreadSanitizer**, with ASLR fixed so its shadow memory maps; the grader handles that. A data race of any kind — one lock held by the wrong thread, a missed atomic, a counter read without its lock — is detected and fatal: TSan aborts the run. A race is not a slowly-bitten bug here; it is an instant FAIL.
- The determinism rule from M2 still rules: no PIDs/timestamps in output, no sleeps, and aggregate results are printed only by `main` after every thread is joined, never inside a thread body.
- There is no “usually works” in this module: the grader flips the thread count and order, and TSan still has to stay silent.

## Prerequisites

Before starting M4, you should be comfortable with everything from M0–M3, plus:

- Explain what a data race is: two threads writing the same variable without synchronization, and why the result depends on timing.
- Call `pthread_create(&tid, NULL, fn, arg)` and `pthread_join(tid, NULL)`.
  (If you've never used pthreads, the exercises walk you through it — but you should know
  that threads exist and that `pthread_create` takes a function pointer.)
- Explain what a mutex is: a lock that only one thread can hold at a time.
- Explain `lock`/`unlock` semantics: if thread A holds the lock, thread B blocks until A
  releases it.
- Compile with `-pthread` in your Makefile (M0/M1 taught Makefiles; this flag is new).
- Read a Makefile `CFLAGS` line and add a new flag to it.

You do NOT need to know: read-write locks, atomics, condition variables, or producer–consumer
patterns. You will build all of those here from scratch.

## So what? (interview / portfolio)

Concurrency bugs are among the highest-frequency failures in real systems, and this module
makes race-free behavior a verifiable claim rather than a hope. The gate — a bounded
producer-consumer queue driven under ThreadSanitizer — is the artifact an interviewer means
when asking you to design a thread-safe queue and prove it.

**Interview questions this module arms you for:**
- Mutex vs. read-write lock vs. atomic: when is each the right tool?
- What is a condition variable doing that a spin-lock can't, and why does it need a mutex?
- How does a data race differ from a logical race, and how does TSan detect one?
- Why must a producer-consumer queue block (not spin) when full or empty?

**Portfolio artifact:** M4-ex05 `queue.c` — a bounded, race-free FIFO (mutex + condvar,
cap 16) proven clean under ThreadSanitizer.