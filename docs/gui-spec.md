# Systems-Forge Desktop Study App — Spec & Decision Log

> Status: **in development**
> Owner: user + coding agent
> Canonical repo: `/home/salmane/Repositories/Systems-Forge`

## Goal

Replace the "terrible" terminal-only experience with a **polished native desktop
study app** for the curriculum: progress tracking, curriculum browsing, in-app
grading, spaced-repetition flashcards, an in-app code editor with a full
solve loop, a pomodoro/focus mode with time tracking, and a profile/skills view.

Everything remains **local and offline**: one binary, one SQLite DB, no daemon,
no accounts. The existing CLI + TUI stay untouched and usable.

## Stack (decided 2026-09-25)

| Layer     | Choice |
|-----------|--------|
| Shell     | **Wails v2** native desktop window (system WebKitGTK on Linux) |
| Backend   | Go (in-process `check.Engine`, `cur.Load`, `store`) |
| Frontend  | **Svelte 5 + Vite + TypeScript** |
| Styling   | **Tailwind CSS v4 + DaisyUI** (component library) |
| Editor    | **CodeMirror 6** (langs: C, Go, Python, Bash/Shell) |
| Markdown  | markdown-it + highlight.js |
| Routing   | In-Svelte store/hash router (no server routing) |
| Persist   | Same `progress.db` (SQLite); **new tables only** |

Verified on this machine: `webkit2gtk-4.1` ✅, GTK3 ✅, Go 1.27.1 ✅,
node v26.10.0 ✅, yarn ✅, `DISPLAY=:0` ✅.

## Views

1. **Dashboard** — overall progress ring + bar, current streak, "continue where
   you left off", activity heatmap, time spent this week, recent sessions.
2. **Path** — M0→M17 vertical module stepper; exercises with pass/fail/untried
   dots; gate exercises styled as milestones; module cards expand in place.
3. **Exercise detail** — rendered spec (markdown + syntax highlight), readings
   with links, Q&A list, deliverable files, run history, **Run grader** button
   rendering method → part results inline; live progress events while running.
4. **Editor** — CodeMirror 6 editing solution files under `answers/` (path-
   traversal safe); Save → Grade → result; full solve loop in-app.
5. **Review** — flashcards built from every exercise's `answers[]` Q&A metadata;
   SM-2 scoring (again/hard/good/easy); per-module filter; state persisted.
6. **Focus** — pomodoro timer (default 25/5, configurable) + distraction-free
   fullscreen (spec + editor + grader + timer); active sessions auto-logged.
7. **Profile** — skills per module (reuses `tui/moduleSkills` data), proficiency
   bars, run stats, time per module.

## Backend & data

- Wails App bindings: `GetCurriculum`, `GetExercise`, `ListFiles/ReadFile/
  WriteFile`, `RunCheck`, `GetProgress`, `QuizDeck`, `AnswerCard`,
  `StartFocus`, `EndFocus`, `GetSettings`, `SaveSettings`.
- `store` additions (same SQLite file, new tables):
  - `review_state(card_key PK, ease, interval_days, due, last_review, reps, lapses)`
  - `sessions(id PK, kind, exercise, started_at, ended_at, minutes)`
  - `settings(key PK, value)`
- New queries: streak (consecutive active days from `runs`), activity by day
  (exists), time by day/week/module (from `sessions`).

## Repo layout (new pieces)

```
forge/gui/
  main.go            Wails window config + asset embed
  app.go             App struct + bound methods (all backend calls)
  focus.go           session bookkeeping helpers
  frontend/          Svelte+Vite app (built → dist, embedded by Wails)
    src/lib/api.ts   typed wrappers over generated Wails bindings
    src/lib/sm2.ts   (algorithm lives in Go; JS only calls AnswerCard)
    src/views/…      Dashboard, Path, Exercise, Editor, Review, Focus, Profile
Makefile             + `gui` and `gui-dev` targets
bin/forge-gui        built desktop app (gitignored)
```

## Build / run / verify

- Toolchain: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`, `wails doctor`
- Dev: `make gui-dev` (hot reload); Ship: `make gui` → `bin/forge-gui`
- Verify: `go build ./...`, `go vet ./...`, `wails build` clean; launch on
  `DISPLAY=:0` and smoke-test every view; `progress.db` schema migrates cleanly.

## Order of work

1. Toolchain proof: install wails, scaffold window, render anything from Go.
2. Store tables + backend bindings.
3. Frontend views (dashboard → path → exercise → editor → review → focus → profile).
4. Polish, docs, full smoke test.

## Fallback

If the Wails/webview path hits a blocker, the same frontend + binding layer
ports to a localhost web app (`forge web`) with the identical look and data
model.