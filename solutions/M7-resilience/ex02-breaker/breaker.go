package main

import "time"

type State int

const (
	Closed   State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	threshold int
	cooldown  int
	now       func() time.Time
	failures  int
	state     State
	openUntil time.Time
}

func NewBreaker(threshold, cooldown int, now func() time.Time) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		cooldown:  cooldown,
		now:       now,
		state:     Closed,
	}
}

// Allow checks whether a call may proceed. While open it lets the clock reach
// the cooldown boundary and then transitions to half-open; that next call is
// the probe. The boundary is inclusive: openUntil is when probes resume.
func (b *CircuitBreaker) Allow() State {
	if b.state == Open && !b.now().Before(b.openUntil) {
		b.state = HalfOpen
	}
	return b.state
}

func (b *CircuitBreaker) State() State {
	return b.state
}

func (b *CircuitBreaker) Success() {
	b.failures = 0
	if b.state == HalfOpen {
		b.state = Closed
	}
}

func (b *CircuitBreaker) Failure() {
	if b.state == HalfOpen {
		b.state = Open
		b.openUntil = b.now().Add(time.Duration(b.cooldown) * time.Second)
		return
	}
	b.failures++
	if b.failures >= b.threshold {
		b.state = Open
		b.openUntil = b.now().Add(time.Duration(b.cooldown) * time.Second)
	}
}