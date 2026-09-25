# M7-ex03 · Health Checks

ex01 and ex02 built the client's side of resilience; a **health-check endpoint**
is the other half — what a real load balancer or supervisor probes before
routing traffic. Your Go TCP server answers `GET /health` with live or dead,
and can be flipped over the wire with explicit routes, keeping its state
across connections. It is the smallest possible version of the liveness probe
that your future supervisor (ex05) learns how to act on.

## Shape

Write a Go program, whole-program. Whoever checks you out over the wire will
drive this exact dialogue (each line is a **fresh connection**, one request
each):

```text
GET /health → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"
GET /down   → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"   (now unhealthy)
GET /health → HTTP/1.1 503 Service Unavailable  Content-Length: 0  (dead)
GET /up     → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"   (healthy again)
GET /health → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"
```

(Exact response bytes, `\r\n` line endings, are in `exercise.json`.) Rules:

- Bind `127.0.0.1` on the port in the `TARGETPORT` environment variable;
  refuse to run without it (exit 1, one stderr line).
- One **request line** is enough: read the first `\r\n`-terminated line and
  respond to `GET /health`, `GET /down`, `GET /up`; anything else →
  `404 Not Found`.
- **State lives across connections**: `/down` means the *next* `/health` from
  any client answers 503, until `/up`. A server that forgets its state fails
  the dialogue.
- Accept and serve forever; a disconnect mid-request must not kill the loop.
- No output on stdout/stderr during normal operation.

## Acceptance

Graded `build` + `net` + `quiz`: `go build -o test .` must succeed, the `net`
grader spawns your server and replays the five connections byte-exactly with
a bounded deadline, and the diff must be clean:

- The five-step dialogue passes byte-exact (200/503/200 with exact bodies).
- `/down` state persists until `/up`, across connections. No state flip → the
  third `/health` answers 200 instead of 503 → FAIL.
- Unknown path → 404.
- The server keeps accepting after each disconnect.

`quiz.txt` is complete (see Quiz).

## Readings

- RFC 9110 — the HTTP statuses the endpoint speaks: 200 (§15.3.1), 404 (§15.5.5),
  503 (§15.6.4) semantics.
- Kubernetes *Configure Probes* — liveness vs readiness: the operational vocabulary this
  exercise is the smallest example of.
- TLPI §56 — the accept-loop skeleton you reuse from M6.

## Quiz

1. Which routing path reports whether the service is currently healthy?
2. Which HTTP status does an unhealthy health endpoint return?
3. Which path flips the service unhealthy for the next probe?