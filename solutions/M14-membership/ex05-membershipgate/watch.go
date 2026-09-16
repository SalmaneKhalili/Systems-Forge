package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
)

type watcher struct {
	mu     sync.Mutex
	states map[string]string
}

func (w *watcher) beat(name string) {
	w.mu.Lock()
	w.states[name] = "alive"
	w.mu.Unlock()
}

func (w *watcher) suspect(name string) {
	w.mu.Lock()
	if w.states[name] == "alive" {
		w.states[name] = "suspect"
	}
	w.mu.Unlock()
}

func (w *watcher) drop(name string) {
	w.mu.Lock()
	w.states[name] = "dead"
	w.mu.Unlock()
}

func (w *watcher) status(name string) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	s, ok := w.states[name]
	if !ok {
		return "dead"
	}
	return s
}

func (w *watcher) list() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	names := []string{}
	for n, s := range w.states {
		if s == "alive" || s == "suspect" {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return strings.Join(names, "|")
}

func main() {
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17435"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	w := &watcher{states: map[string]string{}}

	// Graceful shutdown: on SIGTERM, stop accepting and exit 0. This
	// re-deploys the M7-ex04 skill — a server should not die mid-flight.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		<-stop
		ln.Close()
		close(done)
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-done:
				return
			default:
				continue
			}
		}
		go serve(conn, w)
	}
}

func serve(conn net.Conn, w *watcher) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		case "list":
			fmt.Fprintf(conn, "%s\n", w.list())
		case "status":
			if len(f) < 2 {
				continue
			}
			fmt.Fprintf(conn, "%s\n", w.status(f[1]))
		case "beat", "suspect", "drop":
			if len(f) < 2 {
				continue
			}
			switch f[0] {
			case "beat":
				w.beat(f[1])
			case "suspect":
				w.suspect(f[1])
			case "drop":
				w.drop(f[1])
			}
			fmt.Fprintf(conn, "ok\n")
		}
	}
}