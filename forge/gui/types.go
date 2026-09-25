package main

// DTOs shared between the Go backend and the Svelte frontend. These are
// mirrored into frontend models by Wails binding generation.

// Curriculum is the full module tree plus per-exercise status.
type Curriculum struct {
	Modules []*ModuleInfo `json:"modules"`
	Total   int           `json:"total"`
}

type ModuleInfo struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Skill     string          `json:"skill"`
	Done      int             `json:"done"`
	Total     int             `json:"total"`
	Gate      string          `json:"gate,omitempty"`
	Exercises []*ExerciseInfo `json:"exercises"`
}

type ExerciseInfo struct {
	ID       string `json:"id"`
	Key      string `json:"key"` // globally-unique "module/exercise" path
	Title    string `json:"title"`
	Lang     string `json:"lang"`
	Status   string `json:"status"` // pass | fail | untried
	Attempts int    `json:"attempts"`
	LastAt   string `json:"lastAt,omitempty"`
	IsGate   bool   `json:"isGate"`
}

// ExerciseDetail is everything the exercise view needs.
type ExerciseDetail struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Module   string        `json:"module"`
	Lang     string        `json:"lang"`
	Spec     string        `json:"spec"` // raw subject.md markdown
	Readings []*Reading    `json:"readings"`
	QA       []*QA         `json:"qa"`
	Methods  []*MethodInfo `json:"methods"`
	Files    []*FileInfo   `json:"files"`
	Status   string        `json:"status"`
	LastRun  *RunInfo      `json:"lastRun,omitempty"`
	History  []*RunInfo    `json:"history"`
}

type Reading struct {
	Source  string `json:"source"`
	Author  string `json:"author,omitempty"`
	Chapter string `json:"chapter,omitempty"`
	Section string `json:"section,omitempty"`
	URL     string `json:"url,omitempty"`
}

type QA struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type MethodInfo struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Parts int    `json:"parts"`
}

type FileInfo struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Dir  bool   `json:"dir"`
}

type RunInfo struct {
	Status     string `json:"status"`
	At         string `json:"at"`
	DurationMs int64  `json:"durationMs"`
	Output     string `json:"output"`
}

// Progress aggregates overall stats for the dashboard.
type Progress struct {
	Passed    int              `json:"passed"`
	Attempted int              `json:"attempted"`
	Total     int              `json:"total"`
	Streak    int              `json:"streak"`
	TodayRuns int              `json:"todayRuns"`
	Activity  []*DayActivity   `json:"activity"`   // last 90 days
	Minutes   []*DayActivity   `json:"minutes"`    // focused minutes per day, last 7 days
	WeekMin   int              `json:"weekMin"`
	PerModule []*ModuleBar     `json:"perModule"`
}

type ModuleBar struct {
	ID   string `json:"id"`
	Done int    `json:"done"`
	Total int   `json:"total"`
}

type DayActivity struct {
	Day   string `json:"day"`
	Count int    `json:"count"`
}

// CheckResult is one grading run surfaced to the UI.
type CheckResult struct {
	ID         string          `json:"id"`
	Pass       bool            `json:"pass"`
	DurationMs int64           `json:"durationMs"`
	Methods    []*MethodResult `json:"methods"`
	Output     string          `json:"output"`
}

type MethodResult struct {
	Type   string        `json:"type"`
	Label  string        `json:"label"`
	Pass   bool          `json:"pass"`
	DurMs  int64         `json:"durMs"`
	Output string        `json:"output"`
	Parts  []*PartResult `json:"parts"`
}

type PartResult struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Detail string `json:"detail"`
}

// Card is one spaced-repetition review card.
type Card struct {
	Key        string  `json:"key"` // "<exerciseID>:<index>"
	Module     string  `json:"module"`
	ExID       string  `json:"exID"`
	Title      string  `json:"title"`
	Q          string  `json:"q"`
	A          string  `json:"a"`
	Ease       float64 `json:"ease"`
	Interval   float64 `json:"interval"`
	Due        string  `json:"due,omitempty"`
	Reps       int     `json:"reps"`
	Lapses     int     `json:"lapses"`
	New        bool    `json:"new"`
	Reviewable bool    `json:"reviewable"`
}

// StudySession is one recorded focus/pomodoro block.
type StudySession struct {
	ID        int64  `json:"id"`
	Kind      string `json:"kind"`
	Exercise  string `json:"exercise"`
	StartedAt string `json:"startedAt"`
	Minutes   int    `json:"minutes"`
	Open      bool   `json:"open"`
}

type Stats struct {
	TotalMinutes int                 `json:"totalMinutes"`
	PerExercise  map[string]int      `json:"perExercise"`
	Recent       []*StudySession     `json:"recent"`
}