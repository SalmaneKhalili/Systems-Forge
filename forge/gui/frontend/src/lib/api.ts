// Thin typed wrappers over the Wails-generated bindings (window.go.main.App).
// Calling the runtime object directly (rather than importing the generated
// wailsjs module) keeps the app working in vite dev even before `wails build`
// regenerates the bindings.

import type {
  Card,
  CheckResult,
  Curriculum,
  DailyPlan,
  ExerciseDetail,
  FileInfo,
  HintsResult,
  Progress,
  Stats,
  StudySession,
  UserCardInfo,
  WorkDirInfo,
  XPInfo,
} from "./types";

export class BackendError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "BackendError";
  }
}

function call<R>(name: string, ...args: unknown[]): Promise<R> {
  const fn = (window as any).go?.main?.App?.[name];
  if (typeof fn !== "function") {
    const err = new BackendError(
      `Backend binding "${name}" is not available. ` +
        `This window may have been opened outside the Wails runtime.`,
    );
    logError(name, err);
    return Promise.reject(err);
  }
  return Promise.resolve(fn(...args)).catch((err: unknown) => {
    const e = new BackendError(err instanceof Error ? err.message : String(err));
    logError(name, e);
    throw e;
  });
}

// logError forwards a frontend error to the backend stderr (via the Log
// binding) so runtime issues show up in the app's output as "[frontend] …".
export function logError(tag: string, e: unknown): void {
  const msg = e instanceof Error ? e.message : String(e);
  console.error(`[frontend] ${tag}: ${msg}`);
  try {
    (window as any)?.go?.main?.App?.Log?.(`${tag}: ${msg}`);
  } catch {
    /* the runtime may be torn down; ignore */
  }
}

export const api = {
  curriculum: () => call<Curriculum>("Curriculum"),
  exercise: (id: string) => call<ExerciseDetail>("GetExercise", id),
  files: (id: string) => call<FileInfo[]>("ListFiles", id),
  readFile: (id: string, rel: string) => call<string>("ReadFile", id, rel),
  writeFile: (id: string, rel: string, content: string) =>
    call<void>("WriteFile", id, rel, content),
  runCheck: (id: string) => call<CheckResult>("RunCheck", id),
  runCheckMethod: (id: string, idx: number) =>
    call<CheckResult>("RunCheckMethod", id, idx),
  hints: (id: string) => call<HintsResult>("Hints", id),
  readNotes: (id: string) => call<string>("ReadNotes", id),
  writeNotes: (id: string, content: string) =>
    call<void>("WriteNotes", id, content),
  progress: () => call<Progress>("Progress"),
  quizDeck: (filter: string) => call<Card[]>("QuizDeck", filter),
  answerCard: (key: string, rating: number) =>
    call<Card>("AnswerCard", key, rating),
  addCard: (q: string, a: string) => call<Card>("AddCard", q, a),
  listUserCards: () => call<UserCardInfo[]>("ListUserCards"),
  updateCard: (id: number, q: string, a: string) =>
    call<void>("UpdateCard", id, q, a),
  deleteCard: (id: number) => call<void>("DeleteCard", id),
  xp: () => call<XPInfo>("XP"),
  startFocus: (kind: string, exID: string) =>
    call<StudySession>("StartFocus", kind, exID),
  endFocus: (id: number, minutes: number) =>
    call<StudySession>("EndFocus", id, minutes),
  stats: () => call<Stats>("Stats"),
  settings: () => call<Record<string, string>>("GetSettings"),
  saveSettings: (kv: Record<string, string>) => call<void>("SaveSettings", kv),
  dailyPlan: () => call<DailyPlan>("DailyPlan"),
  workDirs: () => call<WorkDirInfo[]>("WorkDirs"),
  termStart: (id: string, dir: string, command: string[]) =>
    call<void>("TermStart", id, dir, command),
  termInput: (id: string, data: string) => call<void>("TermInput", id, data),
  termResize: (id: string, cols: number, rows: number) =>
    call<void>("TermResize", id, cols, rows),
  termStop: (id: string) => call<void>("TermStop", id),
};

export function fmtDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`;
  return `${(ms / 1000).toFixed(1)}s`;
}

export function fmtDate(iso?: string): string {
  if (!iso) return "—";
  return new Date(iso).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

export function fmtMins(mins: number): string {
  const m = Math.round(mins);
  if (m < 60) return `${m}m`;
  const h = Math.floor(m / 60);
  const r = m % 60;
  return r ? `${h}h ${r}m` : `${h}h`;
}