package methods

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"forge/internal/cur"
	"forge/internal/sandbox"
)

// runProcess spawns a program, probes it while running, then stops it.
// The spawned process receives TARGET_PID (its own pid / process group).
func runProcess(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	spec := m.Process
	res := &Result{}

	env := append([]string{}, c.BaseEnv...)
	for k, v := range spec.TargetEnv {
		env = append(env, k+"="+v)
	}

	h, err := sandbox.Background(sandbox.Options{Cmd: spec.Start, Dir: c.Dir, Env: env})
	if err != nil {
		res.Parts = append(res.Parts, &Part{Name: "start", Pass: false, Detail: err.Error()})
		return res, nil
	}
	// TARGET_PID is the process group id == child pid.
	env = append(env, "TARGET_PID="+strconv.Itoa(h.Pid()))
	defer h.Stop(spec.Kill != "kill")

	select {
	case <-time.After(time.Duration(spec.WaitMs) * time.Millisecond):
	case <-ctx.Done():
		res.Parts = append(res.Parts, &Part{Name: "start", Pass: false, Detail: "cancelled during startup"})
		return res, nil
	}

	if len(spec.Probes) == 0 {
		res.Parts = append(res.Parts, &Part{Name: "started", Pass: true, Detail: "process spawned and alive"})
		return res, nil
	}

	for _, probe := range spec.Probes {
		penv := append([]string{}, c.BaseEnv...)
		for k, v := range probe.Env {
			penv = append(penv, k+"="+v)
		}
		penv = append(penv, "TARGET_PID="+strconv.Itoa(h.Pid()))
		o := c.Opts()
		o.Cmd = probe.Exec
		o.Env = penv
		pr := sandbox.Exec(ctx, o)
		part := &Part{Name: strings.Join(probe.Exec, " ")}
		if !pr.ExitOK(probe.ExpectExit) {
			part.Pass = false
			if pr.TimedOut {
				part.Detail = "probe timed out"
			} else {
				part.Detail = fmt.Sprintf("probe exit %d (want %d)", pr.ExitCode, probe.ExpectExit)
			}
			res.Parts = append(res.Parts, part)
			res.Stderr = strings.TrimSpace(pr.Stdout + "\n" + pr.Stderr)
			return res, nil
		}
		if probe.ExpectOut != "" && !strings.Contains(pr.Stdout, probe.ExpectOut) {
			part.Pass = false
			part.Detail = fmt.Sprintf("probe output missing %q:\n%s", probe.ExpectOut, strings.TrimSpace(pr.Stdout))
			res.Parts = append(res.Parts, part)
			res.Stderr = strings.TrimSpace(pr.Stdout + "\n" + pr.Stderr)
			return res, nil
		}
		part.Pass = true
		part.Detail = "ok"
		res.Parts = append(res.Parts, part)
	}
	return res, nil
}

// runNet drives a scripted protocol against a spawned server.
func runNet(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	spec := m.Net
	res := &Result{}

	host := spec.Host
	if host == "" {
		host = "127.0.0.1"
	}

	env := append([]string{}, c.BaseEnv...)
	for k, v := range spec.StartEnv {
		env = append(env, k+"="+v)
	}
	env = append(env, "TARGETPORT="+strconv.Itoa(spec.Port), "TARGETHOST="+host)

	h, err := sandbox.Background(sandbox.Options{Cmd: spec.Start, Dir: c.Dir, Env: env})
	if err != nil {
		res.Parts = append(res.Parts, &Part{Name: "start", Pass: false, Detail: err.Error()})
		return res, nil
	}
	defer h.Stop(true)

	timeout := spec.TimeoutMs
	if timeout <= 0 {
		timeout = 5000
	}
	deadline := time.Now().Add(time.Duration(spec.WaitMs)*time.Millisecond +
		time.Duration(timeout+1000)*time.Millisecond)

	// Retry dialing the first connection until the server is up or cancelled.
	conn, err := dialRetry(ctx, host, spec.Port, deadline)
	if err != nil {
		cause := err
		if ctx.Err() != nil {
			cause = errors.New("cancelled")
		}
		res.Parts = append(res.Parts, &Part{Name: "start", Pass: false,
			Detail: "server never accepted connections: " + cause.Error()})
		return res, nil
	}
	defer conn.Close()

	if len(spec.Connections) == 0 {
		_ = conn.SetDeadline(time.Now().Add(time.Duration(timeout+1000) * time.Millisecond))
		runNetSteps(conn, spec.Steps, "step ", timeout, res)
		return res, nil
	}

	_ = conn.SetDeadline(time.Now().Add(time.Duration(timeout+1000) * time.Millisecond))
	if !runNetSteps(conn, spec.Connections[0].Steps, "conn 1 step ", timeout, res) {
		return res, nil
	}
	for i := 1; i < len(spec.Connections); i++ {
		_ = conn.Close()
		conn, err = dialRetry(ctx, host, spec.Port, time.Now().Add(2*time.Second))
		if err != nil {
			res.Parts = append(res.Parts, &Part{Name: fmt.Sprintf("conn %d connect", i+1), Pass: false,
				Detail: "server stopped accepting: " + err.Error()})
			return res, nil
		}
		_ = conn.SetDeadline(time.Now().Add(time.Duration(timeout+1000) * time.Millisecond))
		if !runNetSteps(conn, spec.Connections[i].Steps, fmt.Sprintf("conn %d step ", i+1), timeout, res) {
			return res, nil
		}
	}
	return res, nil
}

