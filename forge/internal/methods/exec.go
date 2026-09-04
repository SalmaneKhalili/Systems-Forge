package methods

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge/internal/cur"
	"forge/internal/sandbox"
)

// sanitizerFlags maps configured sanitizers to compiler flags.
var sanitizerFlags = map[string]string{
	"address": "-fsanitize=address",
	"undef":   "-fsanitize=undefined",
	"thread":  "-fsanitize=thread",
	"leak":    "-fsanitize=leak",
}

// buildFlags assembles CFLAGS/LDFLAGS for a build method.
func buildFlags(ex *cur.Exercise) (cflags, ldflags string) {
	var cf, lf []string
	// 42 piscine discipline: warnings are errors; modern C11 dialect.
	cf = append(cf, "-std=gnu11", "-Wall", "-Wextra", "-Werror")
	set := map[string]bool{}
	for _, s := range ex.Sanitizers {
		set[s] = true
	}
	if ex.Sanitize {
		set["address"] = true
		set["undef"] = true
	}
	for name := range set {
		if f, ok := sanitizerFlags[name]; ok {
			cf = append(cf, f)
			lf = append(lf, f)
		}
	}
	return strings.Join(cf, " "), strings.Join(lf, " ")
}

// compareRun checks a command's output against an expected file.
// dir is where the command runs (the exercise dir or a build subdir).
func compareRun(ctx context.Context, c *Ctx, name string, cmd []string,
	stdin, expect, norm string, stderrOK bool, wantExit int, dir string) (*Part, string, error) {

	part := &Part{Name: name}
	o := c.Opts()
	o.Dir = dir
	o.Cmd = cmd
	o.Stdin = stdin
	res := sandbox.Exec(ctx, o)

	merged := strings.TrimSpace(res.Stdout)
	if res.Stderr != "" {
		merged += "\n[stderr]\n" + res.Stderr
	}
	if !res.ExitOK(wantExit) {
		if res.TimedOut {
			part.Detail = "timed out"
		} else {
			part.Detail = fmt.Sprintf("exit %d (want %d)", res.ExitCode, wantExit)
		}
		part.Pass = false
		return part, merged, nil
	}
	if !stderrOK && strings.TrimSpace(res.Stderr) != "" {
		part.Pass = false
		part.Detail = "stderr not empty: " + strings.TrimSpace(res.Stderr)
		return part, merged, nil
	}
	if expect != "" {
		want, err := c.ReadExerciseFile(expect)
		if err != nil {
			part.Pass = false
			part.Detail = "cannot read expect file: " + err.Error()
			return part, merged, nil
		}
		expN := normalize(string(want), norm)
		gotN := normalize(res.Stdout, norm)
		if expN != gotN {
			part.Pass = false
			part.Detail = diffSummary(expN, gotN)
			return part, merged, nil
		}
	}
	part.Pass = true
	part.Detail = "ok"
	return part, merged, nil
}

