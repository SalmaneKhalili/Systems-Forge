# M4-ex01 · join and exit

## Goal

Write `main.c` so `make all` produces `./test`: a program that creates 5 threads, makes
each one finish with a value, and reaps them all — proving both ways a thread can end.

- Thread *i* (for `i` in 1..5) must finish with the value **`i * 2`**:
  - thread 2 must end via `pthread_exit((void *)val)`;
  - the other four may finish via `return (void *)val` (same thing under the hood).
- `main` joins every thread in creation order and prints each value:
  `thread <i> returned <i*2>`, then a clean `join errors 0`.
- The five values add to 30; verify and print `sum 30`.

## Constraints

- Each thread gets its index as its argument (`(intptr_t)`/`(void *)` round-trip); never
  share a mutable global that the threads read or write.
- `pthread_join`'s second argument is how the value travels out — no `printf` inside a
  thread, no shared buffers.
- Print only from `main`, strictly after each `join`; `joins` happen in creation order so
  output is deterministic.
- Pass an argument to every thread (avoid the NULL-argument trap), and always check
  `pthread_create` results.
- Your `Makefile` compiles with `-pthread` added to the cc commands (never redefine
  `CFLAGS`/`LDFLAGS`). No `CFLAGS` redefines.

## Acceptance criteria

- [ ] 5 threads, values 2, 4, 6, 8, 10 printed in creation order
- [ ] thread 2 ends via `pthread_exit`, the rest via `return` — all reaped by `join`
- [ ] every `pthread_join` returns 0 → `join errors 0`
- [ ] `sum 30` verified, exit 0, empty stderr
- [ ] ThreadSanitizer silent (no data race possible by construction)
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which pthread call reclaims a thread and captures its exit value?: <answer>
Which function terminates the calling thread and returns a value to pthread_join?: <answer>
```

## Readings

- The Linux Programming Interface, Chapter 29 "Threads: Introduction" — creating,
  terminating, and joining threads (joinable vs detached, `pthread_exit`, exit-value
  plumbing).
- `man 7 pthreads`, `man pthread_create`, `man pthread_join`, `man pthread_exit`.
- The `(void *)`/`intptr_t` value-passing idiom is described in TLPI Chapter 29.

## How you are graded

- `build`: strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) + `./test`
  stdout diff, exit 0, empty stderr — all under ThreadSanitizer. A thread that reads a
  shared index instead of its own argument is a data race: TSan aborts → FAIL.
- `quiz`: `quiz.txt` answers must match.