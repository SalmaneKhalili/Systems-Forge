// The backend is a provided line server: the k-th frame of each connection is
// answered with R<k> <content>. Do not modify it.
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
	port, err := strconv.Atoi(os.Getenv("TARGETPORT"))
	if err != nil || port <= 0 {
		fmt.Fprintln(os.Stderr, "bad TARGETPORT")
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen failed")
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	seq := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		seq++
		content := strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
		fmt.Fprintf(conn, "R%d %s\n", seq, content)
	}
}