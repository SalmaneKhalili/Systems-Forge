package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17410"
	}
	// committed holds the ids whose participants all prepared.
	committed := map[int]bool{}
	f, err := os.Open("votes.txt")
	if err == nil {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if id, e := strconv.Atoi(strings.TrimSpace(sc.Text())); e == nil {
				committed[id] = true
			}
		}
		f.Close()
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go serve(conn, committed)
	}
}

func serve(conn net.Conn, committed map[int]bool) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 || fields[0] != "tx" {
			continue
		}
		id, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		if committed[id] {
			fmt.Fprintln(conn, "COMMIT")
		} else {
			fmt.Fprintln(conn, "ABORT")
		}
	}
}