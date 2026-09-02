package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	mu   sync.Mutex
	data = map[string]string{}
)

func set(k, v string) {
	mu.Lock()
	defer mu.Unlock()
	data[k] = v
}

func get(k string) string {
	mu.Lock()
	defer mu.Unlock()
	return data[k]
}

func has(k string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := data[k]
	return ok
}

func drop(k string) {
	mu.Lock()
	defer mu.Unlock()
	delete(data, k)
}

// list returns the live keys, one per line, sorted; empty if none.
func list() string {
	mu.Lock()
	var ks []string
	for k := range data {
		ks = append(ks, k)
	}
	mu.Unlock()
	sort.Strings(ks)
	return strings.Join(ks, "\n")
}

func handle(c net.Conn) {
	defer c.Close()
	r := bufio.NewScanner(c)
	w := bufio.NewWriter(c)
	inTxn := false
	staged := map[string]string{}

	reply := func(s string) {
		fmt.Fprintln(w, s)
		w.Flush()
	}

	for r.Scan() {
		line := r.Text()
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		case "begin":
			inTxn = true
			staged = map[string]string{}
			reply("ok")
		case "commit":
			if inTxn {
				mu.Lock()
				for k, v := range staged {
					data[k] = v
				}
				mu.Unlock()
				staged = map[string]string{}
				inTxn = false
			}
			reply("ok")
		case "rollback":
			staged = map[string]string{}
			inTxn = false
			reply("ok")
		case "set":
			if len(f) < 3 {
				reply("nil")
				continue
			}
			if inTxn {
				staged[f[1]] = f[2]
			} else {
				set(f[1], f[2])
			}
			reply("ok")
		case "get":
			if len(f) < 2 {
				reply("nil")
				continue
			}
			k := f[1]
			if inTxn {
				if v, ok := staged[k]; ok {
					reply(v)
					continue
				}
			}
			if !has(k) {
				reply("nil")
				continue
			}
			reply(get(k))
		case "status":
			if len(f) < 2 {
				reply("nil")
				continue
			}
			if has(f[1]) {
				reply("alive")
			} else {
				reply("dead")
			}
		case "drop":
			if len(f) < 2 {
				reply("nil")
				continue
			}
			drop(f[1])
			reply("ok")
		case "list":
			reply(list())
		default:
			reply("nil")
		}
	}
}

func main() {
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17440"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(c)
	}
}