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

// persistentKV is a WAL-backed key-value store. Keys stay in memory sorted;
// every mutation is appended to the log before being applied.
type persistentKV struct {
	mu   sync.Mutex
	data map[string]string
	log  *os.File
}

// load replays any prior log into the store.
func load(path string) (*persistentKV, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	s := &persistentKV{data: map[string]string{}, log: f}
	raw, _ := os.ReadFile(path)
	for _, line := range strings.Split(string(raw), "\n") {
		if line == "" {
			continue
		}
		// log line: op|key|val ; DEL uses val=""
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 2 {
			continue
		}
		op, key := parts[0], parts[1]
		if op == "P" {
			val := ""
			if len(parts) == 3 {
				val = parts[2]
			}
			s.data[key] = val
		} else if op == "D" {
			delete(s.data, key)
		}
	}
	return s, nil
}

func (s *persistentKV) put(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Fprintf(s.log, "P|%s|%s\n", key, val)
	s.data[key] = val
}

func (s *persistentKV) get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *persistentKV) del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Fprintf(s.log, "D|%s\n", key)
	delete(s.data, key)
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("TARGETPORT"))
	if port == 0 {
		port = 18000
	}
	store, err := load("store.log")
	if err != nil {
		fmt.Fprintln(os.Stderr, "load:", err)
		os.Exit(1)
	}

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
		go serve(conn, store)
	}
}

func serve(conn net.Conn, store *persistentKV) {
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
		case "PUT":
			if len(fields) >= 3 {
				store.put(fields[1], fields[2])
				reply = "ok"
			} else {
				reply = "err"
			}
		case "GET":
			if len(fields) >= 2 {
				v, ok := store.get(fields[1])
				if ok {
					reply = v
				} else {
					reply = "nil"
				}
			} else {
				reply = "err"
			}
		case "DEL":
			if len(fields) >= 2 {
				store.del(fields[1])
				reply = "ok"
			} else {
				reply = "err"
			}
		default:
			reply = "err"
		}
		fmt.Fprintf(conn, "%s\n", reply)
	}
}