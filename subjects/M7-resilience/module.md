# M7 · Resilience

Last module you made programs *talk* to each other. This module covers the equally
important half of operations: what your program does when the other side misbehaves — slow,
down, half-open, crashing mid-task. The five exercises are five named resilience
primitives you will meet every day in infrastructure work:

| exercise | kind      | what you build                                          |
|----------|-----------|---------------------------------------------------------|
| ex01     | retries   | a client that retries a flaky peer with exponential backoff (Python) |
| ex02     | breaker   | a circuit breaker with closed/open/half-open states (Go) |
| ex03     | health    | a health-check endpoint that reports and switches state (Go, graded over TCP) |
| ex04     | shutdown  | a worker that shuts down gracefully on SIGTERM (Go)      |
| ex05     | **gate**  | a mini-supervisor: watchdog + bounded restarts + graceful shutdown (Go) |

Most exercises are **whole-program**: write `main.go` (ex01: `main.py`) + a `Makefile` so
`make all` produces `./test`, then run and diff your printed transcript. ex02 is a
**harness shape** (like M3–M5): a `main.go` driver is provided, you implement the breaker.

Two engine notes that keep this module deterministic:

- Go is built by the same `make` framework: your `Makefile` runs `go build -o test .`
  (a `go.mod` is provided). The C sanitizer flags the grader injects are irrelevant to Go
  — they sit in env vars Go ignores.
- **Nothing you print may depend on wall-clock time**: no elapsed durations, no timestamps,
  no "after 3 s". Every transcript is event-driven. When you need real backoff/sleep for
  the *behavior*, use millisecond sleeps — but never print them. This is the rule that
  makes each exercise deterministic, both ways.