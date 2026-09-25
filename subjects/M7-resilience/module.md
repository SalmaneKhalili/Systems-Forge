# M7 · Resilience

Last module you made programs *talk* to each other. This module covers the
equally important half of operations: what your program does when the other
side misbehaves — slow, down, half-open, crashing mid-task. The five exercises
are five named resilience primitives you will meet every day in infrastructure
work: backoff, a circuit breaker, health checks, graceful shutdown, and a
watchdog — and it is also where Go is introduced, so each primitive is built
with the goroutines and channels that make supervising concurrent work
tractable.

## The build

- **ex01 · Backoff: your first Go** — a self-contained TCP client that spins
  up a flaky peer and survives by retrying with exponential backoff; the
  gentlest possible Go onboarding (goroutine, `net.Listen`/`net.Dial`,
  `time.Sleep`), and the concurrency pieces every later exercise reuses.
- **ex02 · Circuit breaker** — harness shape: implement `breaker.go` behind a
  provided driver; Closed/Open/HalfOpen with an injected clock, fail-fast, and
  the one-probe decision the supervisor's restart loop depends on.
- **ex03 · Health checks** — a `GET /health` endpoint that reports live or
  dead and can be flipped over the wire (`/down`, `/up`), state surviving
  across connections; graded over real TCP.
- **ex04 · Graceful shutdown** — a worker sequences its completions through a
  channel, raises SIGTERM to itself, and drains instead of dying on the spot.
- **ex05 · Gate: mini supervisor** — a watchdog that catches a worker
  crashing mid-shift, restarts it a bounded number of times with real backoff
  without replaying completed work, then shuts down gracefully: the smallest
  honest version of what runs your containers.

## Rules

- Almost every exercise is whole-program: write `main.go` (ex01 included) plus
  a `Makefile` so `make all` produces `./test`, then run and diff your printed
  transcript. ex02 is a **harness shape** (like M3–M5): the `main.go` driver
  is provided, you implement the breaker.
- Go is built by the same `make` framework: your `Makefile` runs
  `go build -o test .` (a `go.mod` is provided). The C sanitizer flags the
  grader injects are irrelevant to Go — they sit in env vars Go ignores.
- **Nothing you print may depend on wall-clock time**: no elapsed durations,
  no timestamps, no "after 3 s". Every transcript is event-driven. When you
  need real backoff/sleep for the *behavior*, use millisecond sleeps — but
  never print them. This is the rule that makes each exercise deterministic,
  both ways.
- ex03 is additionally graded over TCP: bind `127.0.0.1` on the `TARGETPORT`
  environment variable, exactly as in M6.

## Prerequisites

M7 is where **Go is introduced**. If you have never used Go, do the [Go Tour](https://tour.golang.org)
first (takes ~1 hour). Then come back and verify you can do these specific things:

**Go language basics (you must know before starting):**

- Write a function with named parameters and return values: `func foo(x int, y string) (int, error)`.
- Define a struct and add methods to it: `type Foo struct { ... }` / `func (f *Foo) Bar()`.
- Use `iota` to define a set of related constants: `const ( A = iota; B; C )`.
- Write a closure: `func() { ... }` passed as an argument or returned from a function.
- Use `if err != nil { ... }` — Go's error handling pattern.
- Declare a goroutine: `go doSomething()` — starts a new concurrent function.
- Use a channel: `ch <- val` (send), `val := <-ch` (receive), `<-chan int` (type).
- Use a `select` block to wait on multiple channel operations with a `default` case.
- Call `time.Sleep(d)` and `time.Duration` arithmetic (`100 * time.Millisecond`).

**Go standard library (used in this module):**

- `net.Listen("tcp", addr)` — returns a `net.Listener`; call `.Accept()` in a loop to get
  `net.Conn` values.
- `conn.Read(buf)` / `conn.Write(data)` — like file I/O but over the network.
- `bufio.NewScanner(conn)` and `scanner.Text()` — line-based reading.
- `os.Getenv("VARNAME")` — read environment variables.
- `os.Exit(1)` / `syscall.SIGTERM` — process lifecycle signals.
- `os/signal.Notify(ch, syscall.SIGTERM)` — channel receives OS signals.

**Build tooling:**

- Run `go build -o test .` to compile a Go program.
- Know what `go.mod` is and that `module <name>` declares the module path.
- `make` wrapping `go build` (the Makefiles here do this — read them before writing).

**Concurrency patterns (you will learn the details here, but know the shape):**

- "Fan out" = multiple goroutines reading from different connections.
- "Graceful shutdown" = receive SIGTERM, drain connections, exit cleanly.
- A circuit breaker = track consecutive failures, stop trying after threshold, retry after
  cooldown.

If you've never written Go before, the first exercise (M7-ex01 backoff) is deliberately a
gentle onboarding: goroutine + `net.Listen`/`net.Dial` + `time.Sleep`. The second exercise
(circuit breaker) is where the pace picks up.

## So what? (interview / portfolio)

Resilience primitives are what "production-grade" actually means in infrastructure roles —
half of your on-call job is handling a peer that is slow, down, or half-open. This module
gives you the named vocabulary (retry/backoff, circuit breaker, health check, graceful
shutdown, watchdog) that interviewers and architecture discussions assume you already have.

**Interview questions this module arms you for:**
- How does exponential backoff work, and when is retrying the wrong answer?
- Explain the closed / open / half-open states of a circuit breaker.
- What should a health-check endpoint report, and how does a supervisor act on it?
- Graceful shutdown on SIGTERM: what do you drain, and why does it matter?

**Portfolio artifact:** M7-ex05 `mini-supervisor` — a watchdog that restarts a worker
boundedly and shuts down gracefully.