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

// durable is a transactional key-value gateway hardened by a write-ahead log.
// It composes M15 transactions (staged writes, atomic commit, rollback) with
// M16 durability (every committed write is appended to an in-memory WAL that a
// crash replay could restore). The `wal` command exposes that durable record.
type durable struct {
	mu   sync.Mutex
	data map[string]string
	wal  []string // committed ops, one per line: "set k=v" in commit order
}

func (d *durable) apply(k, v string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[k] = v
	d.wal = append(d.wal, "set "+k+"="+v)
}

func (d *durable) get(k string) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.data[k]
}

func (d *durable) has(k string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.data[k]
	return ok
}

func (d *durable) keys() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	ks := make([]string, 0, len(d.data))
	for k := range d.data {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func (d *durable) walSnapshot() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return strings.Join(d.wal, "\n")
}

func main() {
	d := &durable{data: map[string]string{}}
	port := os.Getenv("TARGETPORT")
	if port == "" {
		port = "17540"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(d, c)
	}
}

func handle(d *durable, c net.Conn) {
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
		f := strings.Fields(r.Text())
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
				for k, v := range staged {
					d.apply(k, v)
				}
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
				d.apply(f[1], f[2])
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
			if !d.has(k) {
				reply("nil")
				continue
			}
			reply(d.get(k))
		case "status":
			if len(f) < 2 {
				reply("nil")
				continue
			}
			if d.has(f[1]) {
				reply("alive")
			} else {
				reply("dead")
			}
		case "list":
			reply(strings.Join(d.keys(), "\n"))
		case "wal":
			reply(d.walSnapshot())
		default:
			reply("nil")
		}
	}
}
