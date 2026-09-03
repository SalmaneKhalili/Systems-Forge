package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"forge/internal/cur"
)

// --- styles for the polished views ---------------------------------------

var (
	styleTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	styleSection = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45"))
	styleStat    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
	styleMuted   = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
)

// heat colors: GitHub-style contribution ramp (index 0 = no activity).
var heatColors = []string{"240", "22", "28", "34", "82", "46"}

// heatFgColors: readable foreground ramp for the text-glyph activity strip.
var heatFgColors = []string{"243", "244", "28", "34", "40", "46"}

func heatColor(level int) string {
	if level < 0 {
		level = 0
	}
	if level >= len(heatColors) {
		level = len(heatColors) - 1
	}
	return heatColors[level]
}

// heatLevel maps an activity count to a 0..5 ramp index.
func heatLevel(n int) int {
	switch {
	case n <= 0:
		return 0
	case n <= 2:
		return 1
	case n <= 4:
		return 2
	case n <= 6:
		return 3
	case n <= 9:
		return 4
	default:
		return 5
	}
}

// renderHeatHome renders a compact single-line activity strip (one block per
// day for the last 90 days) for the dashboard.
func (a *app) renderHeatHome() string {
	activity := a.loadActivity()
	today := time.Now().UTC()
	start := today.AddDate(0, 0, -89)
	var b strings.Builder
	for d := start; !d.After(today); d = d.AddDate(0, 0, 1) {
		n := activity[d.UTC().Format("2006-01-02")]
		block := "▂"
		switch heatLevel(n) {
		case 0:
			block = "·"
		case 1:
			block = "▁"
		case 2:
			block = "▂"
		case 3:
			block = "▄"
		case 4:
			block = "▆"
		default:
			block = "█"
		}
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(heatFgColors[heatLevel(n)])).Render(block))
	}
	b.WriteString(styleMuted.Render("  last 90 days"))
	return b.String()
}

// renderHeatmap renders a GitHub-style year grid: columns are weeks, rows are
// weekdays. Each cell is a colored block for one day's activity.
func (a *app) renderHeatmap() string {
	activity := a.loadActivity()
	today := time.Now().UTC()
	weeks := 26 // half-year rolling window; GitHub uses ~52

	// Find the Sunday that starts the oldest week of the window.
	end := today
	start := today.AddDate(0, 0, -(7*weeks - 1))
	for start.Weekday() != time.Sunday {
		start = start.AddDate(0, 0, -1)
	}

	cell := func(d time.Time) string {
		n := activity[d.UTC().Format("2006-01-02")]
		c := heatColor(heatLevel(n))
		return lipgloss.NewStyle().Background(lipgloss.Color(c)).
			Foreground(lipgloss.Color("0")).Render("   ")
	}

	// Build the week columns, tracking each week's Sunday for month labels.
	type wk struct {
		sunday time.Time
		cells  [7]string
	}
	var weeksList []wk
	for sun := start; !sun.After(end); sun = sun.AddDate(0, 0, 7) {
		var w wk
		w.sunday = sun
		for off := 0; off < 7; off++ {
			d := sun.AddDate(0, 0, off)
			w.cells[off] = cell(d)
		}
		weeksList = append(weeksList, w)
	}

	// Month-label row: at the first week containing the 1st of a new month,
	// write the month abbreviation (up to 3 chars per 3-char cell).
	monthRow := make([][]rune, len(weeksList))
	for current := range monthRow {
		monthRow[current] = []rune("   ")
	}
	prevMonth := -1
	for i, w := range weeksList {
		m := int(w.sunday.AddDate(0, 0, 6).Month()) // month of the week's last day
		if m != prevMonth {
			label := []rune(time.Month(m).String()[:3])
			copy(monthRow[i], label)
			prevMonth = m
		}
	}
	monthLine := ""
	for i := range monthRow {
		monthLine += string(monthRow[i])
	}

	weekdayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	var b strings.Builder
	b.WriteString(styleSection.Render(" Activity") + "\n\n")
	b.WriteString("     " + monthLine + "\n")
	for day := 0; day < 7; day++ {
		row := ""
		for _, w := range weeksList {
			row += w.cells[day]
		}
		b.WriteString(fmt.Sprintf("%-4s %s\n", weekdayNames[day], row))
	}
	b.WriteString("\n" + heatLegend() + "\n")
	return b.String()
}

func heatLegend() string {
	levels := []int{0, 1, 3, 5, 8, 12}
	var b strings.Builder
	b.WriteString(styleMuted.Render("Less "))
	for _, lv := range levels {
		c := heatColor(heatLevel(lv))
		b.WriteString(lipgloss.NewStyle().Background(lipgloss.Color(c)).Render("  "))
	}
	b.WriteString(styleMuted.Render(" More"))
	return b.String()
}

