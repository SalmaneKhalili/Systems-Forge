package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// sequencer keeps one logical clock for the whole server lifetime.
type sequencer struct {
	seq int
}

// stamp applies the Lamport receive rule: seq = max(seq, claimed) + 1.
func (s *sequencer) stamp(claimed int) int {
	if claimed > s.seq {
		s.seq = claimed
	}
	s.seq++
	return s.seq
}

func main() {
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17405"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	s := &sequencer{}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go serve(conn, s)
	}
}

func serve(conn net.Conn, s *sequencer) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 {
			continue
		}
		claimed, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		n := s.stamp(claimed)
		fmt.Fprintf(conn, "R%d %s %d\n", n, fields[0], claimed)
	}
}