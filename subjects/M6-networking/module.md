# M6 · Networking

Sockets are file descriptors with a destination. The whole dance is five calls
(`socket`, `bind`, `listen`, `accept`, and `read`/`write` on what `accept` hands you), and
the moment you see it, servers stop being magic: a server is a program that waits for a
connection, then talks over an fd it never had to `open`.

The thread through this module: a bare **echo** server, a **line-chat** server with
per-connection state, an **HTTP-ish** server that parses request lines, a server that
**times out** half-open clients instead of hanging forever — and a stateful **gate** server
that greets, acknowledges, and re-answers across several connections, the skeleton of the
switch (module M8) the rest of the curriculum grades through.

| exercise | kind   | what you build                                              |
|----------|--------|-------------------------------------------------------------|
| ex01     | echo   | single-connection echo server                                |
| ex02     | linechat | multi-connection line server with per-connection counters  |
| ex03     | http   | HTTP/1.1-ish request/response (200 and 404)                 |
| ex04     | timeout| `SO_RCVTIMEO`, a server that never lets a half-client wedge it |
| ex05     | **gate** | stateful protocol server with greeting + commands         |

Every exercise ships a `main.c` server you write and a `Makefile`. Grading is
`build` + `net` + `quiz`: your binary is compiled, **spawned**, and driven by a scripted
client — the grader does a fresh TCP connection, plays the exact dialogue in
`exercise.json`, and compares what it gets, byte for byte. No sleeps, no "wait for it to
finish": the network run has a built-in connect retry, so timing cannot make a correct
server flaky.

Two rules make this deterministic, and they are worth internalizing now:

- Your server binds **`127.0.0.1` on the port in the `TARGETPORT` environment variable**
  (the grader injects it — how real tooling hands servers their ports).
- Nothing your server prints to stdout or stderr is graded (the grader closes those), so
  the CPU you spend on "correct" output is better spent on exact **socket** semantics.

---

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

---

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