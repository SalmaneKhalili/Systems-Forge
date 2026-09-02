package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"forge/internal/cur"
	"forge/internal/sandbox"
)

// Op is one operation in a history: a write (sets Val) or a read (returns
// Val) over a single logical key/register.
//
// Start/End are real-time interval bounds in the caller's timebase; the
// checker only uses their ordering (a completed-write edge).
type Op struct {
	Type  string  `json:"type"`
	Key   string  `json:"key"`
	Val   *string `json:"val"`
	Start int64   `json:"start"`
	End   int64   `json:"end"`

	idx     int
	succ    []*Op
	pending int // remaining unplaced predecessors
}

// History is the JSON document produced by scenario drivers.
type History struct {
	Ops []*Op `json:"ops"`
}

// linearizable implements a backtracking search for a total order of ops
// that (a) respects real-time order and (b) makes every read return the
// value of the last preceding write to its key.
//
// It returns linearizable=false with a witness when no such order exists,
// and (false, "") with err set when the history is too complex to decide
// within budget (the caller reports an indeterminate error, never a false
// pass/fail).
func linearizable(ops []*Op) (bool, string, error) {
	n := len(ops)
	if n == 0 {
		return true, "", nil
	}
	if n > maxLinOps {
		return false, "", fmt.Errorf("history too large (%d ops, max %d)", n, maxLinOps)
	}
	// Adopt the caller's Op structs (they may point into a shared slice).
	// Compute real-time successors: A must precede B if A finished before
	// B started.
	for i := 0; i < n; i++ {
		ops[i].idx = i
		ops[i].succ = nil
		ops[i].pending = 0
	}
	for i := 0; i < n; i++ {
		for s := 0; s < n; s++ {
			if i != s && ops[i].End < ops[s].Start {
				ops[i].succ = append(ops[i].succ, ops[s])
			}
		}
	}
	for i := 0; i < n; i++ {
		for p := 0; p < n; p++ {
			if i != p && ops[p].End < ops[i].Start {
				ops[i].pending++
			}
		}
	}

	// current[regKey] is the value of the last write placed so far.
	current := map[string]*string{}
	var valLog []struct {
		key string
		old *string
	}
	order := make([]int, 0, n)
	placed := make([]bool, n)
	iterations := 0

	var search func() bool
	search = func() bool {
		if iterations > searchBudget {
			return false
		}
		iterations++
		if len(order) == n {
			return true
		}
		// Candidates: ready ops (all predecessors placed). Reads whose
		// register value doesn't match are placeable only after the right
		// write lands, so treat them as non-candidates for now.
		var writes []int
		var reads []int
		for i := 0; i < n; i++ {
			if placed[i] || ops[i].pending != 0 {
				continue
			}
			if ops[i].Type == "write" {
				writes = append(writes, i)
			} else if sameVal(current[ops[i].Key], ops[i].Val) {
				reads = append(reads, i)
			}
		}
		if len(writes)+len(reads) == 0 {
			return false
		}
		for _, i := range append(writes, reads...) {
			op := ops[i]
			placed[i] = true
			if op.Type == "write" {
				valLog = append(valLog, struct {
					key string
					old *string
				}{op.Key, current[op.Key]})
				current[op.Key] = op.Val
			}
			for _, s := range op.succ {
				ops[s.idx].pending--
			}
			order = append(order, i)
			if search() {
				return true
			}
			order = order[:len(order)-1]
			for _, s := range op.succ {
				ops[s.idx].pending++
			}
			if op.Type == "write" {
				lv := valLog[len(valLog)-1]
				valLog = valLog[:len(valLog)-1]
				current[lv.key] = lv.old
			}
			placed[i] = false
		}
		return false
	}

	if search() {
		return true, "", nil
	}
	if iterations > searchBudget {
		return false, "", fmt.Errorf("search budget exceeded (%d ops)", n)
	}
	return false, witnessString(ops, order, placed), nil
}

// These are deliberately small for deterministic graders; scenario drivers
// that need larger verdicts sample histories into windows.
const (
	maxLinOps    = 220
	searchBudget = 2_000_000
)

