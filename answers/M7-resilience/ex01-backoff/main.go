package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// flaky accepts exactly 3 connections: it closes the first two without a
// byte of reply and answers "connected\n" on the third.
func flaky() (string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	go func() {
		for seen := 0; seen < 3; seen++ {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			if seen < 2 {
				conn.Close()
				continue
			}
			buf := make([]byte, 4096)
			for {
				n, err := conn.Read(buf)
				if err != nil || n == 0 || buf[n-1] == '\n' {
					break
				}
			}
			conn.Write([]byte("connected\n"))
			conn.Close()
		}
		ln.Close()
	}()
	return ln.Addr().String(), nil
}

// attempt does one try: false means "closed before a response" (EOF or RST).
func attempt(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	defer conn.Close()
	if _, err := conn.Write([]byte("PING\n")); err != nil {
		return false
	}
	data := make([]byte, 0, 64)
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
			if data[len(data)-1] == '\n' {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return string(data) == "connected\n"
}

func main() {
	addr, err := flaky()
	if err != nil {
		fmt.Fprintln(os.Stderr, "server failed")
		os.Exit(1)
	}
	base := 10 * time.Millisecond
	for n := 1; n <= 3; n++ {
		if attempt(addr) {
			fmt.Printf("attempt %d: connected\n", n)
			return
		}
		fmt.Printf("attempt %d: closed before response\n", n)
		if n < 3 {
			time.Sleep(time.Duration(1<<(uint(n-1))) * base) // 10ms, then 20ms
		}
	}
}