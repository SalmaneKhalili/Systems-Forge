# M6-ex05 · Gate: Stateful Server

The gate of M6 is the skeleton of the infrastructure the rest of the
curriculum grades *through*: a **stateful protocol server** whose connection
counter survives across connections, and that greets, acknowledges, and
re-answers. You write it now so you can point to it in M8, where the switch
takes the same shape and adds faults.

## Shape

Write `main.c` so `make all` produces `./test`. The protocol, on each
connection:

1. On accept, send the greeting `HELLO <n>\n`, where `<n>` counts connections
   made to this server process (1 for the first, 2 for the second, … never
   reset).
2. Then loop over complete lines:
   - `PING\n` → reply `PONG\n`
   - `ECHO <text>\n` → reply `ECHO <text>\n` (only the text, no extra spaces)
   - anything else → reply `BAD\n`
3. On EOF (client closed), close the connection and `accept` the next one.

- The greeting **must be sent immediately on accept**, not on demand; state is
  shared across connections (the counter lives in `main`, or in a static).
- Line framing, echo reply, `BAD` fallback, counter — all exact.
- On a disconnect mid-line: discard the partial line, close, accept again.

Compile flags: `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan; never redefine
`CFLAGS`/`LDFLAGS`.

## Acceptance

Graded `build` + `net` + `quiz`. Two `net` runs — the second proves state
survives the new connection:

```
run 1:  (send "")   → expect "HELLO 1\n"
        send "PING\n"        → expect "PONG\n"
        send "ECHO awesome\n" → expect "ECHO awesome\n"
run 2:  (send "")   → expect "HELLO 2\n"
        send "PING\n"        → expect "PONG\n"
```

The `""` send means *just connect*: your server must send the greeting on its
own when a connection is accepted, without waiting for a request line.

- Greeting on demand, never on accept → the `""` step times out, FAIL.
- No stateful counter → the second run sees `HELLO 1`, FAIL.
- Wrong PING/ECHO reply → byte mismatch, FAIL.
- The server survives both disconnects (never exits on EOF).

`quiz.txt` is complete (see Quiz).

## Readings

- `man 2 getenv`, `man 2 getsockopt` — the injected `TARGETPORT`, `SO_REUSEADDR`.
- TLPI §56.5 — the accept loop; the greeting-on-accept is the same pattern as any
  protocol server (SMTP banner, FTP 220, SSH banner…).
- This is the server the rest of the curriculum grades *through* — make it robust.

## Quiz

1. Which environment variable carries the port the grader injects?
2. Which socket option lets a server rebind an address still in TIME_WAIT?