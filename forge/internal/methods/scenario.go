package methods

import (
	"bufio"
	"context"
	"strings"
	"time"

	"forge/internal/cur"
	"forge/internal/sandbox"
)

// runScenario / runFault execute a driver script that boots nodes, injects
// faults, and asserts outcomes.
//
// Driver contract:
//   - exit 0            -> pass
//   - nonzero exit      -> fail
//   - a line "FORGE_RESULT: pass|fail" overrides the exit code (useful when
//     the driver intentionally leaves torn-down state)
//   - lines "FORGE_META: key=value" are captured into the result detail and
//     surfaced in reports and the TUI
func runScenario(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	var driver []string
	var env map[string]string
	var timeout int
	label := m.Type
	if m.Scenario != nil {
		driver, env, timeout = m.Scenario.Driver, m.Scenario.Env, m.Scenario.TimeoutS
	} else if m.Fault != nil {
		driver, env, timeout = m.Fault.Driver, m.Fault.Env, m.Fault.TimeoutS
	} else {
		return &Result{Parts: []*Part{{Name: label, Pass: false, Detail: "no spec"}}}, nil
	}

	res := &Result{}
	o := c.Opts()
	if timeout > 0 {
		o.Timeout = time.Duration(timeout) * time.Second
	}
	o.Cmd = driver
	for k, v := range env {
		o.Env = append(o.Env, k+"="+v)
	}
	// Expose repo root so drivers can reference the switch and tools.
	o.Env = append(o.Env, "FORGE_ROOT="+c.Root, "FORGE_EXERCISE_DIR="+c.Dir)
	dr := sandbox.Exec(ctx, o)
	res.Stdout = dr.Stdout
	res.Stderr = dr.Stderr

	meta := parseDriverOutput(dr.Stdout)
	pass := dr.Success()
	if r, ok := meta["result"]; ok {
		pass = strings.EqualFold(strings.TrimSpace(r), "pass")
	}
	if dr.TimedOut {
		pass = false
		res.Detail = "driver timed out"
	} else if !pass {
		tail := strings.TrimSpace(dr.Stdout)
		if tail == "" {
			tail = strings.TrimSpace(dr.Stderr)
		}
		if len(tail) > 400 {
			tail = tail[len(tail)-400:]
		}
		res.Detail = "driver failed: " + tail
	}
	res.Parts = append(res.Parts, &Part{Name: label, Pass: pass, Detail: metaDetail(meta)})
	return res, nil
}

func parseDriverOutput(stdout string) map[string]string {
	out := map[string]string{}
	sc := bufio.NewScanner(strings.NewReader(stdout))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "FORGE_RESULT:"):
			out["result"] = strings.TrimSpace(strings.TrimPrefix(line, "FORGE_RESULT:"))
		case strings.HasPrefix(line, "FORGE_META:"):
			kv := strings.TrimPrefix(line, "FORGE_META:")
			k, v, ok := strings.Cut(kv, "=")
			if ok {
				out["meta:"+strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
		}
	}
	return out
}

func metaDetail(meta map[string]string) string {
	if _, ok := meta["result"]; ok {
		return "result: " + meta["result"]
	}
	return "driver ok"
}
