package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge/internal/check"
	"forge/internal/cur"
	"forge/internal/methods"
	"forge/internal/sandbox"
	"forge/internal/store"
	"forge/internal/tui"
)

func usage() {
	fmt.Print(`forge — the systems-forge platform

Usage:
  forge                 launch the TUI
  forge init            create your workspace under answers/ (mirrors subjects/)
  forge list [module]   show modules and exercises (optionally one module)
  forge show <id>       print an exercise specification (subject.md)
  forge check <id>      run the grader for one exercise
  forge score           show passing progress per module
  forge selftest        validate the grader against built-in fixtures
  forge help            this text

Exercise ids look like "M2-ex03" (spelled by their directory basename).
`)
}

func run(args []string) error {
	verb := ""
	rest := args
	if len(args) > 0 {
		verb = args[0]
		rest = args[1:]
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	switch verb {
	case "", "tui":
		return tui.Run(root)
	case "init":
		return cmdInit(root)
	case "list":
		return cmdList(root, rest)
	case "show":
		if len(rest) != 1 {
			return fmt.Errorf("usage: forge show <id>")
		}
		return cmdShow(root, rest[0])
	case "check":
		if len(rest) != 1 {
			return fmt.Errorf("usage: forge check <id>")
		}
		return cmdCheck(root, rest[0])
	case "score":
		return cmdScore(root)
	case "selftest":
		return cmdSelftest(root)
	case "help", "--help", "-h":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown command %q (try: forge help)", verb)
	}
}

func cmdInit(root string) error {
	eng, err := check.New(root)
	if err != nil {
		return err
	}
	defer eng.Close()
	curriculum, err := cur.Load(eng.Subjects)
	if err != nil {
		return err
	}
	if len(curriculum.Modules) == 0 {
		fmt.Println("no modules in subjects/ yet")
		return nil
	}
	count := 0
	for _, mod := range curriculum.Modules {
		for _, ex := range mod.Exercises {
			if _, err := eng.WorkDir(ex); err != nil {
				return err
			}
			count++
		}
	}
	fmt.Printf("workspace ready under answers/ (%d exercises mirrored)\n", count)
	return nil
}

func cmdList(root string, args []string) error {
	eng, err := check.New(root)
	if err != nil {
		return err
	}
	defer eng.Close()
	curriculum, err := cur.Load(eng.Subjects)
	if err != nil {
		return err
	}
	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}
	for _, mod := range curriculum.Modules {
		if filter != "" && !strings.Contains(mod.ID, filter) {
			continue
		}
		fmt.Printf("\n%s\n", mod.ID)
		if mod.Title != "" {
			fmt.Printf("  %s\n", mod.Title)
		}
		passed, total, err := eng.Store.ModuleProgress(exerciseIDs(mod))
		if err != nil {
			return err
		}
		fmt.Printf("  [%d/%d]\n", passed, total)
		for _, ex := range mod.Exercises {
			status := " "
			last, _ := eng.Store.LastRun(ex.ID)
			if last != nil {
				if last.Status == store.StatusPass {
					status = "P"
				} else {
					status = "F"
				}
			}
			fmt.Printf("   [%s] %-24s %s\n", status, ex.ID, ex.Title)
		}
	}
	return nil
}

func cmdShow(root, id string) error {
	eng, err := check.New(root)
	if err != nil {
		return err
	}
	defer eng.Close()
	_, ex, err := eng.FindExercise(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "did you mean: forge list")
		return nil
	}
	subject := filepath.Join(ex.Dir, "subject.md")
	data, err := sandbox.ReadFile(subject)
	if err != nil {
		return fmt.Errorf("exercise has no subject.md: %w", err)
	}
	fmt.Println(strings.TrimRight(string(data), "\n"))
	return nil
}

func cmdCheck(root, id string) error {
	eng, err := check.New(root)
	if err != nil {
		return err
	}
	defer eng.Close()
	_, ex, err := eng.FindExercise(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "did you mean: forge list")
		return nil
	}
	report, err := eng.Check(context.Background(), ex)
	if err != nil {
		return err
	}
	renderReport(report)
	if !report.Pass {
		os.Exit(1)
	}
	return nil
}

