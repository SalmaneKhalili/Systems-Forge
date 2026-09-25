package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/creack/pty"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TermSession is one running pty-backed process (an interactive shell or an
// on-demand nvim). Output streams to the frontend as base64 chunks.
type TermSession struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	ptmx   *os.File
	closed bool
}

var (
	termMu      sync.Mutex
	termSession = map[string]*TermSession{}
)

// TermStart launches a process under a pty in dir. command may be empty to
// spawn the user's interactive shell. Output is emitted as the
// "term:out:<id>" event (base64 payload); the "term:exit:<id>" event fires
// when the process finishes.
func (a *App) TermStart(id, dir string, command []string) error {
	fmt.Fprintf(os.Stderr, "[term] start id=%s dir=%q argv=%v\n", id, dir, command)
	argv := command
	if len(argv) == 0 {
		sh := os.Getenv("SHELL")
		if sh == "" {
			sh = "/bin/sh"
		}
		argv = []string{sh}
	}
	if !filepath.IsAbs(argv[0]) {
		p, err := exec.LookPath(argv[0])
		if err != nil {
			return fmt.Errorf("command %q not found: %w", argv[0], err)
		}
		argv[0] = p
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	if dir != "" {
		if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
			return fmt.Errorf("bad working directory %q", dir)
		}
		cmd.Dir = dir
	}
	cmd.Env = append(os.Environ(), "TERM=xterm-256color", "COLORTERM=truecolor")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[term] %s pty.Start(%q): %v\n", id, argv[0], err)
		return err
	}
	// Give the fresh pty a sane size before the frontend's first resize call,
	// so fullscreen-ish startup probes (fish, nvim) never see a 0x0 window.
	_ = pty.Setsize(ptmx, &pty.Winsize{Cols: 80, Rows: 24})
	termMu.Lock()
	termSession[id] = &TermSession{cmd: cmd, ptmx: ptmx}
	termMu.Unlock()

	go a.pump(id, ptmx)
	return nil
}

// pump forwards pty output to the frontend and cleans up on exit. A panic
// here must never kill the app, so it is recovered and logged.
func (a *App) pump(id string, ptmx *os.File) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "[term] %s pump panic: %v\n", id, r)
		}
	}()
	emit := func(event string, payload any) {
		if a.ctx == nil {
			return
		}
		runtime.EventsEmit(a.ctx, event, payload)
	}
	resp := &termResponder{}
	readPty(ptmx, func(chunk []byte) {
		emit("term:out:"+id, base64.StdEncoding.EncodeToString(chunk))
		if reply := resp.feed(chunk); len(reply) > 0 {
			_, _ = ptmx.Write(reply)
		}
	})
	_ = ptmx.Close()

	termMu.Lock()
	ts := termSession[id]
	delete(termSession, id)
	termMu.Unlock()

	if ts != nil {
		ts.mu.Lock()
		ts.closed = true
		_ = ts.cmd.Wait() // reap after the process has exited
		ts.mu.Unlock()
	}
	emit("term:exit:"+id, true)
}

// readPty forwards pty output to onChunk until EOF (the testable core of
// pump). Each chunk is copied so callers may retain it safely.
func readPty(ptmx *os.File, onChunk func([]byte)) {
	buf := make([]byte, 32*1024)
	for {
		n, err := ptmx.Read(buf)
		if n > 0 {
			onChunk(append([]byte(nil), buf[:n]...))
		}
		if err != nil {
			return
		}
	}
}

// termResponder is a miniature terminal emulator side: interactive shells and
// editors probe the terminal with capability queries (DECRQM, XTGETTCAP, OSC
// color) and some, notably fish and bash, block waiting for a reply. The
// frontend xterm.js answers several of these, but not all, so unanswered ones
// stall the shell mid-input. The responder catches the query bytes flowing
// out of the pty and writes the expected answer back, making the session
// immune to whatever the frontend terminal covers. Known answers match the
// queries fired by fish, bash and nvim with TERM=xterm-256color.
type termResponder struct {
	buf []byte // unparsed output window
}

