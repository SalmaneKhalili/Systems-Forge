# M7-ex02 · Circuit breaker

## Goal

A **circuit breaker** sits in front of a flaky backend and stops calls from piling onto a
service that is already failing. Implement the breaker; the driver is provided.

`forge` compiles `main.go` (provided driver) + your `breaker.go`, runs `./test`, diffs stdout.

You must provide:

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

- **Closed**: calls are allowed. Count **consecutive** failures; a success reset is zero.
  When the count reaches `threshold`, trip → **Open**, and remember `openUntil = now() + cooldown`.
- **Open**: `Allow()` fails-fast (`Open`) while `now() < openUntil`; as soon as `now() >=
  openUntil`, the next `Allow()` moves to **HalfOpen** and lets exactly one probe through.
- **HalfOpen**: the probe's result decides — `Success()` closes the breaker and resets the
  failure count; `Failure()` reopens it with a **fresh** cooldown window.

## Constraints

- The `now` function is injected (a fake clock in the driver — the standard way to test
  time-based logic). Never call `time.Now()` inside the breaker; always use the injected one.
- The driver expects the cooldown in the same time units as the fake clock (the driver
  passes integer "clock ticks").
- No `fmt` output from `breaker.go`. Everything printed comes from the driver.
- Provide the `Makefile` so `make all` produces `./test`. Use the module's shared Makefile:

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

## Acceptance criteria

- [ ] the driver transcript prints the 12-line sequence in `expected.txt` exactly
      (closed passes, trip, fast-fails, half-open probe failing then succeeding)
- [ ] success resets the consecutive-failure counter
- [ ] the probe decides half-open → open (fail) or half-open → closed (success)
- [ ] `quiz.txt` complete

Then complete `quiz.txt`:

```
Which breaker state refuses calls without trying the backend?: <answer>
Which single probe call decides whether an open breaker recovers?: <answer>
How many consecutive failures does the harness trip the breaker after?: <answer>
```

## Readings

- Michael Nygard, *Release It!* — the original circuit-breaker pattern ("Circuit Breaker"
  chapter): three states, and why you stop hammering a failing dependency.
- Martin Fowler, "CircuitBreaker" — concise state machine, resets, and probe semantics.

## How you are graded

- `build` compiles and runs `./test`; stdout is diffed against `expected.txt`
  (whitespace-collapsed) + `quiz`. A breaker that never trips, trips too early, or lets
  calls through while open produces a different transcript → FAIL.