# M5-ex05 gate · log parser

## Goal

Implement `log.c` backing the provided `fix.h` + `main.c`: a tiny **key=value log
parser** that reads from a file descriptor. `forge` compiles `main.c` (harness) + `log.c`,
runs `./test config.txt` and `./test apps.txt` (fixtures provided), diffs stdout.

```c
int parse_log(int fd, Tokens *t);  /* 0 ok, -1 on malformed or overflow */
```

- The file is a list of `name=value` lines. Record only the **name** (up to the `=`).
- Blank lines (just `\n`) are skipped.
- A non-blank line with no `=`, an empty name, or a name that doesn't fit `TOKEN_LEN`
  returns **-1**.
- The whole file is read through the given `fd` with `read()` — nothing else. The
  `Tokens` struct is zero-sized until you fill it:

```c
typedef struct { char keys[TOKENS_MAX][TOKEN_LEN]; size_t n; } Tokens;
```

Expected for `config.txt`:

```
host
port
verbose
tokens 3
ALL PASS
```

Expected for `apps.txt`:

```
web
cache
queue
tokens 3
ALL PASS
```

Everything the harness prints comes from your `keys` — so the token names, their order,
and `tokens <n>` must all be exactly right.

## Constraints

- `read()` may return the file in small pieces — your parser must assemble lines across
  reads (buffering!), and handle a final line with or without a trailing newline.
- A name is complete only at `\0` — no stray newline or value bytes in `keys`.
- The `fd` is opened and closed by `main`; your function must not close it.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan. Never redefine `CFLAGS`/`LDFLAGS`.

## Acceptance criteria

- [ ] `config.txt` → `host`, `port`, `verbose`, `tokens 3`
- [ ] `apps.txt` → `web`, `cache`, `queue`, `tokens 3`
- [ ] both runs exit 0, empty stderr
- [ ] malformed input is refused (`-1`), never tolerated
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
What character separates a name from a value in a key=value line?: <answer>
What should a parsing function return for a malformed line?: <answer>
```

## Readings

- `man 2 read`, `man 2 open` — the fd contract, again: partial reads and EOF are normal.
- `man 3 strchr`, `man 3 memchr` — line/segment scanning without string-overflow traps.
- TLPI Chapter 5 — everything in this gate is the "universal I/O model" applied to
  structured text.
- A real `strtok`-style parser is the wrong tool here: you need off-by-one-exact
  delimiters, and `TOKEN_LEN` overflow handling.

## How you are graded

- `build`: strict compile + both fixture runs, each stdout diffed, exit 0, empty stderr.
  Split on the wrong delimiter (`:`): every line is malformed → `parse failed` → FAIL.
  Store the whole `name=value` instead of the name: `keys` hold the full line → diff FAIL.
  Never finish the final unterminated line: output ends up one token short (or the names
  keep their newline) → FAIL.
- `quiz`: `quiz.txt` answers must match.