# M4-ex01 · join and exit

M4 starts with thread lifetime. ex01 makes five threads return values, then joins them in
creation order, establishing the ownership discipline that the mutex, atomic, rwlock, and
queue exercises build on. Write `main.c` so `make all` produces `./test`: the program must
prove both ways a thread can end and report every value only from `main`.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test`. Thread *i* (for `i` in 1..5) must finish with the value **`i * 2`**:

- thread 2 must end via `pthread_exit((void *)val)`;
- the other four may finish via `return (void *)val` (same thing under the hood).

`main` joins every thread in creation order and prints each value as
`thread <i> returned <i*2>`, then prints `join errors 0`. The five values add to 30; verify
that sum and print `sum 30`.

Each thread gets its index as its argument (`(intptr_t)`/`(void *)` round-trip); never
share a mutable global that the threads read or write. `pthread_join`'s second argument is
how the value travels out, so there is no `printf` inside a thread and no shared buffer.
Print only from `main`, strictly after each `join`; `joins` happen in creation order so the
output is deterministic. Pass an argument to every thread, avoiding the NULL-argument trap,
and always check `pthread_create` results.

Your `Makefile` compiles with `-pthread` added to the cc commands (never redefine
`CFLAGS`/`LDFLAGS`). There are no `CFLAGS` redefines.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
thread <i> returned <i*2>
join errors 0
sum 30
```

- The five `thread <i> returned <i*2>` lines are 2, 4, 6, 8, and 10 in creation order, and only `main` prints them.
- Thread 2 ends via `pthread_exit`; the other four end via `return`; `main` reaps all five with `join`.
- Every `pthread_join` returns 0, producing `join errors 0`; the five values are verified to sum to 30.
- Exit is 0 and stderr is empty. ThreadSanitizer is silent because no data race is possible by construction, and `quiz.txt` is complete (see below).

The `build` grade is a strict compile (`-std=gnu11 -Wall -Wextra -Werror -fsanitize=thread`) plus the `./test` stdout diff, exit 0, and empty stderr, all under ThreadSanitizer. A thread that reads a shared index instead of its own argument is a data race: TSan aborts → FAIL. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- The Linux Programming Interface, Chapter 29 "Threads: Introduction" — creating,
  terminating, and joining threads (joinable vs detached, `pthread_exit`, exit-value
  plumbing).
- `man 7 pthreads`, `man pthread_create`, `man pthread_join`, `man pthread_exit`.
- The `(void *)`/`intptr_t` value-passing idiom is described in TLPI Chapter 29.

## Quiz

1. Which pthread call reclaims a thread and captures its exit value?
2. Which function terminates the calling thread and returns a value to pthread_join?