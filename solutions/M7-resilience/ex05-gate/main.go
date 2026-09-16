package main

import (
	"fmt"
	"os"
	"time"
)

// maxCrash is the watchdog's restart budget: two consecutive failures at the
// same task and the supervisor gives up on that worker.
const maxCrash = 2

type worker struct {
	name     string
	total    int
	flakyAt  int // 1-based task index whose first attempt fails; -1 = never
	next     int // next task index to do
	attempt  int // attempts made at the current task
	failures int // consecutive failed attempts at the current task
}

// step runs the next unit of work for this worker. The flake is a *programmed*
// failure: the first attempt of the flaky task reports a crash; the retry
// succeeds. Completed tasks are preserved across the restart.
func (w *worker) step() (retry, ok bool) {
	if w.next >= w.total {
		return false, true
	}
	w.attempt++
	retry = w.attempt > 1
	ok = !(w.next+1 == w.flakyAt && w.attempt == 1)
	if ok {
		w.next++
		w.attempt = 0
		w.failures = 0
	} else {
		w.failures++
	}
	return retry, ok
}

func main() {
	workers := []*worker{
		{name: "a", total: 3, flakyAt: -1},
		{name: "b", total: 3, flakyAt: 2}, // task 2 crashes once, then succeeds
		{name: "c", total: 2, flakyAt: -1},
	}

	for {
		allDone := true
		for _, w := range workers {
			if w.next >= w.total {
				continue
			}
			allDone = false
			retry, ok := w.step()
			if ok {
				if retry {
					fmt.Printf("task %s%d retry ok\n", w.name, w.next)
				} else {
					fmt.Printf("task %s%d ok\n", w.name, w.next)
				}
				continue
			}
			fmt.Printf("task %s%d failed: watchdog restarting worker %s\n",
				w.name, w.next+1, w.name)
			time.Sleep(time.Duration(w.failures) * time.Millisecond)
			if w.failures >= maxCrash {
				fmt.Printf("giving up on worker %s\n", w.name)
				os.Exit(1)
			}
		}
		if allDone {
			break
		}
	}

	fmt.Println("all workers complete")
	fmt.Println("shutdown complete")
}