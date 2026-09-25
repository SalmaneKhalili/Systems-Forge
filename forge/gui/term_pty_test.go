package main

import (
	"bytes"
	"os/exec"
	"testing"
	"time"

	"github.com/creack/pty"
)

// TestTermPtyReadLoop proves the terminal's output core: a real pty-backed
// process has its bytes forwarded by readPty in chunks.
func TestTermPtyReadLoop(t *testing.T) {
	cmd := exec.Command("sh", "-c", "printf 'pty-echo-ok'; exit 0")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("pty.Start: %v", err)
	}
	chunks := make(chan []byte, 8)
	readPty(ptmx, func(b []byte) { chunks <- b })
	if err := cmd.Wait(); err != nil {
		t.Fatalf("cmd.Wait: %v", err)
	}

	var got []byte
	for {
		select {
		case b := <-chunks:
			got = append(got, b...)
		case <-time.After(300 * time.Millisecond):
			goto drained
		}
	}
drained:
	if !bytes.Contains(got, []byte("pty-echo-ok")) {
		t.Fatalf("pty output = %q, want it to contain pty-echo-ok", got)
	}
}

// TestTermInteractiveWrites proves keystrokes written into the pty reach the
// process (the reverse direction of the terminal pipe). The child reads a
// fixed number of bytes then responds, so the pty closes without us needing
// to feed an EOF; written keys carry a trailing newline so the pty's canonical
// line discipline delivers them.
func TestTermInteractiveWrites(t *testing.T) {
	keys := "some-keys\n" // 10 bytes
	cmd := exec.Command("sh", "-c", "head -c 10 >/dev/null; printf 'read-%s' ok")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("pty.Start: %v", err)
	}
	chunks := make(chan []byte, 8)
	done := make(chan struct{})
	go func() {
		readPty(ptmx, func(b []byte) { chunks <- b })
		close(done)
		_ = ptmx.Close()
	}()

	if _, err := ptmx.Write([]byte(keys)); err != nil {
		t.Fatalf("ptmx.Write: %v", err)
	}
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("pty did not close after the child exited")
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("cmd.Wait: %v", err)
	}

	var got []byte
	for {
		select {
		case b := <-chunks:
			got = append(got, b...)
		case <-time.After(200 * time.Millisecond):
			goto drained
		}
	}
drained:
	if !bytes.Contains(got, []byte("read-ok")) {
		t.Fatalf("pty output = %q, want it to contain read-ok", got)
	}
}