func cmdScore(root string) error {
	eng, err := check.New(root)
	if err != nil {
		return err
	}
	defer eng.Close()
	curriculum, err := cur.Load(eng.Subjects)
	if err != nil {
		return err
	}
	var all []string
	for _, mod := range curriculum.Modules {
		ids := exerciseIDs(mod)
		all = append(all, ids...)
		passed, total, err := eng.Store.ModuleProgress(ids)
		if err != nil {
			return err
		}
		bar := progressBar(passed, total)
		fmt.Printf("%-18s %3d/%-3d %s\n", mod.ID, passed, total, bar)
	}
	passed, attempted, total, err := eng.Store.OverallProgress(all)
	if err != nil {
		return err
	}
	fmt.Println(strings.Repeat("-", 44))
	fmt.Printf("%-18s %3d/%-3d   (%d attempted, %d remaining)\n", "TOTAL", passed, total, attempted, total-attempted)
	return nil
}

func exerciseIDs(mod *cur.Module) []string {
	var out []string
	for _, ex := range mod.Exercises {
		out = append(out, ex.ID)
	}
	return out
}

func progressBar(passed, total int) string {
	const width = 24
	if total == 0 {
		return ""
	}
	filled := passed * width / total
	return "[" + strings.Repeat("#", filled) + strings.Repeat(".", width-filled) + "]"
}

func renderReport(r *check.Report) {
	for _, m := range r.Methods {
		mark := "ok "
		if !m.Pass {
			mark = "FAIL"
		}
		fmt.Printf("  [%s] %s (%dms)\n", mark, m.Label, m.DurMs)
		for _, p := range m.Parts {
			pm := "ok "
			if !p.Pass {
				pm = "x"
			}
			fmt.Printf("       [%s] %s\n", pm, p.Name)
			if !p.Pass && p.Detail != "" {
				for _, line := range strings.Split(p.Detail, "\n") {
					fmt.Printf("            %s\n", line)
				}
			}
		}
		if !m.Pass && m.Output != "" {
			fmt.Printf("   --- output ---\n%s\n", m.Output)
		}
	}
	if r.Pass {
		fmt.Printf("\nPASS  %s in %v\n", r.Exercise, r.Duration.Round(0))
	} else {
		fmt.Printf("\nFAIL  %s in %v\n", r.Exercise, r.Duration.Round(0))
	}
}

// cmdSelftest validates the grader against tools/fixtures/<method>/<expect>/<case>.
func cmdSelftest(root string) error {
	fixturesRoot := filepath.Join(root, "tools", "fixtures")
	methodDirs, err := os.ReadDir(fixturesRoot)
	if err != nil {
		return fmt.Errorf("no fixtures found: %w", err)
	}
	failed := 0
	total := 0
	for _, md := range methodDirs {
		if !md.IsDir() {
			continue
		}
		expectationDirs, err := os.ReadDir(filepath.Join(fixturesRoot, md.Name()))
		if err != nil {
			return err
		}
		for _, ed := range expectationDirs {
			if !ed.IsDir() {
				continue
			}
			caseDirs, err := os.ReadDir(filepath.Join(fixturesRoot, md.Name(), ed.Name()))
			if err != nil {
				return err
			}
			for _, cd := range caseDirs {
				if !cd.IsDir() {
					continue
				}
				dir := filepath.Join(fixturesRoot, md.Name(), ed.Name(), cd.Name())
				ex, err := cur.LoadExercise(dir)
				if err != nil {
					return fmt.Errorf("fixture %s: %w", dir, err)
				}
				total++
				expectPass := ex.Expect == "pass"
				fx := &methods.Ctx{Ex: ex, Dir: dir, Root: root}
				allPass := true
				for _, m := range ex.Methods {
					res, err := methods.Dispatch(context.Background(), fx, m)
					if err != nil || !res.Pass {
						allPass = false
					}
				}
				ok := allPass == expectPass
				label := fmt.Sprintf("%s/%s/%s", md.Name(), ed.Name(), cd.Name())
				if ok {
					fmt.Printf("  ok  %-40s (score %s)\n", label, wanted(ex.Expect))
				} else {
					failed++
					got := "pass"
					if !allPass {
						got = "fail"
					}
					fmt.Printf("  BAD %-40s expected %s, grader scored %s\n", label, ex.Expect, got)
				}
			}
		}
	}
	fmt.Printf("\n%d fixtures, %d mismatches\n", total, failed)
	if failed > 0 {
		return fmt.Errorf("grader self-test failed")
	}
	return nil
}

func wanted(expect string) string {
	if expect == "pass" {
		return "PASS expected"
	}
	return "FAIL expected"
}
