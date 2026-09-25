package main

import (
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
		return err
	}
	termMu.Lock()
	termSession[id] = &TermSession{cmd: cmd, ptmx: ptmx}
	termMu.Unlock()

	go a.pump(id, ptmx)
	return nil
}

// pump forwards pty output to the frontend and cleans up on exit.
func (a *App) pump(id string, ptmx *os.File) {
	buf := make([]byte, 32*1024)
	for {
		n, err := ptmx.Read(buf)
		if n > 0 {
			runtime.EventsEmit(a.ctx, "term:out:"+id, base64.StdEncoding.EncodeToString(buf[:n]))
		}
		if err != nil {
			break
		}
	}
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
	runtime.EventsEmit(a.ctx, "term:exit:"+id, true)
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