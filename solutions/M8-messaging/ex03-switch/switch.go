package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("TARGETPORT"))
	if err != nil || port <= 0 {
		fmt.Fprintln(os.Stderr, "bad TARGETPORT")
		os.Exit(1)
	}
	backendPort := startBackend()

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
		go handle(client, backendPort)
	}
}

// startBackend picks a free loopback port, launches ./svc on it, and waits
// until the backend actually accepts connections.
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

// handle proxies one client connection to a fresh backend connection. Both
// directions run at once; the first side to end tears the pair down, and the
// next client connection is unaffected.
func handle(client net.Conn, backendPort int) {
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
	for {
		line, err := cr.ReadString('\n')
		if err != nil {
			return
		}
		if _, err := backend.Write([]byte(line)); err != nil {
			return
		}
	}
}