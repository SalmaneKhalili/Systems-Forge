# M5-ex01 · descriptors

M5 begins by making the descriptor table observable. ex01 starts with stdin, stdout, and
stderr, then follows `open` and `close` as the kernel assigns and reuses the lowest free
number; ex02 uses those handles to compare byte-transfer APIs. Write `main.c` so `make all`
produces `./test`, and read every reported number from the system rather than hard-coding it.

## Shape

This is a whole-program exercise: write `main.c` and its `Makefile`; `make all` produces
`./test`. Your program starts with three handles open before it does anything else: standard
input, standard output, and standard error. Every `open` buys the lowest free number, and
every `close` gives it back.

- Print `fileno(stdin)`, `fileno(stdout)`, and `fileno(stderr)`.
- Call `open("/dev/null", O_RDONLY)` once, print `first open <fd>`, and close it.
- Open `/dev/null` again and print `reuse after close <fd>`.
- Open twice in a row and print the descriptors in ascending order as `two opens <a> and <b>`.
- Close both, open once more, and print `lowest available <fd>`.
- Check every `open`/`close` result and close every fd you own.

Descriptor numbers are exact and deterministic (`0`, `1`, `2`, then `3`, `4`, …); never use
`printf("3")` strings, because the value must come from `open`. Compile with
`-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan, without redefining `CFLAGS`/`LDFLAGS`.

## Acceptance

`make all`, then `./test` must exit 0 and print exactly:

```text
stdin 0
stdout 1
stderr 2
first open 3
reuse after close 3
two opens 3 and 4
lowest available 3
ALL PASS
```

- The first three values come from `fileno`, the open values come from the actual descriptor results, and the simultaneous opens are printed ascending.
- `first open 3`, `reuse after close 3`, `two opens 3 and 4`, and `lowest available 3` demonstrate allocation, release, reuse, and lowest-available selection.
- Exit is 0, stderr is empty, every fd is closed before exit, and `quiz.txt` is complete (see below).

The `build` grade is a strict compile (`-std=gnu11 -Wall -Wextra -Werror` + ASan/UBSan) plus the `./test` run, stdout diff, exit 0, and empty stderr. Hard-coding the numbers can make the diff look right, but the descriptor contract and the quiz still require the values to come from the system; a reference that misorders the two simultaneous opens fails with a wrong `two opens` line. The `quiz` grade checks that `quiz.txt` answers match.

## Readings

- `man 2 open`, `man 2 close` — the "lowest-numbered file descriptor not currently open"
  contract, straight from the manual.
- TLPI §5.4 "Relationship Between File Descriptors and Open Files" — the table of
  numbers, and why storing them in shell redirections (`>`, `<`) is just `open` +
  table surgery done for you.
- The Linux Programming Interface, Chapter 4 "File I/O: The Universal I/O Model".

## Quiz

1. What is a file descriptor number?
2. Which integer does the first open() usually return?