// dialRetry keeps dialing until it succeeds or the deadline passes.
func dialRetry(ctx context.Context, host string, port int, deadline time.Time) (net.Conn, error) {
	var err error
	for {
		var conn net.Conn
		conn, err = net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), 300*time.Millisecond)
		if err == nil {
			return conn, nil
		}
		if time.Now().After(deadline) {
			return nil, err
		}
		select {
		case <-time.After(200 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// runNetSteps plays one conn's send/expect script, appending parts to res.
// Returns false on the first failed step.
func runNetSteps(conn net.Conn, steps []*cur.NetStep, prefix string, timeout int, res *Result) bool {
	for i, step := range steps {
		send := decodeBytes(step.Send)
		var sendErr error
		if step.SendHex != "" {
			send, sendErr = hex.DecodeString(step.SendHex)
			if sendErr != nil {
				res.Parts = append(res.Parts, &Part{Name: prefix + itoa(i+1) + " send", Pass: false,
					Detail: "bad send_hex: " + sendErr.Error()})
				return false
			}
		}
		if len(send) > 0 {
			if _, sendErr = conn.Write(send); sendErr != nil {
				res.Parts = append(res.Parts, &Part{Name: prefix + itoa(i+1) + " send", Pass: false,
					Detail: sendErr.Error()})
				return false
			}
		}
		if step.SleepMs > 0 {
			time.Sleep(time.Duration(step.SleepMs) * time.Millisecond)
		}

		var want []byte
		if step.ExpectHex != "" {
			want, sendErr = hex.DecodeString(step.ExpectHex)
			if sendErr != nil {
				res.Parts = append(res.Parts, &Part{Name: prefix + itoa(i+1) + " expect", Pass: false,
					Detail: "bad expect_hex: " + sendErr.Error()})
				return false
			}
		} else if step.Expect != "" {
			want = decodeBytes(step.Expect)
		}

		part := &Part{Name: prefix + itoa(i+1)}
		got, ok := netTransaction(conn, timeout, want, step.Contains)
		if step.Contains != "" {
			if ok {
				part.Pass = true
				part.Detail = "contains ok"
			} else {
				part.Pass = false
				part.Detail = fmt.Sprintf("substring %q not received (got %d bytes)", step.Contains, len(got))
			}
		} else if len(want) > 0 {
			if ok {
				part.Pass = true
				part.Detail = "ok"
			} else {
				part.Pass = false
				part.Detail = fmt.Sprintf("expected %d bytes, got %d:\n  want %q\n  got  %q",
					len(want), len(got), string(want), string(got))
			}
		} else {
			part.Pass = len(got) > 0
			if !part.Pass {
				part.Detail = "no data received"
			}
		}
		res.Parts = append(res.Parts, part)
		if !part.Pass {
			return false
		}
	}
	return true
}

// netTransaction reads from conn until the expectation is satisfied or the
// overall timeout elapses, then reports whether it matched. want == nil and
// contains == "" means "any data".
func netTransaction(conn net.Conn, timeout int, want []byte, contains string) ([]byte, bool) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	var got []byte
	buf := make([]byte, 4096)
	for {
		if contains != "" && ContainsBytes(got, []byte(contains)) {
			return got, true
		}
		if len(want) > 0 && bytesEqual(got, want) {
			return got, true
		}
		if contains == "" && len(want) == 0 && len(got) > 0 {
			return got, true
		}
		if time.Now().After(deadline) {
			return got, false
		}
		_ = conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		n, rerr := conn.Read(buf)
		if n > 0 {
			got = append(got, buf[:n]...)
			continue
		}
		if rerr != nil {
			if netErr, ok := rerr.(net.Error); ok && netErr.Timeout() {
				if len(got) > 0 && contains == "" && len(want) == 0 {
					return got, true
				}
				continue
			}
			return got, false
		}
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
