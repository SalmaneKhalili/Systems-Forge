package main

import (
	"fmt"
	"time"
)

// results is the scripted backend: each element is consumed by exactly one
// call that is allowed (closed pass or half-open probe). Fast-fails consume
// nothing, so the transcript stays deterministic.
var results = []bool{true, true, false, false, false, true, true, true}

func main() {
	clock := 0
	now := func() time.Time {
		return time.Unix(int64(clock), 0) // reads, does not advance
	}

	b := NewBreaker(2, 3, now)

	next := 0
	var ok bool

	for i := 0; i < 12; i++ {
		clock++ // one logical tick per call
		switch b.Allow() {
		case Open:
			fmt.Printf("call %d: open fast-fail\n", i+1)
			continue
		case HalfOpen:
			ok = results[next]
			next++
			if ok {
				b.Success()
			} else {
				b.Failure()
			}
			fmt.Printf("call %d: half-open probe %s -> %s\n", i+1, result(ok), stateName(b.State()))
		case Closed:
			ok = results[next]
			next++
			if ok {
				b.Success()
			} else {
				b.Failure()
			}
			fmt.Printf("call %d: closed %s -> %s\n", i+1, result(ok), stateName(b.State()))
		}
	}
}

func result(ok bool) string {
	if ok {
		return "ok"
	}
	return "fail"
}

func stateName(s State) string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	default:
		return "half-open"
	}
}