const (
	termMaxWindow = 8 * 1024 // drop a window that never terminates (queries are tiny)
)

// feed scans a chunk of pty output for query sequences and returns the canned
// replies to write back into the pty. Non-query bytes are consumed silently;
// incomplete sequences are kept for the next chunk.
func (r *termResponder) feed(b []byte) []byte {
	if len(r.buf) > termMaxWindow {
		r.buf = r.buf[:0]
	}
	r.buf = append(r.buf, b...)
	var replies [][]byte
	for len(r.buf) > 0 {
		if r.buf[0] != 0x1b {
			r.buf = r.buf[1:]
			continue
		}
		n := matchQuery(r.buf)
		if n == 0 {
			break // waiting for the sequence to complete
		}
		q := r.buf[:n]
		r.buf = r.buf[n:]
		if ans, ok := termAnswer(q); ok {
			replies = append(replies, ans)
		}
	}
	return bytes.Join(replies, nil)
}

// matchQuery reports the byte length of a complete query sequence at the start
// of buf: OSC string (ESC ] … ST), DCS string (ESC P … ST), or CSI (ESC [ …
// final byte). Returns 0 when more bytes are needed.
func matchQuery(buf []byte) int {
	switch {
	case buf[0] == 0x1b && len(buf) > 1 && buf[1] == ']': // OSC … ST
		for i := 2; i < len(buf); i++ {
			if buf[i] == 0x07 {
				return i + 1
			}
			if i+1 < len(buf) && buf[i] == 0x1b && buf[i+1] == '\\' {
				return i + 2
			}
		}
	case buf[0] == 0x1b && len(buf) > 1 && buf[1] == 'P': // DCS … ST
		for i := 2; i < len(buf); i++ {
			if i+1 < len(buf) && buf[i] == 0x1b && buf[i+1] == '\\' {
				return i + 2
			}
		}
	case buf[0] == 0x1b && len(buf) > 1 && buf[1] == '[': // CSI … final
		for i := 2; i < len(buf); i++ {
			if c := buf[i]; c >= 0x40 && c <= 0x7e {
				return i + 1
			}
		}
	}
	return 0
}

// termAnswer maps a known query to its reply. Unrecognized queries return ok
// == false and are passed on to the frontend unmodified.
func termAnswer(q []byte) ([]byte, bool) {
	// Keyboard-protocol query (fish, bash 5.2+): CSI ? u → basic kitty-style
	// keyboard reporting so the shell proceeds past its probe.
	if bytes.Equal(q, []byte{0x1b, '[', '?', 'u'}) {
		return []byte{0x1b, '[', '?', '1', 'u'}, true
	}
	// DA (device attributes) queries come in several shapes.
	if bytes.Equal(q, []byte{0x1b, '[', '0', 'c'}) || bytes.Equal(q, []byte{0x1b, '[', 'c'}) {
		return []byte{0x1b, '[', '?', '6', '2', ';', '2', '2', ';', 'c'}, true
	}
	// DA1 variant used by a few shells/configs: CSI ? 1 ; 2 c.
	if bytes.Equal(q, []byte{0x1b, '[', '?', '1', ';', '2', 'c'}) || bytes.Equal(q, []byte{0x1b, '[', '?', 'c'}) {
		return []byte{0x1b, '[', '?', '6', '2', ';', '2', '2', ';', 'c'}, true
	}
	// DA2 (secondary device attributes): CSI > c. A bare xterm-ish answer is
	// enough to satisfy probes that otherwise stall waiting for a reply.
	if bytes.Equal(q, []byte{0x1b, '[', '>', 'c'}) {
		return []byte{0x1b, '[', '>', '0', ';', '0', ';', '0', 'c'}, true
	}
	// XTGETTCAP by terminal name (CSI > 0 q): answer as an xterm-ish device
	// with KITTY attribute so private-mode queries are honored.
	if bytes.HasPrefix(q, []byte{0x1b, '[', '>', '0', 'q'}) {
		return []byte{0x1b, '[', '>', '0', ';', '0', ';', '5', 'q'}, true
	}
	// DECRQM (CSI ? Ps $ p): report "not recognized / default" (Pa=0) so the
	// caller proceeds; also covers the DECRQM *p variant from some tools.
	if len(q) > 6 && q[0] == 0x1b && q[1] == '[' && q[2] == '?' && q[len(q)-2] == '$' && (q[len(q)-1] == 'p' || q[len(q)-1] == '*') {
		digits := q[3 : len(q)-2]
		return append(append([]byte{0x1b, '[', '?'}, digits...), []byte{';', '0', '$', 'y'}...), true
	}
	// DSR device-status reports: CSI 5 n / CSI 6 n (cursor position).
	if bytes.Equal(q, []byte{0x1b, '[', '5', 'n'}) {
		return []byte{0x1b, '[', '0', 'n'}, true
	}
	if bytes.Equal(q, []byte{0x1b, '[', '6', 'n'}) {
		return []byte{0x1b, '[', '1', ';', '1', 'R'}, true
	}
	// OSC 11 background-colour query.
	if bytes.HasPrefix(q, []byte{0x1b, ']', '1', '1', ';', '?'}) {
		return []byte{0x1b, ']', '1', '1', ';', 'r', 'g', 'b', ':', '0', 'f', '/', '0', 'f', '/', '1', '5', 0x1b, '\\'}, true
	}
	// XTGETTCAP DCS (ESC P + q <name> ST): echo an empty but valid response.
	if bytes.HasPrefix(q, []byte{0x1b, 'P', '+', 'q'}) {
		name := bytes.TrimSuffix(q[3:], []byte{0x1b, '\\'})
		return append(append([]byte{0x1b, 'P', '1', '+', 'r'}, name...), append([]byte{'='}, append(name, []byte{0x1b, '\\'}...)...)...), true
	}
	return nil, false
}

