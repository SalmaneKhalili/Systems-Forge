# M6-ex03 · HTTP-ish server

## Goal

Write `main.c` so `make all` produces `./test` — an **HTTP/1.1-ish server**: reads a full
request (blank line `"\r\n\r\n"` ends it), responds to

```
GET /ping HTTP/1.1
Host: loopback
```

with the exact response

```
HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\npong\n
```

and to any other path/method with

```
HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n
```

The graded dialogue (single run):

```
send "GET /ping HTTP/1.1\r\nHost: loopback\r\n\r\n"
  → expect "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\npong\n"
send "GET /nope HTTP/1.1\r\nHost: loopback\r\n\r\n"
  → expect "HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"
```

## Constraints

- Parse until you find `"\r\n\r\n"` in the receive buffer; if the buffer fills before that,
  reply 400 or just close — anything is legal (the grader never sends a request that long).
- Compare only the method and path; `Host` and other headers are ignored but must not crash
  the parser.
- Reply is exactly the raw bytes of the HTTP response, including `\r\n` line endings.
  - 200: the body is `pong\n`, length 5
  - 404: empty body, `Content-Length: 0`
- One request per connection is fine (HTTP/0.9 style is acceptable to this grader).
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan.

## Acceptance criteria

- [ ] GET /ping → 200 + body; GET /anything → 404 — exact bytes
- [ ] server re-accepts after EOF (loop survives the first disconnect)
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which HTTP request method does this exercise serve?: <answer>
Which number means the resource was not found?: <answer>
```

## Readings

- `man 2 read`, `man 2 write`; `memchr`/`strstr` to find `"\r\n\r\n"`.
- RFC 9110 §4 / §6 — the semantics of GET, 200, 404, and the role of Content-Length.
  TLPI §59.2 (the "scatter/gather" read loop from the client perspective) is also relevant.

## How you are graded

- `build` strict; `net` sends the raw HTTP request bytes, reads raw response bytes with a
  deadline. No `\r\n` in the 200 response: response length won't match Content-Length:5 →
  FAIL. 404 wrong number or missing body: byte mismatch → FAIL.
- `quiz`: `quiz.txt` answers must match.