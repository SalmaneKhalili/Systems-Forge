# M5-ex01 · descriptors

## Goal

Write `main.c` so `make all` produces `./test` — a tour of the **file descriptor table**.
Your program has three handles open before it starts (let alone what the shell handed it):
standard input, standard output, standard error. Every `open` you do buys the *lowest
free number*; every `close` gives it back.

```
stdin 0
stdout 1
stderr 2
first open 3
reuse after close 3
two opens 3 and 4
lowest available 3
ALL PASS
```

- print `fileno(stdin)`, `fileno(stdout)`, `fileno(stderr)`;
- `open("/dev/null", O_RDONLY)` once → print `first open <fd>`, then `close` it;
- open `/dev/null` again → print `reuse after close <fd>`;
- open twice in a row → print them **ascending** as `two opens <a> and <b>`;
- close both, open once more → print `lowest available <fd>`;
- every `open`/`close` result must be checked and every fd you own must be closed.

## Constraints

- fd numbers are exact and deterministic (0, 1, 2, then 3, 4…) — no guessing, no
  `printf("3")` strings: read the value out of `open`.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan. Never redefine `CFLAGS`/`LDFLAGS`.

## Acceptance criteria

- [ ] `stdin 0`, `stdout 1`, `stderr 2`
- [ ] `first open 3`, `reuse after close 3`, `two opens 3 and 4`, `lowest available 3`
- [ ] exit 0, empty stderr, every fd closed before exit
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
What is a file descriptor number?: <answer>
Which integer does the first open() usually return?: <answer>
```

## Readings

- `man 2 open`, `man 2 close` — the "lowest-numbered file descriptor not currently open"
  contract, straight from the manual.
- TLPI §5.3 "File Descriptors Are Not Magic" — the table of numbers, and why storing them
  in shell redirections (`>`, `<`) is just `open` + table surgery done for you.
- The Linux Programming Interface, Chapter 5 "File I/O: The Universal I/O Model".

## How you are graded

- `build`: strict compile (`-std=gnu11 -Wall -Wextra -Werror` + ASan/UBSan) + `./test`
  run, stdout diff, exit 0, empty stderr. Hard-code the numbers: the diff sees
  `first open 3`… wait, no — the diff sees what you print, so hard-coding looks right —
  but the *quiz* and every later module assume you know where numbers come from. A real
  (buggy) reference that misorders the two simultaneous opens fails with a wrong `two
  opens` line.
- `quiz`: `quiz.txt` answers must match.