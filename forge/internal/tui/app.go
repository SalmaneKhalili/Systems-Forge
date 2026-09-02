// Package tui renders the curriculum in the terminal: module tree, exercise
// specs (markdown), on-demand grading, and live progress.
package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"forge/internal/check"
	"forge/internal/cur"
	"forge/internal/sandbox"
	"forge/internal/store"
)

// status symbols
const (
	markPass    = "✓"
	markFail    = "✗"
	markUntried = "·"
)

var (
	styleHeader = lipgloss.NewStyle().Bold(true)
	styleModule = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	styleActive = lipgloss.NewStyle().Background(lipgloss.Color("237")).Foreground(lipgloss.Color("15"))
	styleDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	stylePass   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	styleFail   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	styleEx     = lipgloss.NewStyle()
	styleFooter = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	styleBanner = lipgloss.NewStyle().Background(lipgloss.Color("248")).Foreground(lipgloss.Color("0"))
)

type entryKind int

const (
	entryModule entryKind = iota
	entryExercise
)

type entry struct {
	kind   entryKind
	modID  string
	exID   string
	title  string
	done   int // completed exercises under a module
	total  int
	status string // store status for an exercise: "", pass, fail
}

type app struct {
	root string
	eng  *check.Engine
	cur  *cur.Curriculum
	rmd  *glamour.TermRenderer

	entries []*entry
	sel     int
	winTop  int
	winH    int

	viewPort  viewport.Model
	mode      string // browse | spec | score
	spinner   spinner.Model
	checking  bool
	lastCheck string

	width  int
	height int
	leftW  int
}

// Run launches the TUI. root is the repository root.
func Run(root string) error {
	m, err := newApp(root)
	if err != nil {
		return err
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func newApp(root string) (*app, error) {
	eng, err := check.New(root)
	if err != nil {
		return nil, err
	}
	curriculum, err := cur.Load(eng.Subjects)
	if err != nil {
		eng.Close()
		return nil, err
	}
	rmd, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(100))
	if err != nil {
		eng.Close()
		return nil, err
	}
	a := &app{
		root:    root,
		eng:     eng,
		cur:     curriculum,
		rmd:     rmd,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		leftW:   30,
		winTop:  0,
		winH:    20,
	}
	rebuild(a)
	a.viewPort = viewport.New(80, 24)
	a.selectFirstExercise()
	return a, nil
}

// selectFirstExercise puts the cursor on the first exercise so the spec pane
// shows subject content immediately on startup.
func (a *app) selectFirstExercise() {
	for i, e := range a.entries {
		if e.kind == entryExercise {
			a.setSel(i)
			return
		}
	}
}

func rebuild(a *app) {
	a.entries = a.entries[:0]
	for _, mod := range a.cur.Modules {
		done := 0
		for _, ex := range mod.Exercises {
			last, _ := a.eng.Store.LastRun(ex.ID)
			if last != nil && last.Status == store.StatusPass {
				done++
			}
		}
		a.entries = append(a.entries, &entry{
			kind: entryModule, modID: mod.ID, title: mod.Title,
			done: done, total: len(mod.Exercises),
		})
		for _, ex := range mod.Exercises {
			last, _ := a.eng.Store.LastRun(ex.ID)
			st := ""
			if last != nil {
				st = store.StatusPass
				if last.Status != store.StatusPass {
					st = store.StatusFail
				}
			}
			a.entries = append(a.entries, &entry{
				kind: entryExercise, exID: ex.ID, title: ex.Title, status: st,
			})
		}
	}
}

func (a *app) selected() *entry {
	if a.sel < 0 || a.sel >= len(a.entries) {
		return nil
	}
	return a.entries[a.sel]
}

func (a *app) findEx(id string) *cur.Exercise {
	for _, mod := range a.cur.Modules {
		if ex := mod.Exercise(id); ex != nil {
			return ex
		}
	}
	return nil
}

func (a *app) initViewport() {
	a.viewPort = viewport.New(a.width-a.leftW-4, a.height-3)
	a.viewPort.SetContent("")
	a.viewPort.Style = lipgloss.NewStyle().Padding(0, 1)
}

func (a *app) loadSpec() {
	e := a.selected()
	if e == nil || e.kind == entryModule {
		a.mode = "browse"
		return
	}
	ex := a.findEx(e.exID)
	if ex == nil {
		return
	}
	data, err := sandbox.ReadFile(ex.Dir + "/subject.md")
	if err != nil {
		a.viewPort.SetContent("no subject.md yet for " + ex.ID)
		return
	}
	rendered, err := a.rmd.Render(string(data))
	if err != nil {
		rendered = string(data)
	}
	a.viewPort.SetContent(rendered)
}

func (a *app) loadScore() {
	a.mode = "score"
	var b strings.Builder
	for _, mod := range a.cur.Modules {
		done := 0
		for _, ex := range mod.Exercises {
			last, _ := a.eng.Store.LastRun(ex.ID)
			if last != nil && last.Status == store.StatusPass {
				done++
			}
		}
		bar := bar(done, len(mod.Exercises))
		fmt.Fprintf(&b, "%s  %d/%d  %s\n", mod.ID, done, len(mod.Exercises), bar)
	}
	a.viewPort.SetContent(b.String())
}

func bar(done, total int) string {
	const width = 20
	if total == 0 {
		return ""
	}
	filled := done * width / total
	return "[" + strings.Repeat("#", filled) + strings.Repeat(".", width-filled) + "]"
}

// msgCheck signals a completed grading run.
type msgCheck struct {
	id     string
	report *check.Report
	err    error
}

func (a *app) Init() tea.Cmd {
	return a.spinner.Tick
}

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = m.Width, m.Height
		a.winH = m.Height - 2
		a.leftW = minInt(34, a.width/3)
		if a.width < 60 {
			a.leftW = 22
		}
		a.initViewport()
		if a.mode == "score" {
			a.loadScore()
		} else {
			a.loadSpec()
		}
		return a, nil

	case tea.KeyMsg:
		return a.handleKey(m)

	case spinner.TickMsg:
		var cmd tea.Cmd
		if a.checking {
			a.spinner, cmd = a.spinner.Update(m)
		}
		return a, cmd

	case msgCheck:
		a.checking = false
		if m.err != nil {
			a.lastCheck = "error: " + m.err.Error()
		} else {
			if m.report.Pass {
				a.lastCheck = "PASS " + m.id
			} else {
				a.lastCheck = "FAIL " + m.id
			}
		}
		rebuild(a)
		return a, nil
	}
	return a, nil
}