func witnessString(ops []*Op, order []int, placed []bool) string {
	var b strings.Builder
	b.WriteString("no linearization exists; op order so far:\n")
	for i, idx := range order {
		o := ops[idx]
		fmt.Fprintf(&b, "  %3d. [%s] %s=%s\n", i, o.Type, o.Key, valOrNil(o.Val))
	}
	b.WriteString("remaining unplaceable ops:\n")
	for i, o := range ops {
		if !placed[i] {
			fmt.Fprintf(&b, "  - [%s] %s=%s (time %d..%d)\n", o.Type, o.Key, valOrNil(o.Val), o.Start, o.End)
		}
	}
	return b.String()
}

func valOrNil(v *string) string {
	if v == nil {
		return "⊥"
	}
	return *v
}

func sameVal(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// isWrittenHist helpers ---------------------------------------------------

// parseHistory reads a history document (either {"ops":[...] } or a bare
// JSON array).
func parseHistory(data []byte) (*History, error) {
	var h History
	if err := json.Unmarshal(data, &h); err == nil && h.Ops != nil {
		return &h, normalizeVals(&h)
	}
	var arr []*Op
	dec := json.NewDecoder(strings.NewReader(string(data)))
	if err := dec.Decode(&arr); err != nil {
		return nil, fmt.Errorf("history decode: %w", err)
	}
	h.Ops = arr
	return &h, normalizeVals(&h)
}

// normalizeVals coerces JSON scalars to strings and validates op types.
func normalizeVals(h *History) error {
	for i, op := range h.Ops {
		switch op.Type {
		case "read", "write":
		default:
			return fmt.Errorf("op %d: unknown type %q", i, op.Type)
		}
		if op.Key == "" {
			return fmt.Errorf("op %d: missing key", i)
		}
	}
	return nil
}

// runLincheck grades histories found by globbing spec.History.
func runLincheck(ctx context.Context, c *Ctx, m *cur.Method) (*Result, error) {
	spec := m.Lincheck
	res := &Result{}
	files, err := globHistoryFiles(c.Abs(spec.History))
	if err != nil {
		res.Parts = append(res.Parts, &Part{Name: "history", Pass: false, Detail: err.Error()})
		return res, nil
	}
	if len(files) == 0 {
		res.Parts = append(res.Parts, &Part{Name: "history", Pass: false,
			Detail: fmt.Sprintf("no history files matched %q", spec.History)})
		return res, nil
	}
	sort.Strings(files)
	for _, f := range files {
		data, err := sandbox.ReadFile(f)
		if err != nil {
			res.Parts = append(res.Parts, &Part{Name: filepath.Base(f), Pass: false, Detail: err.Error()})
			return res, nil
		}
		h, err := parseHistory(data)
		if err != nil {
			res.Parts = append(res.Parts, &Part{Name: filepath.Base(f), Pass: false, Detail: "corrupt history: " + err.Error()})
			return res, nil
		}
		lin, witness, lerr := linearizable(h.Ops)
		name := filepath.Base(f)
		switch {
		case lerr != nil:
			res.Parts = append(res.Parts, &Part{Name: name, Pass: false, Detail: lerr.Error()})
			return res, nil
		case !lin:
			res.Parts = append(res.Parts, &Part{Name: name, Pass: false, Detail: witness})
			return res, nil
		default:
			res.Parts = append(res.Parts, &Part{Name: name, Pass: true,
				Detail: fmt.Sprintf("linearizable (%d ops)", len(h.Ops))})
		}
	}
	return res, nil
}

// globHistoryFiles resolves a literal file, directory (all *.json inside),
// or glob pattern.
func globHistoryFiles(p string) ([]string, error) {
	if strings.ContainsAny(p, "*?[") {
		return filepath.Glob(p)
	}
	st, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return []string{p}, nil
	}
	ents, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range ents {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			out = append(out, filepath.Join(p, e.Name()))
		}
	}
	return out, nil
}