// loadActivity loads the per-day run counts for the whole window (returns a
// map keyed by "2006-01-02"; fills every day in the past 6 months with 0).
func (a *app) loadActivity() map[string]int {
	today := time.Now().UTC()
	start := today.AddDate(0, 0, -183)
	got, err := a.eng.Store.ActivityByDay(start, today)
	if err != nil {
		got = map[string]int{}
	}
	// fill any missing days with 0 so colors render consistently
	out := map[string]int{}
	for d := start; !d.After(today); d = d.AddDate(0, 0, 1) {
		key := d.UTC().Format("2006-01-02")
		out[key] = got[key]
	}
	return out
}

// bigBar renders a wide ASCII/SVG-style progress bar with a percent label.
func bigBar(done, total int, width int) string {
	if total <= 0 {
		return ""
	}
	filled := done * width / total
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar + fmt.Sprintf("  %d/%d (%.0f%%)", done, total, pct(done, total))
}

// renderHome renders the dashboard: banner, overall progress, activity strip,
// module completion banners, and a key. This is the landing view.
func (a *app) renderHome() string {
	var allIDs []string
	for _, mod := range a.cur.Modules {
		for _, ex := range mod.Exercises {
			allIDs = append(allIDs, ex.ID)
		}
	}
	passed, attempted, total, _ := a.eng.Store.OverallProgress(allIDs)

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(styleTitle.Render(" systems-forge "))
	b.WriteString(styleMuted.Render("    bootcamp-style distributed systems") + "\n")
	b.WriteString(styleMuted.Render(" 92 exercises · 18 modules · every skill re-deployed") + "\n\n")

	// Overall progress bar
	b.WriteString(styleSection.Render(" Progress") + "\n")
	b.WriteString(" " + bigBar(passed, total, 60) + "\n")
	if attempted < total {
		b.WriteString(styleMuted.Render(fmt.Sprintf(" %d/%d attempted · %d not started\n", attempted, total, total-attempted)))
	}

	// Activity strip
	b.WriteString("\n" + styleSection.Render(" Activity") + "\n")
	b.WriteString(" " + a.renderHeatHome() + "\n")

	// Module completion banners (95-style "Mx complete")
	var completed []string
	for _, mod := range a.cur.Modules {
		done, cnt, _ := a.eng.Store.ModuleProgress(a.exerciseIDs(mod))
		if cnt > 0 && done == cnt {
			completed = append(completed, stylePass.Render("  ✓ "+mod.ID+" complete"))
		} else if cnt > 0 && done > 0 {
			completed = append(completed, styleMuted.Render(fmt.Sprintf("  · %s %d/%d", mod.ID, done, cnt)))
		}
	}
	if len(completed) > 0 {
		b.WriteString("\n" + styleSection.Render(" Modules") + "\n")
		b.WriteString(strings.Join(completed, "\n") + "\n")
	}

	b.WriteString("\n" + styleMuted.Render(" tab: browse exercises · 2: activity · 3: profile · ?: help") + "\n")
	return b.String()
}

// renderProfile renders the skills/profile view: per-module skill, proficiency
// bar, and earned badge when the module is complete.
func (a *app) renderProfile() string {
	var b strings.Builder
	b.WriteString("\n" + styleTitle.Render(" Profile ") + "\n")
	b.WriteString(styleMuted.Render(" skills earned by completing each module's gate") + "\n\n")

	earned := 0
	totalSkills := 0
	for _, mod := range a.cur.Modules {
		done, cnt, _ := a.eng.Store.ModuleProgress(a.exerciseIDs(mod))
		if cnt == 0 {
			continue
		}
		totalSkills++
		complete := done == cnt
		if complete {
			earned++
		}
		badge := styleDim.Render(" ○ ")
		if complete {
			badge = stylePass.Render(" ✓ ")
		} else if done > 0 {
			badge = styleDim.Render(" ◔ ")
		}
		name := skillFor(mod.ID)
		barW := 24
		fill := done * barW / cnt
		bar := strings.Repeat("█", fill) + strings.Repeat("░", barW-fill)
		pctv := pct(done, cnt)
		line := badge + " " + name
		b.WriteString(line + "\n")
		b.WriteString("     " + styleMuted.Render(bar+" "+fmt.Sprintf("%.0f%%", pctv)) + "\n")
	}
	if totalSkills > 0 {
		b.WriteString("\n" + bigBar(earned, totalSkills, 60) + styleMuted.Render(" skills earned") + "\n")
	}
	return b.String()
}

// exerciseIDs returns the exercise IDs for a module.
func (a *app) exerciseIDs(mod *cur.Module) []string {
	ids := make([]string, len(mod.Exercises))
	for i, ex := range mod.Exercises {
		ids[i] = ex.ID
	}
	return ids
}
