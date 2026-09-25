# M7-ex02 · Circuit Breaker

ex01 retried politely; a **circuit breaker** stops calls from piling onto a
service that is already failing. It sits in front of the flaky backend, tracks
consecutive failures, trips **Open** to fail fast, and lets one **probe**
through to decide when to recover — the lesson the supervisor's restart loop
(ex05) and every rate-limited client build on. This is a harness shape: the
driver is provided, you implement the breaker.

## Shape

`forge` compiles `main.go` (provided driver) + your `breaker.go`, runs
`./test`, diffs stdout. You must provide:

```go
type State int

const (
	Closed   State = iota // every call passes; backend failures are counted
	Open                  // calls fail fast without touching the backend
	HalfOpen              // exactly one probe call is let through
)

func NewBreaker(threshold, cooldown int, now func() time.Time) *CircuitBreaker
func (b *CircuitBreaker) Allow() State     // decision for the next call
func (b *CircuitBreaker) State() State     // current state (for the driver transcript)
func (b *CircuitBreaker) Success()         // record a successful backend call
func (b *CircuitBreaker) Failure()         // record a failed backend call
```

Semantics the driver will exercise:

- **Closed**: calls are allowed. Count **consecutive** failures; a success
  resets them to zero. When the count reaches `threshold`, trip → **Open**,
  and remember `openUntil = now() + cooldown`.
- **Open**: `Allow()` fails-fast (`Open`) while `now() < openUntil`; as soon
  as `now() >= openUntil`, the next `Allow()` moves to **HalfOpen** and lets
  exactly one probe through.
- **HalfOpen**: the probe's result decides — `Success()` closes the breaker
  and resets the failure count; `Failure()` reopens it with a **fresh**
  cooldown window.

Rules:

- The `now` function is injected (a fake clock in the driver — the standard
  way to test time-based logic). Never call `time.Now()` inside the breaker;
  always use the injected one.
- The driver expects the cooldown in the same time units as the fake clock
  (the driver passes integer "clock ticks").
- No `fmt` output from `breaker.go`. Everything printed comes from the driver.
- Provide the `Makefile` so `make all` produces `./test`. Use the module's
  shared Makefile:

```make
PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go
NAME  = test

all: $(NAME)
$(NAME):
	$(GO) build -o $(NAME) .
clean:
	rm -f $(NAME)
fclean: clean
re: fclean all
.PHONY: all clean fclean re
```

## Acceptance

Graded `build` + `quiz`: the driver compiles and runs `./test`, and its stdout
is diffed against `expected.txt` (whitespace-collapsed); exit 0, empty stderr.

- The driver transcript prints the 12-line sequence in `expected.txt` exactly
  (closed passes, trip, fast-fails, half-open probe failing then succeeding).
- Success resets the consecutive-failure counter.
- The probe decides half-open → open (fail) or half-open → closed (success).

A breaker that never trips, trips too early, or lets calls through while open
produces a different transcript → FAIL. `quiz.txt` is complete (see Quiz).

## Readings

- Michael Nygard, *Release It!* — the original circuit-breaker pattern ("Circuit Breaker"
  chapter): three states, and why you stop hammering a failing dependency.
- Martin Fowler, "CircuitBreaker" — concise state machine, resets, and probe semantics.

## Quiz

1. Which breaker state refuses calls without trying the backend?
2. Which single probe call decides whether an open breaker recovers?
3. How many consecutive failures does the harness trip the breaker after?