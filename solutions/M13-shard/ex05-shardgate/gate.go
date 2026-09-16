package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type gateway struct {
	mu     sync.Mutex
	shards []map[string]string
}

func newGateway() *gateway {
	return &gateway{shards: []map[string]string{{}, {}, {}}}
}

func shardOf(key string) int {
	if key <= "m" {
		return 0
	}
	if key <= "t" {
		return 1
	}
	return 2
}

func (g *gateway) set(key, val string) int {
	n := shardOf(key)
	g.mu.Lock()
	g.shards[n][key] = val
	g.mu.Unlock()
	return n
}

func (g *gateway) get(key string) (int, string) {
	n := shardOf(key)
	g.mu.Lock()
	val, ok := g.shards[n][key]
	g.mu.Unlock()
	if !ok {
		val = "?"
	}
	return n, val
}

func main() {
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17425"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	g := newGateway()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go serve(conn, g)
	}
}

func serve(conn net.Conn, g *gateway) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		case "set":
			if len(f) < 3 {
				continue
			}
			fmt.Fprintf(conn, "set %d\n", g.set(f[1], f[2]))
		case "get":
			if len(f) < 2 {
				continue
			}
			n, val := g.get(f[1])
			fmt.Fprintf(conn, "%d:%s\n", n, val)
		}
	}
}