// Type mirrors of the Go DTOs in forge/gui/types.go (also generated into
// wailsjs/models.ts by Wails; these stay in sync with the Go structs).

export interface Curriculum {
  modules: ModuleInfo[];
  total: number;
}

export interface ModuleInfo {
  id: string;
  title: string;
  skill: string;
  done: number;
  total: number;
  gate?: string;
  exercises: ExerciseInfo[];
}

export type ExStatus = "pass" | "fail" | "untried" | "error";

export interface ExerciseInfo {
  id: string;
  key: string;
  title: string;
  lang: string;
  status: ExStatus;
  attempts: number;
  lastAt?: string;
  isGate: boolean;
}

export interface ExerciseDetail {
  id: string;
  title: string;
  module: string;
  lang: string;
  spec: string;
  readings: Reading[];
  qa: QA[];
  methods: MethodInfo[];
  files: FileInfo[];
  status: ExStatus;
  lastRun?: RunInfo;
  history: RunInfo[];
}

export interface Reading {
  source: string;
  author?: string;
  chapter?: string;
  section?: string;
  url?: string;
}

export interface QA {
  q: string;
  a: string;
}

export interface MethodInfo {
  type: string;
  label: string;
  parts: number;
}

export interface FileInfo {
  path: string;
  name: string;
  size: number;
  dir: boolean;
}

export interface RunInfo {
  status: ExStatus;
  at: string;
  durationMs: number;
  output: string;
}

export interface Progress {
  passed: number;
  attempted: number;
  total: number;
  streak: number;
  todayRuns: number;
  activity: DayActivity[];
  minutes: DayActivity[];
  weekMin: number;
  perModule: ModuleBar[];
}

export interface ModuleBar {
  id: string;
  done: number;
  total: number;
}

export interface DayActivity {
  day: string;
  count: number;
}

export interface CheckResult {
  id: string;
  pass: boolean;
  durationMs: number;
  methods: MethodResult[];
  output: string;
}

export interface MethodResult {
  type: string;
  label: string;
  pass: boolean;
  durMs: number;
  output: string;
  parts: PartResult[];
}

export interface PartResult {
  name: string;
  pass: boolean;
  detail: string;
}

export interface Card {
  key: string;
  module: string;
  exID: string;
  title: string;
  q: string;
  a: string;
  ease: number;
  interval: number;
  due?: string;
  reps: number;
  lapses: number;
  new: boolean;
  reviewable: boolean;
}

export interface StudySession {
  id: number;
  kind: string;
  exercise: string;
  startedAt: string;
  minutes: number;
  open: boolean;
}

export interface Stats {
  totalMinutes: number;
  perExercise: Record<string, number>;
  recent: StudySession[];
}

export interface Quest {
  key: string;
  label: string;
  target: number;
  done: number;
  unit: string;
}

export interface DailyPlan {
  date: string;
  quests: Quest[];
  dueToday: number;
  streak: number;
  doneAll: boolean;
}

export interface HintsResult {
  hints: string[];
  fails: number;
}

export interface UserCardInfo {
  id: number;
  q: string;
  a: string;
  reps: number;
  interval: number;
  due: string;
}

export interface BadgeInfo {
  id: string;
  name: string;
  desc: string;
  icon: string;
  earned: boolean;
  unlockedAt?: string;
}

export interface XPInfo {
  xp: number;
  level: number;
  xpIn: number;
  xpNext: number;
  badges: BadgeInfo[];
  passed: number;
  total: number;
}

export interface WorkDirInfo {
  key: string;
  title: string;
  dir: string;
}

export const STATUS_LABEL: Record<ExStatus, string> = {
  pass: "Passed",
  fail: "In progress",
  untried: "Not started",
  error: "Error",
};