# M7-ex03 · Health checks

## Goal

Write a **health-check endpoint** in Go — a TCP server whose `GET /health` route reports
live or dead, and that can be flipped with explicit routes. This is what real load
balancers and supervisors probe before routing traffic.

Whoever checks you out over the wire will drive this exact dialogue (each line is a **fresh
connection**, one request each):

```
GET /health → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"
GET /down   → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"   (now unhealthy)
GET /health → HTTP/1.1 503 Service Unavailable  Content-Length: 0  (dead)
GET /up     → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"   (healthy again)
GET /health → HTTP/1.1 200 OK  Content-Length: 3  body "ok\n"
```

(Exact response bytes, `\r\n` line endings, are in `exercise.json`.)

## Constraints

- Bind `127.0.0.1` on the port in the `TARGETPORT` environment variable; refuse to run
  without it (exit 1, one stderr line).
- One **request line** is enough: read the first `\r\n`-terminated line and respond to
  `GET /health`, `GET /down`, `GET /up`; anything else → `404 Not Found`.
- **State lives across connections**: `/down` means the *next* `/health` from any client
  answers 503, until `/up`. A server that forgets its state fails the dialogue.
- Accept and serve forever; a disconnect mid-request must not kill the loop.
- No output on stdout/stderr during normal operation.

## Acceptance criteria

- [ ] the five-step dialogue passes byte-exact (200/503/200 exact bodies)
- [ ] `/down` state persists until `/up`, across connections
- [ ] unknown path → 404
- [ ] server keeps accepting after each disconnect
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which routing path reports whether the service is currently healthy?: <answer>
Which HTTP status does an unhealthy health endpoint return?: <answer>
Which path flips the service unhealthy for the next probe?: <answer>
```

## Readings

- RFC 9110 — the HTTP statuses the endpoint speaks: 200 (§15.3.1), 404 (§15.5.5),
  503 (§15.6.4) semantics.
- Kubernetes *Configure Probes* — liveness vs readiness: the operational vocabulary this
  exercise is the smallest example of.
- TLPI §56 — the accept-loop skeleton you reuse from M6.

## How you are graded

- `build` compiles (`go build -o test .`); the `net` grader spawns your server and replays
  the five connections byte-exactly with a bounded deadline; `quiz` covers the concepts.
  No state flip → the third `/health` answers 200 instead of 503 → FAIL.