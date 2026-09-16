// Package readings enforces the Reading-Ladder Standard
// (docs/readings-standard.md) over the subjects/ tree.
//
// The lint guarantees structural invariants only; pedagogic judgment about
// which entry source is easiest first stays a human review (rule R5 in the
// standard). Rules implemented here:
//
//	R1 every subjects/<module>/ex*/subject.md contains a "## Readings" section
//	R2 a subject citing the Linux Programming Interface / TLPI also cites at
//	   least one anchor of the form §N.M "Chapter title"
//	R3 the Readings block opens on an entry source (URL, man page, chapter,
//	   marker, book title, RFC) — never on a bare section-level TLPI dive
//	R4 every precise `man [[sect]] name` citation resolves via `man -w`
package readings

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Issue is a single lint finding.
type Issue struct {
	File string
	Rule string
	Msg  string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s  [%s] %s", i.File, i.Rule, i.Msg)
}

// entryRE matches a cold-start "entry source" token on the opening line of a
// Readings block: a URL, a man citation, a bold/keyword marker, an italic book
// title, an RFC number, or a chapter-level citation ("Chapter N", any book).
var entryRE = regexp.MustCompile(`https?://|\bman\s|\*\*|\*[^*\n]+\*|\bRFC\s+\d+|(?:^|[^A-Za-z])[Cc]hapter\s+\d+|(?i)reading ladder|glossary|cold-start`)

// targetRE matches precise man citations: "man 2 fork", "man 7 signal-safety",
// "man bash". Named capture "name" yields the page to resolve.
var targetRE = regexp.MustCompile(`\bman\s+(?:(?:2|3|1|7|8|4|5|6|9)(?:\.[a-z0-9]+)?\s+)?([a-zA-Z0-9_\-]+)`)

// tlpiMarkRE detects any reference to the primary C/OS canon.
var tlpiMarkRE = regexp.MustCompile(`(?i)linux programming interface|TLPI|\bTLPI\b`)

// tlpiAnchorRE accepts any numeric pin for a TLPI/DDIA-style citation:
// "§N", "§N.M", "§N.M–N.M", or "Chapter N".
var tlpiAnchorRE = regexp.MustCompile(`(?:§\d+(?:\.\d+)*(?:[–—.,-]\d+(?:\.\d+)*)?|[Cc]hapter\s+\d+)`)

// blockHeadingRE is the markdown heading that starts a Readings block.
var blockHeadingRE = regexp.MustCompile(`^#+\s+Readings\b`)

var stopwords = map[string]bool{
	"pages": true, "page": true, "manpage": true, "interfaces": true,
}

// Check lints the subjects/ tree rooted at root and reports all issues.
func Check(root string) ([]Issue, error) {
	r := &runner{root: root, manCache: map[string]bool{}}
	return r.run()
}

type runner struct {
	root     string
	manCache map[string]bool
}

func (r *runner) run() ([]Issue, error) {
	var issues []Issue

	subjects, err := filepath.Glob(filepath.Join(r.root, "subjects", "*", "ex*", "subject.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(subjects)

	// Collect every man target across all subjects so we resolve each page once.
	var manTargets []string
	for _, f := range subjects {
		txt, err := readFile(f)
		if err != nil {
			return nil, err
		}
		for _, m := range targetRE.FindAllStringSubmatch(txt, -1) {
			if len(m) > 1 && !stopwords[m[1]] {
				manTargets = append(manTargets, m[1])
			}
		}
	}
	sort.Strings(manTargets)
	manTargets = uniq(manTargets)

	manOK := map[string]bool{}
	for _, name := range manTargets {
		manOK[name] = r.manExists(name)
	}

	for _, f := range subjects {
		txt, err := readFile(f)
		if err != nil {
			return nil, err
		}
		issues = append(issues, r.lintSubject(f, txt, manOK)...)
	}
	return issues, nil
}

func (r *runner) lintSubject(f, txt string, manOK map[string]bool) []Issue {
	var out []Issue
	rel := strings.TrimPrefix(f, r.root+string(filepath.Separator))

	// R1 — the block must exist.
	block := readingsBlock(txt)
	if block == "" {
		out = append(out, Issue{rel, "R1", "subject has no ## Readings block"})
		return out
	}

	// R2 — TLPI cites must carry a numeric pin (§N / Chapter N).
	if tlpiMarkRE.MatchString(txt) && !tlpiAnchorRE.MatchString(block) {
		out = append(out, Issue{rel, "R2", "TLPI cited without a numeric pin (§N / Chapter N) in ## Readings"})
	}

	// R3 — the block must open on an entry source, not a bare section dive.
	first := firstContentLine(block)
	if first == "" {
		out = append(out, Issue{rel, "R3", "## Readings block is empty"})
	} else if !entryRE.MatchString(first) {
		out = append(out, Issue{rel, "R3", "block opens on a non-entry line (want URL, man page, chapter, ** marker, RFC, or book title): "+trunc(first, 90)})
	}

	// R4 — every precise man citation must resolve.
	for _, m := range targetRE.FindAllStringSubmatch(txt, -1) {
		if len(m) < 2 || stopwords[m[1]] {
			continue
		}
		if !manOK[m[1]] {
			out = append(out, Issue{rel, "R4", fmt.Sprintf("man(%s) does not resolve on this system", m[1])})
		}
	}
	return out
}

// readingsBlock returns the text of the ## Readings section (heading excluded,
// running to the next heading or EOF), or "" if the section is absent.
func readingsBlock(txt string) string {
	lines := strings.Split(txt, "\n")
	for i, line := range lines {
		if blockHeadingRE.MatchString(line) {
			var out []string
			for _, rest := range lines[i+1:] {
				if regexp.MustCompile(`^#{1,3}\s`).MatchString(rest) {
					break
				}
				out = append(out, rest)
			}
			return strings.Join(out, "\n")
		}
	}
	return ""
}

// firstContentLine returns the first non-blank line of a block that begins a
// structural element (bullet, ordered item, or heading).
func firstContentLine(block string) string {
	re := regexp.MustCompile(`(?m)^\s*(?:[-*+]\s+|#+\s+|\d+[.)]\s+)`)
	for _, line := range strings.Split(block, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if re.MatchString(line) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

func (r *runner) manExists(name string) bool {
	if hit, ok := r.manCache[name]; ok {
		return hit
	}
	err := exec.Command("man", "-w", name).Run()
	ok := err == nil
	r.manCache[name] = ok
	return ok
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func uniq(in []string) []string {
	var out []string
	for i, s := range in {
		if i == 0 || s != in[i-1] {
			out = append(out, s)
		}
	}
	return out
}

func trunc(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	rs := []rune(s)
	if len(rs) <= n {
		return s
	}
	return string(rs[:n]) + "…"
}