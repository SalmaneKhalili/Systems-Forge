# M6-ex02 · Line Chat

The echo server copied raw bytes; a real protocol needs a unit of meaning —
and a stream has none built in. This exercise adds the first framing rule: you
read bytes, split on `\n`, and answer every complete line with
`<lineno>:<line>`, where the number counts complete lines **within the current
connection** (1, 2, 3, …). The accept loop keeps running after a client
disconnects, and a fresh connection restarts the counter at 1.

## Shape

Write `main.c` so `make all` produces `./test` — a line-oriented chat server.

- The line is the unit: the client may split a message across several `read`s
  (or pack several lines into one), so your server must reassemble lines from
  a buffer either way.
- Reply prefix and line contents must match exactly — including the counter
  going back to `1:` on the second connection.
- `read` may return fewer or multiple lines freely; `write` the reply as its
  exact bytes (prefix + line text + `\n`).
- The accept loop is the point: one `accept`, one `serve`, forever. Losing
  the loop loses the second graded run.

Compile flags: `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan; never redefine
`CFLAGS`/`LDFLAGS`.

## Acceptance

Graded `build` + `net` + `quiz`. Two separate server runs drive the exact
dialogues, byte-compared with a bounded deadline:

```
send "hello\n"  → expect "1:hello\n"
send "world\n"  → expect "2:world\n"
send "again\n"  → expect "3:again\n"
```

then a *second* run (proves you accept a new connection and reset state):

```
send "hi\n"     → expect "1:hi\n"
```

- Accept once and never loop: the second run finds a dead server → FAIL.
- No framing (echo raw bytes): the reply never equals `1:hello\n` → FAIL.
- Reply counter not reset per connection: the second run gets the wrong
  number → FAIL.

`quiz.txt` is complete (see Quiz).

## Readings

- `man 2 accept`, `man 2 read`, `man 2 write`, `man 2 memmove`.
- TLPI §56.4–56.5 — the accept loop pattern you are extending; "stream sockets have no
  message boundaries" is *precisely* why a framing rule exists.
- TLPI §56.5.4 "I/O on Stream Sockets" and §61.1 "Partial Reads and Writes on Stream
  Sockets" if you want to see the same framing problem from the client side
  (readall-style loops).

## Quiz

1. What does the accept loop let a server do?
2. How does a server usually know where one message ends?