// runBuild compiles via the student's Makefile under strict flags and runs
// verification binaries.
func runBuild(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	spec := m.Build
	res := &Result{}

	target := spec.Make
	if target == "" {
		target = "all"
	}
	buildDir := c.Dir
	if spec.Sub != "" {
		buildDir = filepath.Join(c.Dir, spec.Sub)
	}
	cflags, ldflags := buildFlags(c.Ex)

	env := append([]string{}, c.BaseEnv...)
	env = append(env, "CFLAGS="+cflags, "LDFLAGS="+ldflags, "CPPFLAGS="+cflags)

	clean := spec.Clean
	if clean == "" {
		clean = "fclean"
	}
	// A missing fclean target must not silently break the build: the
	// student's Makefile defines targets, so run fclean only if present.
	// We detect it via `make -q <target> 2>/dev/null` style? Simpler:
	// attempt fclean; if it fails, fall through to plain build (pseudo
	// targets like fclean are conventionally expected though).
	fc := c.Opts()
	fc.Dir = buildDir
	fc.Cmd = []string{"make", clean}
	fc.Env = env
	fcRes := sandbox.Exec(ctx, fc)
	if fcRes.TimedOut {
		res.Parts = append(res.Parts, &Part{Name: "make " + clean, Pass: false, Detail: "timed out"})
		res.Detail = fcRes.Stderr + fcRes.Stdout
		return res, nil
	}

	bld := c.Opts()
	bld.Dir = buildDir
	bld.Cmd = []string{"make", target}
	bld.Env = env
	bRes := sandbox.Exec(ctx, bld)
	part := &Part{Name: "make " + target, Pass: bRes.Success()}
	if !part.Pass {
		if bRes.TimedOut {
			part.Detail = "build timed out"
		} else {
			part.Detail = "build failed (exit " + itoa(bRes.ExitCode) + ")"
		}
		res.Parts = append(res.Parts, part)
		res.Stderr = strings.TrimSpace(bRes.Stdout + "\n" + bRes.Stderr)
		return res, nil
	}
	part.Detail = "build ok"
	res.Parts = append(res.Parts, part)

	for _, run := range spec.Runs {
		wantExit := run.ExpectExit
		part, merged, err := compareRun(ctx, c, run.Name, run.Cmd, run.Input, run.Expect, run.Normalize, run.StderrOK, wantExit, buildDir)
		if err != nil {
			return nil, err
		}
		res.Parts = append(res.Parts, part)
		if !part.Pass {
			res.Stderr = strings.TrimSpace(merged)
			break
		}
	}
	return res, nil
}

// runStdout runs a program and diffs stdout against an expected file.
func runStdout(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	spec := m.Stdout
	res := &Result{}
	part, merged, err := compareRun(ctx, c, "stdout", spec.Cmd, spec.Input, spec.Expect, spec.Normalize, spec.StderrOK, 0, c.Dir)
	if err != nil {
		return nil, err
	}
	res.Parts = append(res.Parts, part)
	res.Stdout = merged
	return res, nil
}

// runArtifact asserts properties of produced files.
func runArtifact(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	spec := m.Artifact
	res := &Result{}
	path := c.Abs(spec.Path)

	exists := true
	if spec.Exists != nil {
		exists = *spec.Exists
	}
	_, err := os.Stat(path)
	if exists != (err == nil) {
		part := &Part{Name: spec.Path, Pass: false}
		if exists {
			part.Detail = "missing"
		} else {
			part.Detail = "exists but should not"
		}
		res.Parts = append(res.Parts, part)
		return res, nil
	}
	if !exists {
		res.Parts = append(res.Parts, &Part{Name: spec.Path, Pass: true})
		return res, nil
	}

	st, _ := os.Stat(path)
	if spec.MinSize > 0 && st.Size() < spec.MinSize {
		res.Parts = append(res.Parts, &Part{Name: spec.Path, Pass: false,
			Detail: fmt.Sprintf("size %d < min %d", st.Size(), spec.MinSize)})
		return res, nil
	}
	if spec.MaxSize > 0 && st.Size() > spec.MaxSize {
		res.Parts = append(res.Parts, &Part{Name: spec.Path, Pass: false,
			Detail: fmt.Sprintf("size %d > max %d", st.Size(), spec.MaxSize)})
		return res, nil
	}

	if spec.IsA != "" {
		o := c.Opts()
		o.Cmd = []string{"file", "-b", path}
		fr := sandbox.Exec(ctx, o)
		if !fr.Success() || !strings.HasPrefix(strings.TrimSpace(fr.Stdout), spec.IsA) {
			res.Parts = append(res.Parts, &Part{Name: spec.Path, Pass: false,
				Detail: fmt.Sprintf("`file` = %q (want prefix %q)", strings.TrimSpace(fr.Stdout), spec.IsA)})
			return res, nil
		}
	}
	if len(spec.Symbols) > 0 {
		o := c.Opts()
		o.Cmd = []string{"nm", path}
		nr := sandbox.Exec(ctx, o)
		names := symbolNames(nr.Stdout)
		for _, want := range spec.Symbols {
			if !names[want] {
				res.Parts = append(res.Parts, &Part{Name: "symbol " + want, Pass: false, Detail: "not found in symbol table"})
				return res, nil
			}
		}
	}
	if len(spec.Contains) > 0 {
		data, err := sandbox.ReadFile(path)
		if err != nil {
			res.Parts = append(res.Parts, &Part{Name: spec.Path, Pass: false, Detail: err.Error()})
			return res, nil
		}
		for _, want := range spec.Contains {
			if !bytes.Contains(data, []byte(want)) {
				res.Parts = append(res.Parts, &Part{Name: "contains " + want, Pass: false, Detail: "substring not found"})
				return res, nil
			}
		}
	}
	res.Parts = append(res.Parts, &Part{Name: spec.Path, Pass: true, Detail: "ok"})
	return res, nil
}