func (a *app) handleKey(m tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.String() {
	case "q", "ctrl+c":
		return a, tea.Quit

	case "j", "down":
		if a.mode == "browse" {
			a.move(1)
		}
	case "k", "up":
		if a.mode == "browse" {
			a.move(-1)
		}
	case "g":
		if a.mode == "browse" {
			a.setSel(0)
		} else {
			a.viewPort.GotoTop()
		}
	case "G":
		if a.mode == "browse" {
			a.setSel(len(a.entries) - 1)
		} else {
			a.viewPort.GotoBottom()
		}
	case "n":
		if a.mode == "browse" {
			for i := a.sel + 1; i < len(a.entries); i++ {
				if a.entries[i].kind == entryExercise {
					a.setSel(i)
					break
				}
			}
		}
	case "enter", "l", "right":
		a.enter()
	case "tab":
		a.enter()
	case "e", "h", "left":
		a.mode = "browse"
	case "x":
		return a, a.checkSel()
	case "s":
		a.loadScore()
	case "b":
		a.mode = "browse"
		a.loadSpec()
	case "?":
		a.viewPort.SetContent(helpText())
		a.mode = "spec"
	case "pgdown", " ":
		if a.mode != "browse" {
			a.viewPort.LineDown(1)
		}
	case "pgup":
		if a.mode != "browse" {
			a.viewPort.LineUp(1)
		}
	}
	return a, nil
}

func (a *app) move(delta int) {
	i := a.sel + delta
	if i < 0 {
		i = 0
	}
	if i >= len(a.entries) {
		i = len(a.entries) - 1
	}
	if i == a.sel {
		return
	}
	a.setSel(i)
}

// setSel moves the cursor to entries[i], keeps it on screen, and refreshes
// the spec pane with the selected exercise's subject markdown.
func (a *app) setSel(i int) {
	if i < 0 || i >= len(a.entries) {
		return
	}
	a.sel = i
	a.scrollIntoView()
	if a.entries[i].kind == entryExercise {
		e := a.findEx(a.entries[i].exID)
		if e != nil {
			if data, err := sandbox.ReadFile(e.Dir + "/subject.md"); err == nil {
				if r, err := a.rmd.Render(string(data)); err == nil {
					a.viewPort.SetContent(r)
				}
			}
		}
	}
}

func (a *app) enter() {
	e := a.selected()
	if e == nil {
		return
	}
	if e.kind == entryModule {
		return
	}
	a.mode = "spec"
}

