package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type replica struct {
	mu      sync.Mutex
	entries []string
	commit  int
}

func (r *replica) append(cmd string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	idx := len(r.entries)
	r.entries = append(r.entries, cmd)
	return idx
}

func (r *replica) setCommit(idx int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if idx > r.commit {
		r.commit = idx
	}
	return r.commit
}

func (r *replica) read() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := r.commit + 1
	if n > len(r.entries) {
		n = len(r.entries)
	}
	return strings.Join(r.entries[:n], "|")
}

func main() {
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17415"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	r := &replica{commit: -1}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go serve(conn, r)
	}
}

func serve(conn net.Conn, r *replica) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		case "append":
			if len(f) < 2 {
				continue
			}
			idx := r.append(f[1])
			fmt.Fprintf(conn, "ok %d\n", idx)
		case "commit":
			if len(f) < 2 {
				continue
			}
			idx, err := strconv.Atoi(f[1])
			if err != nil {
				continue
			}
			fmt.Fprintf(conn, "commit %d\n", r.setCommit(idx))
		case "read":
			fmt.Fprintf(conn, "%s\n", r.read())
		}
	}
}