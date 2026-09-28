# M6-ex05 · Gate: stateful server

## Goal

Write `main.c` so `make all` produces `./test` — a **stateful protocol server** with a
connection counter that survives across connections. This is the skeleton of the grading
**gate** server that switches later modules' exercises (see PLAN.md) — the server that
greets, acknowledges, and cross-connects state.

Protocol, on each connection:

1. On accept, send the greeting `HELLO <n>\n`, where `<n>` counts connections made to
   this server process (1 for the first, 2 for the second, … never reset).
2. Then loop over complete lines:
   - `PING\n` → reply `PONG\n`
   - `ECHO <text>\n` → reply `ECHO <text>\n` (only the text, no extra spaces)
   - anything else → reply `BAD\n`
3. On EOF (client closed), close the connection and `accept` the next one.

Graded dialogue — two `net` runs (the second proves state survives the new connection):

```
run 1:  (send "")   → expect "HELLO 1\n"
        send "PING\n"        → expect "PONG\n"
        send "ECHO awesome\n" → expect "ECHO awesome\n"
run 2:  (send "")   → expect "HELLO 2\n"
        send "PING\n"        → expect "PONG\n"
```

The `""` send means *just connect*: your server must send the greeting on its own when a
connection is accepted, without waiting for a request line.

## Constraints

- Greeting **must be sent immediately on accept**, state shared across connections
  (the counter lives in `main`, or in a static).
- Line framing, echo reply, `BAD` fallback, counter — all must be exact.
- On a disconnect mid-line: discard the partial line, close, accept again.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan.

## Acceptance criteria

- [ ] `HELLO 1\n` and `HELLO 2\n` in two runs — greeting cross-connection counter
- [ ] `PONG\n` and `ECHO awesome\n` exact
- [ ] server survives both disconnects (never exits on EOF)
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which environment variable carries the port the grader injects?: <answer>
Which socket option lets a server rebind an address still in TIME_WAIT?: <answer>
```

## Readings

- `man 2 getenv`, `man 2 getsockopt` — the injected `TARGETPORT`, `SO_REUSEADDR`.
- TLPI §56.5 — the accept loop; the greeting-on-accept is the same pattern as any
  protocol server (SMTP banner, FTP 220, SSH banner…).
- This is the server the rest of the curriculum grades *through* — make it robust.

## How you are graded

- `build` strict; two `net` sessions. Greeting on demand, never on accept → `""` step
  times out, FAIL. No stateful counter → second run sees `HELLO 1`, FAIL. Wrong PING/ECHO
  reply → byte mismatch, FAIL.
- `quiz`: `quiz.txt` answers must match.