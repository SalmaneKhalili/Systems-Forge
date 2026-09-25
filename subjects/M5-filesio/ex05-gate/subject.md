# M5-ex05 · Gate: Log Parser

The gate applies the descriptor model to structured input. ex05 reads a `name=value` file
through a descriptor, assembles lines across arbitrary `read()` boundaries, and returns only
names in order; the harness checks the parser without giving up ownership of the fd. You
deliver `log.c` against the provided `fix.h` and `main.c`.

## Shape

Implement `log.c` behind the provided `fix.h` + `main.c`: a tiny **key=value log parser**
that reads from a file descriptor. `forge` compiles `main.c` (harness) + `log.c`, runs
`./test config.txt` and `./test apps.txt` (fixtures provided), and diffs stdout.

```c
int parse_log(int fd, Tokens *t);  /* 0 ok, -1 on malformed or overflow */
```

The file is a list of `name=value` lines. Record only the **name**, up to the `=`. Skip
blank lines (just `\n`). A non-blank line with no `=`, an empty name, or a name that does
not fit `TOKEN_LEN` returns **-1**. The whole file is read through the given `fd` with
`read()` — nothing else. The `Tokens` struct is zero-sized until you fill it:

```c
typedef struct { char keys[TOKENS_MAX][TOKEN_LEN]; size_t n; } Tokens;
```

`read()` may return the file in small pieces, so assemble lines across reads (buffering) and
handle a final line with or without a trailing newline. A name is complete only at `\0`, with
no stray newline or value bytes in `keys`. The `fd` is opened and closed by `main`; your
function must not close it.

Compile with `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan, without redefining
`CFLAGS`/`LDFLAGS`.

## Acceptance

`make all`, then `./test config.txt` and `./test apps.txt` must each exit 0 and print exactly:

```text
host
port
verbose
tokens 3
ALL PASS
```

```text
web
cache
queue
tokens 3
ALL PASS
```

- The `config.txt` transcript is `host`, `port`, `verbose`, `tokens 3`; the `apps.txt` transcript is `web`, `cache`, `queue`, `tokens 3`, with names, order, and count exact.
- Both fixture runs exit 0 with empty stderr; malformed input is refused with `-1`, never tolerated.
- The final unterminated line is handled, the fd is left to `main`, and `quiz.txt` is complete (see below).

The `build` grade is a strict compile plus both fixture runs, each stdout diffed, with exit 0
and empty stderr. Splitting on the wrong delimiter (`:`) makes every line malformed →
`parse failed` → FAIL. Storing the whole `name=value` instead of the name makes `keys` hold
the full line → diff FAIL. Never finishing the final unterminated line leaves output one
token short, or leaves the names with their newline → FAIL. The `quiz` grade checks that
`quiz.txt` answers match.

## Readings

- `man 2 read`, `man 2 open` — the fd contract, again: partial reads and EOF are normal.
- `man 3 strchr`, `man 3 memchr` — line/segment scanning without string-overflow traps.
- TLPI Chapter 4 "File I/O: The Universal I/O Model" — everything in this gate is the
  "universal I/O model" applied to structured text.
- A real `strtok`-style parser is the wrong tool here: you need off-by-one-exact
  delimiters, and `TOKEN_LEN` overflow handling.

## Quiz

1. What character separates a name from a value in a key=value line?
2. What should a parsing function return for a malformed line?