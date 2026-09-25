# M6-ex03 · HTTP-Ish Server

Line chat invented one framing rule; HTTP is a framing rule everyone knows.
This exercise takes you to the other end of a request: read a full HTTP
request (the blank line `"\r\n\r\n"` ends it), answer `GET /ping` with a
`200` and its body, and everything else with a `404` — in exact raw HTTP
bytes, `\r\n` line endings included.

## Shape

Write `main.c` so `make all` produces `./test` — an HTTP/1.1-ish server.

- Parse until you find `"\r\n\r\n"` in the receive buffer; if the buffer
  fills before that, reply `400` or just close — anything is legal (the
  grader never sends a request that long).
- Compare only the method and path; `Host` and other headers are ignored but
  must not crash the parser.
- Reply exactly the raw bytes of the HTTP response, `\r\n` line endings
  included. A request

```text
GET /ping HTTP/1.1
Host: loopback
```

is answered with

```text
HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\npong\n
```

and any other path or method with

```text
HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n
```
- One request per connection is fine (HTTP/0.9 style is acceptable to this
  grader).

Compile flags: `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan; never redefine
`CFLAGS`/`LDFLAGS`.

## Acceptance

Graded `build` + `net` + `quiz`. The single-run dialogue:

```
send "GET /ping HTTP/1.1\r\nHost: loopback\r\n\r\n"
  → expect "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\npong\n"
send "GET /nope HTTP/1.1\r\nHost: loopback\r\n\r\n"
  → expect "HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"
```

The `net` runner sends the raw bytes above and reads the raw response with a
deadline, comparing byte-for-byte:

- GET /ping → 200 + body; GET /anything → 404 — exact bytes.
- Server re-accepts after EOF (the loop survives the first disconnect).
- No `\r\n` in the 200 response → response length won't match
  `Content-Length: 5` → FAIL.
- 404 with the wrong number or a missing body → byte mismatch → FAIL.

`quiz.txt` is complete (see Quiz).

## Readings

- `man 2 read`, `man 2 write`; `memchr`/`strstr` to find `"\r\n\r\n"`.
- RFC 9110 — the semantics of GET (§9.3.1), 200 (§15.3.1), 404 (§15.5.5), the
  Content-Length header (§8.6), and message framing (§6.1).
- TLPI §56.5.4 "I/O on Stream Sockets" (the read-loop, framing, and partial-read
  concerns from the client perspective) is also relevant.

## Quiz

1. Which HTTP request method does this exercise serve?
2. Which number means the resource was not found?