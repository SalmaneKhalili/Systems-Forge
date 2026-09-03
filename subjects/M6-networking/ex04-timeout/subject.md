# M6-ex04 · read timeout

## Goal

Write `main.c` so `make all` produces `./test` — a **line-echo server with a read
timeout**: every connection that has no data for **400 ms** is terminated with
`TIMEOUT\n`.

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

Graded dialogue (single run, one connection):

```
send "ready\n"
  → expect "ready\n"
send "half"
  (grader sleeps 1200ms — server times out at 400ms)
  → expect "TIMEOUT\n"
```

## Constraints

- **`SO_RCVTIMEO`** (or `poll`/`select` with timeout — both are acceptable) on the read
  side. The exact mechanism is yours; the point is that a stuck client cannot wedge the
  server.
- Timeout = **400 ms**. Your server must fail before the grader's 1200 ms delay ends.
- A complete line echoed with correct newline; a timed-out connection must produce exactly
  `TIMEOUT\n`.
- After the timeout event: close the client; the accept loop continues, so a *new*
  connection can proceed normally.
- `-std=gnu11 -Wall -Wextra -Werror`, ASan/UBSan.

## Acceptance criteria

- [ ] `ready\n` echoed correctly; `half` (no newline) triggers `TIMEOUT\n` within 1200 ms
- [ ] server re-accepts after a timeout — subsequent connection must still work
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which socket option makes a blocking read give up after a delay?: <answer>
What errno does a timed-out recv return?: <answer>
```

## Readings

- `man 2 setsockopt` (`SO_RCVTIMEO`); TLPI §23.3 "Setting Timeouts on Blocking
  Operations" for the timeout-on-sockets mechanism and why `EAGAIN`/`EWOULDBLOCK`
  are the same on Linux.
- A `poll` loop with a deadline achieves the same effect and is a common alternative
  (TLPI §63.2). Both are acceptable for this exercise.

## How you are graded

- `build` strict; `net` sends the exact script, sleeps between the two sends, and expects
  the server to have timed out before its own deadline. No timeout mechanism: server hangs
  forever on the second `read` → `net` deadline fires, FAIL. Timeout too long: `net`
  deadline fires before `TIMEOUT\n` → FAIL.
- `quiz`: `quiz.txt` answers must match.