# M7-ex01 · Backoff — your first Go

This is the module's **Go onboarding**. Everything before this was C; here you write a
Go program for the first time. It is deliberately the gentlest possible Go: a TCP
client that talks to a local **flaky peer** it spins up itself, and survives by
**retrying with exponential backoff**. You need exactly three new Go pieces:
goroutine, channel-free TCP dial/listen, and `time.Sleep`. No structs, no `iota`,
no closures yet — those arrive in M7-ex02.

## Goal

Write `main.go` — a self-contained program that (1) starts a flaky TCP server in a
goroutine, and (2) runs a client retry loop against it:

1. **The flaky server**: listens on `127.0.0.1`, accepts exactly **3** connections.
   The first two are accepted and **closed without a byte of reply** (the flake). The
   third is accepted, reads your request, and answers `connected\n`.
2. **The client retry loop**: dials that port up to **3 attempts**; between attempts it
   waits exponentially longer — base delay **10 ms**, doubling each time (10 ms, then
   20 ms). Real sleeps; *nothing about their duration may be printed*.
3. Print exactly:

```
attempt 1: closed before response
attempt 2: closed before response
attempt 3: connected
```

An attempt that gets zero response bytes is `closed before response`; one that receives
`connected\n` is a success and stops the loop.

`forge` runs `make all` to build `./test`, then runs it and diffs stdout.

## Constraints

- One port is enough: bind the flaky listener to `127.0.0.1` with port `0` (let the OS
  choose), read the chosen port back, and have the client dial it.
- No output beyond the three lines above; nothing on stderr; no timestamps or durations.
- The backoff must be real: `time.Sleep(10 * time.Millisecond)` then `2*`, `4*`… the point
  is the retry loop, not a mock.
- `Write` the full request, close the write side, read until EOF or a full `connected\n`
  answer. Treat a peer that closes without replying as a failed attempt.
- `go.mod` is provided; use only the standard library. Your `Makefile` must do exactly
  `go build -o test .` and provide `fclean`/`re`.

## Acceptance criteria

- [ ] exactly the three graded lines, in order, exit 0
- [ ] retry loop (a single-attempt program fails); success detected on the third
- [ ] increasing sleep between attempts (real exponential backoff) — verified by `go vet` +
  the grader's reference timing
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which algorithm multiplies the wait between retries as they keep failing?: <answer>
What does a successful attempt mean for a retry loop?: <answer>
Which peer behaviour does this exercise classify as a failed attempt?: <answer>
```

## Readings — the three Go pieces you need

- **Goroutine**: use `go f()` to run `net.Listen`/`Accept` in the background. Skim
  https://go.dev/tour/concurrency/1 — goroutines are the Go spelling of a thread. Unlike a
  thread you don't join it here; the retry loop is your program's "main thread".
- **TCP in Go**: `net.Listen("tcp", "127.0.0.1:0")`, `listener.Accept()`, `net.Dial("tcp", addr)`,
  `conn.Read`/`conn.Write`. Your M6 socket instincts map 1:1: read-to-EOF still means the
  peer closed its write side. https://pkg.go.dev/net
- **Timing**: `time.Sleep`, `time.Millisecond`. https://pkg.go.dev/time
- Cloudflare's "Exponential Backoff" — the industrial version of the doubling wait you
  just wrote.
- TLPI §56.5 "Stream Sockets" and §61.1 "Partial Reads and Writes on Stream Sockets" on the
  EOF/hang-up failure modes that make retries necessary.

## How you are graded

- `build` runs `make all` (must produce `./test`), then runs `./test` and diffs its stdout
  against `expected.txt` (whitespace collapsed, stderr must stay empty). No wall-clock is
  asserted — only the event transcript. `quiz` covers the concepts.