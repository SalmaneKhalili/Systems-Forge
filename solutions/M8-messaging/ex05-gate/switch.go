package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type directive struct {
	op    string // dup | drop | hold
	frame int    // 1-based client frame index within the current connection
}

func main() {
	port, err := strconv.Atoi(os.Getenv("TARGETPORT"))
	if err != nil || port <= 0 {
		fmt.Fprintln(os.Stderr, "bad TARGETPORT")
		os.Exit(1)
	}
	backendPort := startBackend()
	faults := loadFaults(os.Getenv("TARGETFAULTS"))

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen failed")
		os.Exit(1)
	}
	for {
		client, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(client, backendPort, faults)
	}
}

func startBackend() int {
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintln(os.Stderr, "backend: no free port")
		os.Exit(1)
	}
	port := probe.Addr().(*net.TCPAddr).Port
	probe.Close()

	cmd := exec.Command("./svc")
	cmd.Env = append(os.Environ(), "TARGETPORT="+strconv.Itoa(port))
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "backend: cannot start")
		os.Exit(1)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 200*time.Millisecond)
		if err == nil {
			conn.Close()
			return port
		}
		if time.Now().After(deadline) {
			fmt.Fprintln(os.Stderr, "backend: not ready")
			os.Exit(1)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// loadFaults reads one op:index per line. A missing file, a blank line, or an
// unparsable index is skipped — the relay stays transparent for that frame.
func loadFaults(path string) []directive {
	if path == "" {
		path = "faults.txt"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var ds []directive
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		sep := strings.IndexByte(line, ':')
		if sep < 0 {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(line[sep+1:]))
		if err != nil || n <= 0 {
			continue
		}
		ds = append(ds, directive{op: strings.TrimSpace(line[:sep]), frame: n})
	}
	return ds
}

// handle proxies one client connection to a fresh backend connection,
// applying the fault directives to the client→backend direction. Frame
// numbers restart at 1 per connection; a held frame is released only after
// the following forwarded frame, in that exact order.
func handle(client net.Conn, backendPort int, faults []directive) {
	defer client.Close()

	backend, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", backendPort))
	if err != nil {
		return
	}
	defer backend.Close()

	go func() {
		br := bufio.NewReader(backend)
		for {
			line, err := br.ReadString('\n')
			if err != nil {
				client.Close()
				return
			}
			if _, err := io.WriteString(client, line); err != nil {
				return
			}
		}
	}()

	cr := bufio.NewReader(client)
	k := 0
	var held []byte
	for {
		line, err := cr.ReadString('\n')
		if err != nil {
			return
		}
		k++
		if op, ok := hit(faults, k); ok {
			switch op {
			case "drop":
				continue
			case "dup":
				if _, err := backend.Write([]byte(line)); err != nil {
					return
				}
				if _, err := backend.Write([]byte(line)); err != nil {
					return
				}
				continue
			case "hold":
				held = []byte(line)
				continue
			}
		}
		if _, err := backend.Write([]byte(line)); err != nil {
			return
		}
		if held != nil {
			if _, err := backend.Write(held); err != nil {
				return
			}
			held = nil
		}
	}
}

func hit(faults []directive, k int) (string, bool) {
	for _, d := range faults {
		if d.frame == k {
			return d.op, true
		}
	}
	return "", false
}