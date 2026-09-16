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

type node struct {
	mu      sync.Mutex
	term    int
	entries []string
}

func (n *node) see(t int) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	if t > n.term {
		n.term = t
	}
	return n.term
}

func (n *node) append(cmd string) (int, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.term < 1 {
		return 0, false
	}
	idx := len(n.entries)
	n.entries = append(n.entries, cmd)
	return idx, true
}

func (n *node) read() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	return strings.Join(n.entries, "|")
}

func main() {
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17420"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	n := &node{}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go serve(conn, n)
	}
}

func serve(conn net.Conn, n *node) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		case "see":
			if len(f) < 2 {
				continue
			}
			t, err := strconv.Atoi(f[1])
			if err != nil {
				continue
			}
			fmt.Fprintf(conn, "term %d\n", n.see(t))
		case "append":
			if len(f) < 2 {
				continue
			}
			if idx, ok := n.append(f[1]); ok {
				fmt.Fprintf(conn, "ok %d\n", idx)
			} else {
				fmt.Fprintf(conn, "not leader\n")
			}
		case "read":
			fmt.Fprintf(conn, "%s\n", n.read())
		}
	}
}