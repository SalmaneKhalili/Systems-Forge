package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/creack/pty"
)

// TestTermShellResponder reproduces the app's exact terminal setup end to end:
// bash is launched interactively (it reads ~/.bashrc, which execs fish), and
// the always-on termResponder answers fish's terminal capability probes. Falls
// back to plain bash when fish is missing. The shell must reach a prompt and
// respond to typed input.
func TestTermShellResponder(t *testing.T) {
	argv := "/bin/bash"
	if _, err := exec.LookPath("fish"); err != nil {
		t.Log("no fish; testing plain interactive bash instead")
	} else {
		argv = "/bin/bash" // .bashrc execs fish when present
	}
	cmd := exec.Command(argv)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("pty.Start: %v", err)
	}
	_ = pty.Setsize(ptmx, &pty.Winsize{Rows: 24, Cols: 80})

	resp := &termResponder{}
	chunks := make(chan []byte, 256)
	go func() {
		readPty(ptmx, func(b []byte) {
			chunks <- b
			if reply := resp.feed(b); len(reply) > 0 {
				_, _ = ptmx.Write(reply)
			}
		})
		_ = ptmx.Close()
	}()

	got := []byte{}
	prompt := func() bool {
		return bytes.Contains(got, []byte("❯")) || bytes.Contains(got, []byte("> "))
	}
	deadline := time.After(8 * time.Second)
	for !prompt() {
		select {
		case b := <-chunks:
			got = append(got, b...)
		case <-deadline:
			t.Fatalf("no shell prompt; output so far (%d bytes): %q", len(got), got)
		}
	}
	t.Log("shell prompt reached")

	if _, err := ptmx.Write([]byte("echo pty-shell-ok\r")); err != nil {
		t.Fatalf("write: %v", err)
	}
	deadline = time.After(5 * time.Second)
	for !bytes.Contains(got, []byte("pty-shell-ok")) {
		select {
		case b := <-chunks:
			got = append(got, b...)
		case <-deadline:
			t.Fatalf("typed command produced no reply; output: %q", got)
		}
	}
	t.Log("shell responded to typed input")
}

// TestTermResponderAnswers sanity-checks the canned replies so protocol
// regressions are caught without launching a shell.
func TestTermResponderAnswers(t *testing.T) {
	cases := []struct {
		query []byte
		want  []byte
	}{
		{[]byte{0x1b, '[', '?', 'u'}, []byte{0x1b, '[', '?', '1', 'u'}},
		{[]byte{0x1b, '[', '0', 'c'}, []byte{0x1b, '[', '?', '6', '2', ';', '2', '2', ';', 'c'}},
		{[]byte{0x1b, ']', '1', '1', ';', '?', 0x1b, '\\'}, []byte{0x1b, ']', '1', '1', ';', 'r', 'g', 'b', ':', '0', 'f', '/', '0', 'f', '/', '1', '5', 0x1b, '\\'}},
	}
	for i, c := range cases {
		ans, ok := termAnswer(c.query)
		if !ok {
			t.Errorf("case %d: query %q not answered", i, c.query)
			continue
		}
		if !bytes.Equal(ans, c.want) {
			t.Errorf("case %d: answer = %q, want %q", i, ans, c.want)
		}
	}
	// An unrecognized query must be left untouched.
	if ans, ok := termAnswer([]byte{0x1b, '[', '2', '5', 'l'}); ok {
		t.Errorf("rmcup was answered unexpectedly: %q", ans)
	}
}