// TermInput writes frontend keystrokes into the pty.
func (a *App) TermInput(id, data string) {
	termMu.Lock()
	ts := termSession[id]
	termMu.Unlock()
	if ts == nil {
		return
	}
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if !ts.closed {
		_, _ = ts.ptmx.Write([]byte(data))
	}
}

// TermResize adjusts the pty window size to match the xterm element.
func (a *App) TermResize(id string, cols, rows int) {
	termMu.Lock()
	ts := termSession[id]
	termMu.Unlock()
	if ts == nil || cols < 2 || rows < 2 {
		return
	}
	_ = pty.Setsize(ts.ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

// TermStop kills the process; the pump goroutine emits the exit event.
func (a *App) TermStop(id string) {
	termMu.Lock()
	ts := termSession[id]
	termMu.Unlock()
	if ts == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "[term] stop id=%s\n", id)
	ts.mu.Lock()
	if ts.cmd.Process != nil {
		_ = ts.cmd.Process.Kill()
	}
	_ = ts.ptmx.Close()
	ts.mu.Unlock()
}

// WorkDirInfo describes one terminal working directory (a repo-root or
// exercise workspace entry).
type WorkDirInfo struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Dir   string `json:"dir"`
}

// WorkDirs lists the terminal's openable working directories: the repository
// root plus every exercise's answers workspace.
func (a *App) WorkDirs() []*WorkDirInfo {
	out := []*WorkDirInfo{{Key: "", Title: "Repository root", Dir: a.root}}
	for _, mod := range a.cur.Modules {
		for _, ex := range mod.Exercises {
			work, err := a.eng.WorkDir(ex)
			if err != nil {
				continue
			}
			out = append(out, &WorkDirInfo{
				Key:   mod.ID + "/" + ex.ID,
				Title: mod.ID + " · " + ex.Title,
				Dir:   work,
			})
		}
	}
	return out
}

// killAllTerms stops every running session (called on app shutdown).
func killAllTerms() {
	termMu.Lock()
	ids := make([]string, 0, len(termSession))
	for id, ts := range termSession {
		ids = append(ids, id)
		if ts.cmd.Process != nil {
			_ = ts.cmd.Process.Kill()
		}
		_ = ts.ptmx.Close()
	}
	termMu.Unlock()
}