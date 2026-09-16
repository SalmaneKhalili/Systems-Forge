# M6-ex02 · line chat

## Goal

Write `main.c` so `make all` produces `./test` — a **line-oriented chat server**: it reads
bytes, splits them on `\n` (a framing rule), and for every complete line replies
`<lineno>:<the line>\n`, where `<lineno>` counts complete lines **within the current
connection** (1, 2, 3, …). The accept loop keeps running after a client disconnects —
a fresh connection restarts the counter at 1.

Graded dialogues (two separate server runs):

```
send "hello\n"  → expect "1:hello\n"
send "world\n"  → expect "2:world\n"
send "again\n"  → expect "3:again\n"
```
then a *second* run (proves you accept a new connection and reset state):
```
send "hi\n"     → expect "1:hi\n"
```

## Constraints

- The line is the unit: the client may split a message across several `read`s (or pack
  several lines in one), and your server must reassemble lines from a buffer either way.
- Reply prefix and line contents must match exactly — including the counter going back to
  `1:` on the second connection.
- `read` returns fewer/multiple lines freely; `write` the reply as its exact bytes
  (prefix + line text + `\n`).
- The accept loop is the point: one `accept`, one `serve`, forever. Losing the loop loses
  the second graded run.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan. Never redefine `CFLAGS`/`LDFLAGS`.

## Acceptance criteria

- [ ] `1:hello\n` / `2:world\n` / `3:again\n` on the first connection — exact
- [ ] second server run answers `1:hi\n` — the counter is per-connection
- [ ] server survives EOF (client disconnect) and continues accepting
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
What does the accept loop let a server do?: <answer>
How does a server usually know where one message ends?: <answer>
```

## Readings

- `man 2 accept`, `man 2 read`, `man 2 write`, `man 2 memmove`.
- TLPI §56.4–56.5 — the accept loop pattern you are extending; "stream sockets have no
  message boundaries" is *precisely* why a framing rule exists.
- TLPI §56.5.4 "I/O on Stream Sockets" and §61.1 "Partial Reads and Writes on Stream
  Sockets" if you want to see the same framing problem from the client side
  (readall-style loops).

## How you are graded

- `build` strict; two `net` sessions drive the exact dialogues above, byte-compared with a
  bounded deadline. Accept once and never loop: the second run finds a dead server → FAIL.
  No framing (echo raw bytes): reply never equals `1:hello\n` → FAIL. Reply counter not
  reset per connection: second run gets the wrong number → FAIL.
- `quiz`: `quiz.txt` answers must match.