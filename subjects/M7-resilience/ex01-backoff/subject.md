# M7-ex01 · Backoff

## Goal

Write `main.py` — a client that talks to a **flaky peer** (a local TCP service that closes
the socket instead of answering) and survives by **retrying with exponential backoff**.

Your program must:

1. Start a server in a background thread: it accepts exactly 3 connections. The first two
   are accepted and **closed without a byte of reply** (the flake). The third is accepted,
   reads your request, and answers `connected\n`.
2. Run a client retry loop against that server: up to **3 attempts**, and **between
   attempts, wait exponentially longer** (a small backoff — base delay 10 ms, doubling
   each time: 10 ms, then 20 ms…). Real `sleep`s; nothing about their duration may be
   printed.
3. Classify each attempt and print exactly:

```
attempt 1: closed before response
attempt 2: closed before response
attempt 3: connected
```

An attempt that gets zero response bytes is `closed before response`; an attempt that
receives `connected\n` is a success and stops the loop.

## Constraints

- One port is enough: bind the flaky server to `127.0.0.1`, let it choose the port, and let
  the client dial that port.
- No output beyond the three lines above; nothing on stderr; no timestamps or durations.
- The backoff must be real (the flaky peer is only a simulation — the point is the retry
  loop). Between attempts, `time.sleep` for `base * 2**n` ms.
- `send` the full request; expect a line-terminated reply; treat EOF as failure, a full
  `connected\n` as success.

## Acceptance criteria

- [ ] exactly the three graded lines, in order, exit 0
- [ ] retry loop (single-attempt programs fail), success detected on the third
- [ ] increasing sleep between attempts (real exponential backoff)
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which algorithm multiplies the wait between retries as they keep failing?: <answer>
What does a successful attempt mean for a retry loop?: <answer>
Which peer behaviour does this exercise classify as a failed attempt?: <answer>
```

## Readings

- `man 2 socket`, `man 2 connect`, `man 2 recv`, `man 2 send`.
- Cloudflare's "Exponential Backoff" guidance — the canonical pairing of retry counts and
  doubling waits, plus jitter — is the industrial version of what you just wrote.
- TLPI §62 on the failure modes of a peer (EOF/hang-up) that make retries necessary.

## How you are graded

- `stdout` runs `python3 main.py` and diffs output against `expected.txt` (whitespace
  collapsed, stderr must stay empty) + `quiz`. Nothing about timing is asserted — only the
  event transcript — which is what keeps the retry behaviour deterministic to grade.