func (a *app) checkSel() tea.Cmd {
	e := a.selected()
	if e == nil || e.kind != entryExercise {
		return nil
	}
	ex := a.findEx(e.exID)
	if ex == nil {
		return nil
	}
	a.checking = true
	a.lastCheck = ""
	id := ex.ID
	return func() tea.Msg {
		report, err := a.eng.Check(context.Background(), ex)
		return msgCheck{id: id, report: report, err: err}
	}
}

func (a *app) scrollIntoView() {
	if a.sel < a.winTop {
		a.winTop = a.sel
	}
	if a.sel >= a.winTop+a.winH {
		a.winTop = a.sel - a.winH + 1
	}
}

func (a *app) View() string {
	if a.width == 0 {
		return "loading…"
	}
	return a.renderHeader() + "\n" + a.renderBody() + "\n" + a.renderFooter()
}

func (a *app) renderHeader() string {
	var allIDs []string
	for _, mod := range a.cur.Modules {
		for _, ex := range mod.Exercises {
			allIDs = append(allIDs, ex.ID)
		}
	}
	done, _, total, _ := a.eng.Store.OverallProgress(allIDs)
	left := styleHeader.Render(" forge — systems-forge ")
	right := fmt.Sprintf(" %d/%d (%.0f%%) ", done, total, pct(done, total))
	line := styleHeader.Render(left)
	pad := a.width - lipgloss.Width(line) - lipgloss.Width(right)
	if pad < 1 {
		pad = 1
	}
	return line + strings.Repeat(" ", pad) + right
}

func pct(d, t int) float64 {
	if t == 0 {
		return 0
	}
	return 100 * float64(d) / float64(t)
}

func (a *app) renderFooter() string {
	keys := " ↑↓/jk move · enter view · x check · s score · b browse · ? help · q quit "
	if a.checking {
		return styleBanner.Render(a.spinner.View() + " checking " + selectedID(a) + "...")
	}
	f := " " + keys
	if a.lastCheck != "" {
		f += " │ " + a.lastCheck
	}
	pad := a.width - lipgloss.Width(f)
	if pad < 1 {
		pad = 1
	}
	return f + strings.Repeat(" ", pad)
}

func selectedID(a *app) string {
	e := a.selected()
	if e == nil || e.kind != entryExercise {
		return ""
	}
	return e.exID
}

func (a *app) renderBody() string {
	bodyH := a.height - 2 // header + footer occupy their own rows
	if bodyH < 1 {
		bodyH = 1
	}
	leftW := a.leftW
	rightW := a.width - leftW
	if rightW < 20 {
		rightW = 20
		leftW = a.width - rightW - 1
	}
	if a.winH != bodyH {
		a.winH = bodyH
	}
	a.viewPort.Width = rightW
	a.viewPort.Height = bodyH
	return lipgloss.JoinHorizontal(lipgloss.Top, a.renderTree(), a.viewPort.View())
}

func (a *app) renderTree() string {
	n := a.winH
	if n < 1 {
		n = 1
	}
	end := a.winTop + n
	if end > len(a.entries) {
		end = len(a.entries)
	}
	width := a.leftW
	var out []string
	for i := a.winTop; i < end; i++ {
		e := a.entries[i]
		var line string
		if e.kind == entryModule {
			bar := bar(e.done, e.total)
			line = styleModule.Render(fmt.Sprintf(" %s %s %d/%d", e.modID, bar, e.done, e.total))
		} else {
			mark := styleDim.Render(markUntried)
			m := e.status
			if m == store.StatusPass {
				mark = stylePass.Render(markPass)
			} else if m == store.StatusFail {
				mark = styleFail.Render(markFail)
			}
			line = fmt.Sprintf("   %s %s  %s", mark, styleEx.Render(e.exID), e.title)
		}
		if i == a.sel {
			line = styleActive.Render(" " + line + " ")
		}
		out = append(out, truncatePad(line, width))
	}
	for len(out) < n {
		out = append(out, strings.Repeat(" ", width))
	}
	return strings.Join(out, "\n")
}

func truncatePad(s string, width int) string {
	w := lipgloss.Width(s)
	if w > width {
		return s[:maxInt(0, width-1)] + "…"
	}
	return s + strings.Repeat(" ", width-w)
}

func (a *app) Close() {
	if a.eng != nil {
		a.eng.Close()
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func helpText() string {
	return `
# keys

| key | action |
|-----|--------|
| j / k / ↑ / ↓ | move in the tree |
| enter, tab, l | open the exercise spec |
| x | run the grader for the selected exercise |
| s | progress score |
| b / e / h | back to the tree |
| g / G | top / bottom |
| n | next exercise |
| ? | this help |
| q / ctrl+c | quit |
`
}
