package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const total = 5

	// The worker does simulated work and reports each completion through the
	// channel; main prints only once the completion is received, so the
	// transcript is ordered by the worker, not by the scheduler.
	done := make(chan int)
	go func() {
		for t := 1; t <= total; t++ {
			time.Sleep(time.Millisecond)
			done <- t
		}
	}()

	for t := 1; t <= total; t++ {
		<-done
		fmt.Printf("task %d done\n", t)
	}

	// In-flight work has finished; nothing is left to drain. Negotiated here:
	// the graceful shutdown is a handled SIGTERM — install the handler so the
	// signal becomes a queued event instead of killing the process, raise it
	// to ourselves, then confirm and exit 0.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	<-sig
	fmt.Println("received SIGTERM")
	fmt.Println("shutdown complete")
}