// Package sandbox runs commands under resource limits, timeouts, and fresh
// process groups so a check can never hang or leak processes.
package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Options controls a sandboxed command invocation.
type Options struct {
	Cmd     []string
	Dir     string
	Env     []string // extra environment, e.g. "CFLAGS=-Wall -Werror"
	Stdin   string
	Timeout time.Duration
	AS      int64 // address space limit, bytes (0 = inherited)
	FSize   int64 // file size limit, bytes
	NProc   int64 // process/thread limit
	Nofile  int64 // open-file limit

	// Background launches into its own session and returns immediately.
	// The caller must call Handle.Stop; otherwise a process group leaks.
	Background bool
}

// Handle wraps a running background process group.
type Handle struct {
	pgid     int
	done     chan struct{}
	once     sync.Once
	exitCode int
}

// Pid returns the process group id of the launched command (all children
// share this group).
func (h *Handle) Pid() int {
	if h == nil {
		return -1
	}
	return h.pgid
}

// ExitCode returns the process exit code after it has exited.
// Returns -1 if the handle is nil or the process hasn't exited yet.
func (h *Handle) ExitCode() int {
	if h == nil {
		return -1
	}
	<-h.done
	return h.exitCode
}

// Wait blocks until the process exits and returns the exit code.
// Returns -1 and false if the timeout elapses before exit.
func (h *Handle) Wait(timeout time.Duration) (int, bool) {
	if h == nil {
		return -1, false
	}
	select {
	case <-h.done:
		return h.exitCode, true
	case <-time.After(timeout):
		return -1, false
	}
}

// Stop terminates the process group. When graceful is true it sends SIGTERM
// first and escalates to SIGKILL after a short grace period.
func (h *Handle) Stop(graceful bool) {
	if h == nil {
		return
	}
	h.once.Do(func() {
		_ = syscall.Kill(-h.pgid, syscall.SIGTERM)
		if graceful {
			select {
			case <-h.done:
				return
			case <-time.After(2 * time.Second):
			}
		}
		_ = syscall.Kill(-h.pgid, syscall.SIGKILL)
		<-h.done
	})
}

// Result of a sandboxed execution.
type Result struct {
	ExitCode int
	TimedOut bool
	Stdout   string
	Stderr   string
	Duration time.Duration
}

// Success reports a clean run with exit code 0.
func (r Result) Success() bool {
	return !r.TimedOut && r.ExitCode == 0
}

// ExitOK reports whether the command finished (not timed out) and exited
// with the wanted code.
func (r Result) ExitOK(want int) bool {
	return !r.TimedOut && r.ExitCode == want
}

var baseOnce sync.Once

// envFor returns a deterministic environment overlay for a command.
func envFor(extra []string) []string {
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"TERM=dumb",
		"LC_ALL=C",
		"TZ=UTC",
	}
	return append(env, extra...)
}

// execPrefix builds the wrapper that applies ulimits then execs the command
// through setsid so every launched process lands in a fresh session whose
// group leader pid we know.
func execPrefix(o Options) (argv []string) {
	var pre []string
	if o.AS > 0 {
		pre = append(pre, "ulimit -v "+strconv.FormatInt(o.AS, 10)+" 2>/dev/null;")
	}
	if o.FSize > 0 {
		pre = append(pre, "ulimit -f "+strconv.FormatInt(o.FSize, 10)+" 2>/dev/null;")
	}
	if o.NProc > 0 {
		pre = append(pre, "ulimit -u "+strconv.FormatInt(o.NProc, 10)+" 2>/dev/null;")
	}
	if o.Nofile > 0 {
		pre = append(pre, "ulimit -n "+strconv.FormatInt(o.Nofile, 10)+" 2>/dev/null;")
	}
	script := strings.Join(pre, " ") + " exec setsid \"$@\""
	argv = append(argv, "bash", "-c", script, "--")
	argv = append(argv, o.Cmd...)
	return argv
}

func sanitizeTimeout(t time.Duration) time.Duration {
	if t <= 0 {
		return 30 * time.Second
	}
	return t
}

// Exec runs a command to completion under the sandbox.
func Exec(ctx context.Context, o Options) Result {
	o.Timeout = sanitizeTimeout(o.Timeout)

	var stdout, stderr bytes.Buffer
	argv := execPrefix(o)
	cmd := exec.Command("bash", argv[1:]...)
	cmd.Dir = o.Dir
	cmd.Env = envFor(o.Env)
	cmd.Stdin = strings.NewReader(o.Stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	if err := cmd.Start(); err != nil {
		return Result{ExitCode: -1, Stderr: "spawn error: " + err.Error()}
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	timedOut := false
	var timer <-chan time.Time
	if o.Timeout > 0 {
		timer = time.After(o.Timeout)
	}
	select {
	case <-ctx.Done():
		timedOut = true
	case <-timer:
		timedOut = true
	case <-done:
	}
	if timedOut {
		_ = cmd.Process.Kill()
		<-done
		// Make sure the whole session tree is gone.
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	res := Result{
		ExitCode: -1,
		TimedOut: timedOut,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: time.Since(start),
	}
	if cmd.ProcessState != nil {
		res.ExitCode = cmd.ProcessState.ExitCode()
	}
	return res
}

// Background launches a command into its own session and returns a handle.
// All stdin/stdout are /dev/null. Failure to Stop leaks a process group.
func Background(o Options) (*Handle, error) {
	argv := execPrefix(o)
	cmd := exec.Command("bash", argv[1:]...)
	cmd.Dir = o.Dir
	cmd.Env = envFor(o.Env)
	cmd.Stdin = nil // /dev/null
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	h := &Handle{pgid: cmd.Process.Pid, done: make(chan struct{})}
	go func() {
		_ = cmd.Wait()
		if cmd.ProcessState != nil {
			h.exitCode = cmd.ProcessState.ExitCode()
		}
		close(h.done)
	}()
	return h, nil
}

// ReadFile is a tiny helper for tool code that reads a (possibly large)
// file without losing errors.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(path))
}

// FormatBytes renders a size for diagnostics.
func FormatBytes(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1fGiB", float64(n)/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1fMiB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1fKiB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%dB", n)
	}
}
