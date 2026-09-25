# M6-ex01 · Echo Server

Sockets are file descriptors with a destination, and this exercise builds the
whole dance end to end. Your TCP echo server binds `127.0.0.1` on the
injected `TARGETPORT` and sends back whatever bytes a client sends — the same
copy loop that sits at the heart of almost every networked program. Every
later exercise in this module reuses this skeleton.

## Shape

Write `main.c` so `make all` produces `./test`. The five-call sequence, in
order:

```
socket(SOCK_STREAM) → setsockopt(SO_REUSEADDR) → bind → listen → accept → read → write
```

- `socket(AF_INET, SOCK_STREAM, 0)`, then `bind`, then `listen(fd, 4)`, then
  `accept` in a loop. `SO_REUSEADDR` before bind (so restarting is painless).
- The port comes from `getenv("TARGETPORT")`; reject a missing or invalid
  value with exit 1.
- Loop forever: `read` whatever arrives, `write` the same bytes back, until
  the client closes; then accept the next connection and repeat.
- `read` may return fewer bytes than the client "sent" — echo exactly the
  bytes you got (`write` must write `n`, not an assumed buffer size).
- A client that disconnects must not kill the server: `read` returns 0 (EOF) →
  close and accept again.
- `signal(SIGPIPE, SIG_IGN)` so a write to a peer that just went away fails
  gracefully instead of killing the process.

Compile flags: `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan; never redefine
`CFLAGS`/`LDFLAGS`.

## Acceptance

Graded `build` + `net` + `quiz`. `make all` produces `./test`, and the `net`
runner dials your server and plays two scripted sessions, comparing replies
byte-for-byte with a bounded deadline:

- client sends `hello\n` → server sends back exactly `hello\n`
- client sends `ping\n` → server sends back exactly `ping\n`
- the accept loop runs forever, re-serving after EOF
- binds loopback only, on the injected port

A server that replies from memory (e.g. always `PONG\n`), echoes the wrong
slice, or never listens → FAIL. `quiz.txt` is complete (see Quiz).

## Readings

- `man 2 socket`, `man 2 bind`, `man 2 listen`, `man 2 accept`, `man 2 read`, `man 2 write`.
- TLPI §56.4–56.5 "sockets: stream sockets" — the socket/bind/listen/accept skeleton and
  the accept loop; §56.2 for `socket()` itself.
- The echo server is the Unix miniature of every networked program: a listener plus a
  per-connection copy loop.

## Quiz

1. Which call puts a listening socket's address family, port and address on it?
2. Which call accepts a pending connection and returns a new socket for it?