// symbolNames extracts defined symbol names from nm output.
func symbolNames(nmOut string) map[string]bool {
	set := map[string]bool{}
	for _, line := range strings.Split(nmOut, "\n") {
		f := strings.Fields(line)
		if len(f) >= 3 {
			set[f[2]] = true
		}
	}
	return set
}

// parseKV parses "key: value" or "key=value" lines.
func parseKV(s string) map[string]string {
	out := map[string]string{}
	for _, ln := range strings.Split(s, "\n") {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		var k, v string
		if i := strings.Index(ln, ":"); i >= 0 {
			k, v = strings.TrimSpace(ln[:i]), strings.TrimSpace(ln[i+1:])
		} else if i := strings.Index(ln, "="); i >= 0 {
			k, v = strings.TrimSpace(ln[:i]), strings.TrimSpace(ln[i+1:])
		} else {
			continue
		}
		out[strings.ToLower(k)] = v
	}
	return out
}

// runReport parses structured output and compares to a reference file.
func runReport(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	spec := m.Report
	res := &Result{}
	refData, err := c.ReadExerciseFile(spec.Reference)
	if err != nil {
		res.Parts = append(res.Parts, &Part{Name: spec.Reference, Pass: false, Detail: "missing reference: " + err.Error()})
		return res, nil
	}
	ref := parseKV(string(refData))

	o := c.Opts()
	o.Cmd = spec.Cmd
	rr := sandbox.Exec(ctx, o)
	res.Stdout = rr.Stdout
	res.Stderr = rr.Stderr
	if !rr.Success() {
		res.Parts = append(res.Parts, &Part{Name: spec.Cmd[0], Pass: false,
			Detail: fmt.Sprintf("exit %d: %s", rr.ExitCode, strings.TrimSpace(rr.Stderr))})
		return res, nil
	}
	got := parseKV(rr.Stdout)

	referenced := map[string]bool{}
	for key, want := range ref {
		referenced[key] = true
		gotVal, found := got[key]
		if !found {
			res.Parts = append(res.Parts, &Part{Name: "field " + key, Pass: false, Detail: "absent from output"})
			continue
		}
		if tol, ok := spec.Tolerate[key]; ok {
			if !withinTolerance(gotVal, want, tol) {
				res.Parts = append(res.Parts, &Part{Name: "field " + key, Pass: false,
					Detail: fmt.Sprintf("%s vs %s (tol %.2f)", gotVal, want, tol)})
			} else {
				res.Parts = append(res.Parts, &Part{Name: "field " + key, Pass: true})
			}
			continue
		}
		if strings.TrimSpace(gotVal) != strings.TrimSpace(want) {
			res.Parts = append(res.Parts, &Part{Name: "field " + key, Pass: false,
				Detail: fmt.Sprintf("got %q want %q", gotVal, want)})
			continue
		}
		res.Parts = append(res.Parts, &Part{Name: "field " + key, Pass: true})
	}
	return res, nil
}

func withinTolerance(got, want string, tol float64) bool {
	var g, w float64
	if _, err := fmt.Sscanf(got, "%f", &g); err != nil {
		return strings.TrimSpace(got) == strings.TrimSpace(want)
	}
	if _, err := fmt.Sscanf(want, "%f", &w); err != nil {
		return false
	}
	if w == 0 {
		return g == 0
	}
	diff := g - w
	if diff < 0 {
		diff = -diff
	}
	return diff/absf(w) <= tol
}

func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
