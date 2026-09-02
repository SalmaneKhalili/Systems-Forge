package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Metrics is a mutex-safe registry of named counters (from ex01).
type Metrics struct {
	mu   sync.Mutex
	data map[string]int64
}

func newMetrics() *Metrics {
	return &Metrics{data: make(map[string]int64)}
}

func (m *Metrics) add(name string, n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[name] += n
}

// inc counts a connection: it increments conn_total and returns its value.
func (m *Metrics) incConn() {
	m.add("conn_total", 1)
}

func (m *Metrics) get(name string) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[name]
}

// snap returns all metrics as name=value lines sorted by name.
func (m *Metrics) snap() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	names := make([]string, 0, len(m.data))
	for k := range m.data {
		names = append(names, k)
	}
	sort.Strings(names)
	out := make([]string, 0, len(names))
	for _, k := range names {
		out = append(out, fmt.Sprintf("%s=%d", k, m.data[k]))
	}
	return out
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("TARGETPORT"))
	if port == 0 {
		port = 18500
	}
	reg := newMetrics()

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}
		reg.incConn()
		go serve(conn, reg)
	}
}

func serve(conn net.Conn, reg *Metrics) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		var reply string
		switch fields[0] {
		case "INC":
			if len(fields) >= 3 {
				n, _ := strconv.ParseInt(fields[2], 10, 64)
				reg.add(fields[1], n)
				reply = "ok"
			} else {
				reply = "err"
			}
		case "GET":
			if len(fields) >= 2 {
				reply = strconv.FormatInt(reg.get(fields[1]), 10)
			} else {
				reply = "err"
			}
		case "SNAP":
			reply = strings.Join(reg.snap(), "\n")
		default:
			reply = "err"
		}
		fmt.Fprintf(conn, "%s\n", reply)
	}
}