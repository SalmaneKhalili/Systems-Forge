# M6-ex01 · echo server

## Goal

Write `main.c` so `make all` produces `./test` — a **TCP echo server** on `127.0.0.1`:
whatever bytes a client sends, the server sends back. Two scripted visits are graded:

- client sends `hello\n` → server must send back exactly `hello\n`
- client sends `ping\n` → server must send back exactly `ping\n`

```
socket(SOCK_STREAM) → setsockopt(SO_REUSEADDR) → bind → listen → accept → read → write
```

Your server binds **`127.0.0.1`, port from the `TARGETPORT` env var** (the grader injects
it). It then loops: `read` whatever arrives, `write` the same bytes back, until the client
closes; then accept the next connection and repeat — forever.

## Constraints

- The five-call sequence, in order: `socket(AF_INET, SOCK_STREAM, 0)`, `bind`,
  `listen(fd, 4)`, then `accept` in a loop. `SO_REUSEADDR` before bind (so restarting is
  painless). Port from `getenv("TARGETPORT")`; reject a missing/invalid value with exit 1.
- `read` may return fewer bytes than the client "sent" — echo exactly the bytes you got
  (`write` must write `n`, not an assumed buffer size).
- A client that disconnects must not kill the server: `read` returns 0 (EOF) → close and
  accept again.
- `signal(SIGPIPE, SIG_IGN)` so a write to a peer that just went away fails gracefully
  instead of killing the process.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan. Never redefine `CFLAGS`/`LDFLAGS`.

## Acceptance criteria

- [ ] `make all` produces `./test`, then `accept` loop runs forever, re-serving after EOF
- [ ] graded dialogue `hello\n`↔`hello\n` and `ping\n`↔`ping\n` passes byte-exact
- [ ] binds loopback only, on the injected port
- [ ] `quiz.txt` complete (see below)

Then complete `quiz.txt`:

```
Which call puts a listening socket's address family, port and address on it?: <answer>
Which call accepts a pending connection and returns a new socket for it?: <answer>
```

## Readings

- `man 2 socket`, `man 2 bind`, `man 2 listen`, `man 2 accept`, `man 2 read`, `man 2 write`.
- TLPI §56.4–56.5 "sockets: stream sockets" — the socket/bind/listen/accept skeleton and
  the accept loop; §56.2 for `socket()` itself.
- The echo server is the Unix miniature of every networked program: a listener plus a
  per-connection copy loop.

## How you are graded

- `build` compiles strict; two `net` scripted sessions then dial your server, scan for the
  expected reply with a bounded deadline, and compare byte-for-byte. A server that replies
  from memory (e.g. always `PONG\n`), echoes the wrong slice, or never listens → FAIL.
- `quiz`: `quiz.txt` answers must match.