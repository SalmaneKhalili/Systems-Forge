# M6-ex04 · Read Timeout

The last three servers trusted the client to be well-behaved. This exercise
covers the operational half: a **line-echo server with a read timeout** —
every connection that has no data for **400 ms** is terminated with
`TIMEOUT\n`. A stuck or half-open client can never wedge your server.

## Shape

Write `main.c` so `make all` produces `./test`. The loop:

```
socket → bind → listen → loop:
  accept
  SO_RCVTIMEO = 400ms on accepted fd (or inherited from listening fd)
  while (read > 0) {
    reassemble line from buffer
    on complete line → write exact line back
  }
  if read timed out (EAGAIN/EWOULDBLOCK) → write "TIMEOUT\n"
  close
```

- **`SO_RCVTIMEO`** (or `poll`/`select` with a timeout — both acceptable) on
  the read side. The exact mechanism is yours; the point is that a stuck
  client cannot wedge the server.
- Timeout = **400 ms**. Your server must fail before the grader's 1200 ms
  delay ends.
- A complete line echoed with the correct newline; a timed-out connection must
  produce exactly `TIMEOUT\n`.
- After the timeout event: close the client; the accept loop continues, so a
  *new* connection can proceed normally.

Compile flags: `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan; never redefine
`CFLAGS`/`LDFLAGS`.

## Acceptance

Graded `build` + `net` + `quiz`. Single run, one connection:

```
send "ready\n"
  → expect "ready\n"
send "half"
  (grader sleeps 1200ms — server times out at 400ms)
  → expect "TIMEOUT\n"
```

- `ready\n` echoed correctly; `half` (no newline) triggers `TIMEOUT\n` within
  1200 ms.
- The server re-accepts after a timeout — a subsequent connection must still
  work.
- No timeout mechanism: the server hangs forever on the second `read` → the
  `net` deadline fires, FAIL. Timeout too long: the `net` deadline fires
  before `TIMEOUT\n` → FAIL.

`quiz.txt` is complete (see Quiz).

## Readings

- `man 2 setsockopt` (`SO_RCVTIMEO`); TLPI §23.3 "Setting Timeouts on Blocking
  Operations" for the timeout-on-sockets mechanism and why `EAGAIN`/`EWOULDBLOCK`
  are the same on Linux.
- A `poll` loop with a deadline achieves the same effect and is a common alternative
  (TLPI §63.2). Both are acceptable for this exercise.

## Quiz

1. Which socket option makes a blocking read give up after a delay?
2. What errno does a timed-out recv return?