# M6 · Networking

Sockets are file descriptors with a destination. The whole dance is five calls
— `socket`, `bind`, `listen`, `accept`, and `read`/`write` on what `accept`
hands you — and the moment you see it, servers stop being magic: a server is a
program that waits for a connection, then talks over an fd it never had to
`open`. This module walks that thread from a bare echo server up to a stateful
gate that greets, acknowledges, and re-answers across several connections —
the skeleton of the fault-injection switch (M8) that the rest of the
curriculum grades through.

## The build

- **ex01 · Echo server** — a single-connection echo server: `socket`→`bind`→
  `listen`→`accept`, then an identical-bytes copy loop.
- **ex02 · Line chat** — the same accept loop with per-connection state: lines
  are numbered within their connection, and a fresh connection restarts the
  counter.
- **ex03 · HTTP-ish server** — parse a request until the `"\r\n\r\n"` blank
  line and answer `200`/`404` in exact raw HTTP bytes.
- **ex04 · Read timeout** — `SO_RCVTIMEO` so a half-open client cannot wedge
  the server: no data for 400 ms means `TIMEOUT\n` and re-accept.
- **ex05 · Gate: stateful server** — a greeting counter that survives
  connections, PING/ECHO/BAD handling, and clean teardown: the skeleton of the
  M8 switch.

## Rules

- Every exercise: you write `main.c` (a full server) plus a `Makefile`. Grading
  is `build` + `net` + `quiz`: the binary is compiled, **spawned**, and driven
  by a scripted client — the grader makes a fresh TCP connection, plays the
  exact dialogue, and compares byte for byte. The `net` run has a built-in
  connect retry, so timing cannot make a correct server flaky.
- Your server binds **`127.0.0.1` on the port in the `TARGETPORT` environment
  variable** (the grader injects it — how real tooling hands servers their
  ports). A missing or invalid `TARGETPORT` is a failure path (exit 1).
- Nothing your server prints to stdout or stderr is graded (the grader closes
  those), so spend your effort on exact **socket** semantics.
- Compile flags: `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan; never
  redefine `CFLAGS`/`LDFLAGS`.

## Prerequisites

Before starting M6, you should be comfortable with everything from M0–M5, plus:

- Explain what a socket is: a file descriptor with a remote address attached.
- Call `socket(AF_INET, SOCK_STREAM, 0)` and check the return value.
- Explain `struct sockaddr_in`: `sin_family`, `sin_port` (network byte order), `sin_addr`.
- Explain `htons`/`htonl`/`ntohs`/`ntohl`: convert between host byte order and network
  (big-endian) byte order.
- Call `bind(sockfd, &addr, sizeof(addr))` and explain why `INADDR_ANY` means "listen on
  all interfaces."
- Call `listen(sockfd, backlog)` — explain that it marks the socket as passive (server side).
- Call `accept(sockfd, ...)` — returns a NEW file descriptor for the connected client.
- Call `read`/`write` on the accepted fd — a TCP connection is just a bidirectional byte stream.
- Explain `SO_REUSEADDR`: why you need it to rebind a port after a crash.
- Handle `SIGPIPE` by ignoring it (`signal(SIGPIPE, SIG_IGN)`) — explain that writing to a
  broken connection would otherwise kill your process.
- Call `getenv("TARGETPORT")` and convert it with `atoi()` or `strtol()`.
- Handle timeouts: explain that `setsockopt` + `SO_RCVTIMEO` makes `read` return -1 with
  `errno == EAGAIN` if no data arrives in time.

You do NOT need to know: non-blocking I/O, `select`/`poll`, UDP, or DNS resolution.
You will learn those concepts later.

## So what? (interview / portfolio)

"Sockets are fds with a destination" is the lens that makes every server — web server, RPC
service, your own gateway — stop being magic. This is also the module where the grader
starts *driving your binary over real TCP*, the same contract production tooling uses, so
you leave with a runnable, scriptable server under your belt.

**Interview questions this module arms you for:**
- Walk through the five calls that make a server: `socket`→`bind`→`listen`→`accept`→`read/write`.
- Why bind to `127.0.0.1:<port>` from an env var rather than a hard-coded port?
- How do you stop a half-open connection from wedging a server forever (`SO_RCVTIMEO`)?
- How does an HTTP request start (the request line)? How would you parse it minimally?

**Portfolio artifact:** M6-ex05 `gate` — a stateful protocol server that greets,
acknowledges, and re-answers across several connections (the skeleton of the M8 switch).