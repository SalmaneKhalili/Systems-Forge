package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

const okBody = "HTTP/1.1 200 OK\r\nContent-Length: 3\r\n\r\nok\n"
const deadBody = "HTTP/1.1 503 Service Unavailable\r\nContent-Length: 0\r\n\r\n"
const notFound = "HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"

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

	healthy := true
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		serve(conn, &healthy)
	}
}

func serve(conn net.Conn, healthy *bool) {
	defer conn.Close()
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	switch line {
	case "GET /health HTTP/1.1\r\n":
		if *healthy {
			fmt.Fprint(conn, okBody)
		} else {
			fmt.Fprint(conn, deadBody)
		}
	case "GET /down HTTP/1.1\r\n":
		*healthy = false
		fmt.Fprint(conn, okBody)
	case "GET /up HTTP/1.1\r\n":
		*healthy = true
		fmt.Fprint(conn, okBody)
	default:
		fmt.Fprint(conn, notFound)
	}
}