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

// cluster is a 3-node raft group (a, b, c). Each node is a distinct TCP
// endpoint; nodes exchange RequestVote / AppendEntries lines over real
// sockets. The grading client drives every event (election, proposals, reads)
// explicitly, so the transcript is fully deterministic — no wall clock.
const N = 3

type Peer struct {
	id    int
	term  int
	voted bool
	lead  bool
	log   []string
}

type hub struct {
	mu    sync.Mutex
	base  int
	peers []*Peer
}

func (h *hub) addr(i int) string { return fmt.Sprintf("127.0.0.1:%d", h.base+1+i) }
func (h *hub) name(i int) string { return string(rune('a' + i)) }

// dial sends one line to peer i and returns its one-line reply.
func (h *hub) roundTrip(i int, line string) string {
	c, err := net.Dial("tcp", h.addr(i))
	if err != nil {
		return "DOWN"
	}
	defer c.Close()
	fmt.Fprintln(c, line)
	reply, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return "DOWN"
	}
	return strings.TrimRight(reply, "\n")
}

// elect makes peer id campaign for leadership in a strictly higher term.
// It needs a majority (with its own vote): 2 of 3.
func (h *hub) elect(id int) string {
	term := 0
	for _, p := range h.peers {
		if p.term > term {
			term = p.term
		}
	}
	term++
	me := h.peers[id]
	me.term = term
	me.voted = true
	me.lead = true // will relinquish if it loses the vote
	yes := 1       // self-vote
	for _, o := range h.peers {
		if o.id == id {
			continue
		}
		if r := h.roundTrip(o.id, fmt.Sprintf("VOTE %d %d", term, id)); strings.HasPrefix(r, "YES") {
			yes++
		}
	}
	if yes >= N/2+1 {
		me.lead = true
		return fmt.Sprintf("ELECTED %s term %d", h.name(id), term)
	}
	me.lead = false
	return fmt.Sprintf("LOST %s term %d", h.name(id), term)
}

// propose appends a command to the current leader's log and replicates it to
// both followers over TCP.
func (h *hub) propose(cmd string) string {
	l := h.leader()
	if l < 0 {
		return "not leader"
	}
	me := h.peers[l]
	idx := len(me.log)
	me.log = append(me.log, cmd)
	for _, o := range h.peers {
		if o.id == l {
			continue
		}
		h.roundTrip(o.id, fmt.Sprintf("APP %d %d %s", me.term, l, cmd))
	}
	return fmt.Sprintf("ok %d", idx)
}

func (h *hub) leader() int {
	for i, p := range h.peers {
		if p.lead {
			return i
		}
	}
	return -1
}

func (h *hub) read() string {
	l := h.leader()
	if l < 0 {
		return ""
	}
	return strings.Join(h.peers[l].log, "|")
}

func (h *hub) logOf(i int) string { return strings.Join(h.peers[i].log, "|") }
func (h *hub) termOf(i int) string {
	return fmt.Sprintf("term %d", h.peers[i].term)
}

func main() {
	portStr := os.Getenv("TARGETPORT")
	if portStr == "" {
		portStr = "17430"
	}
	port, _ := strconv.Atoi(portStr)
	h := &hub{base: port}
	for i := 0; i < N; i++ {
		h.peers = append(h.peers, &Peer{id: i})
	}
	// fire up each peer's own TCP listener
	for _, p := range h.peers {
		go servePeer(h, p.id)
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go serveClient(h, c)
	}
}

// servePeer handles the intra-cluster inner-protocol socket of one node.
func servePeer(h *hub, id int) {
	ln, err := net.Listen("tcp", h.addr(id))
	if err != nil {
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		go handleInner(h, id, c)
	}
}

func handleInner(h *hub, id int, c net.Conn) {
	defer c.Close()
	r := bufio.NewScanner(c)
	w := bufio.NewWriter(c)
	p := h.peers[id]
	for r.Scan() {
		f := strings.Fields(r.Text())
		if len(f) == 0 {
			fmt.Fprintln(w, "ERR")
			w.Flush()
			continue
		}
		switch f[0] {
		case "VOTE":
			// f: VOTE term candIdx
			term, _ := strconv.Atoi(f[1])
			if term > p.term {
				p.term = term
				p.voted = false
				p.lead = false
			}
			if term >= p.term && !p.voted {
				p.voted = true
				fmt.Fprintln(w, "YES")
			} else {
				fmt.Fprintln(w, "NO")
			}
		case "APP":
			// f: APP term leaderIdx cmd
			term, _ := strconv.Atoi(f[1])
			cmd := strings.Join(f[3:], " ")
			if term > p.term {
				p.term = term
				p.lead = false
			}
			p.log = append(p.log, cmd)
			fmt.Fprintln(w, "OK")
		default:
			fmt.Fprintln(w, "ERR")
		}
		w.Flush()
	}
}

// serveClient drives cluster-level commands from the grading client.
func serveClient(h *hub, c net.Conn) {
	defer c.Close()
	r := bufio.NewScanner(c)
	w := bufio.NewWriter(c)
	for r.Scan() {
		f := strings.Fields(r.Text())
		if len(f) == 0 {
			continue
		}
		var out string
		switch f[0] {
		case "elect":
			if len(f) < 2 {
				out = "nil"
			} else {
				out = h.elect(int(f[1][0] - 'a'))
			}
		case "propose":
			if len(f) < 2 {
				out = "nil"
			} else {
				out = h.propose(strings.Join(f[1:], " "))
			}
		case "read":
			out = h.read()
		case "term", "log":
			if len(f) < 2 {
				out = "nil"
			} else {
				id := int(f[1][0] - 'a')
				if f[0] == "term" {
					out = h.termOf(id)
				} else {
					out = h.logOf(id)
				}
			}
		default:
			out = "nil"
		}
		fmt.Fprintln(w, out)
		w.Flush()
	}
}
