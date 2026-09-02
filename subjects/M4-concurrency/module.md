# M4 · Concurrency

One process, many threads. This module makes concurrency *boring*: you create threads, you
wait for them, and you protect shared state with the three tools the C universe gives you —
a mutex, a read-write lock, and an atomic. The gate wires them into a producer–consumer
queue that must ship every byte in order, race-free, no matter how the scheduler interleaves
the threads. After M4, "but it worked on my machine" no longer travels — a data race is a
build error here.

**Rules of the module**

- Two shapes: four exercises are full programs of yours (one `main.c` + `Makefile`); the
  gate gives you a `queue.h` + a harness `main.c` — you implement the queue in your own
  `queue.c` (never edit the harness or header).
- You write every `Makefile` (`all`→`./test`, `fclean`, `re`; never redefine
  `CFLAGS`/`LDFLAGS`). Thread programs must link Posix threads — add `-pthread` to the
  compile and link commands, never to `CFLAGS`.
- `forge` compiles with `-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread` and runs the
  program under **ThreadSanitizer** (with ASLR fixed so its shadow memory maps — the grader
  handles that). A data race of any kind — one lock held by the wrong thread, a missed
  atomic, a counter read without its lock — is *detected and fatal*: TSan aborts the run.
  That is the whole game: a race is not a slowly-bitten bug here, it is an instant FAIL.
- Determinism rule from M2 still rules: no PIDs/timestamps in output, no sleeps, aggregate
  results printed only by `main` after every thread is joined (never inside a thread body).
- There is no "usually works" in this module — the grader flips the thread count and order,
  and TSan still has to stay silent.

| ex | topic | artifact you deliver |
|----|-------|----------------------|
| ex01 | join & exit | whole `main.c`: spawn threads, `pthread_exit`/`return`, join and verify values |
| ex02 | mutex bomb | whole `main.c`: two threads, one `volatile` counter, exactly-total increments under `pthread_mutex_t` |
| ex03 | atomic counter | whole `main.c`: the same race-free total with `_Atomic` + `atomic_fetch_add`, no lock |
| ex04 | readers & writers | whole `main.c`: `pthread_rwlock_t`, invariant-checking readers, version-bumping writers |
| ex05 | **Gate: bounded queue** | `queue.c`: thread-safe FIFO with mutex + condvar, blocking push/pop, cap 16 |