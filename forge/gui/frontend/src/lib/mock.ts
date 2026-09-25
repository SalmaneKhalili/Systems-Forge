/// <reference types="vite/client" />
// Dev-only backend mock. main.ts imports this module dynamically under
// `import.meta.env.DEV`, so production builds never include it (the import
// statement itself is eliminated). It lets the whole app run in a plain
// browser tab without the Wails webview, mirroring the Go DTOs exactly, so
// view rendering can be smoke-tested and screenshotted during development.
import fixture from "./mock-data.json";

type Fn = (...args: unknown[]) => unknown;

const modules = fixture.modules;
const total = modules.reduce((n, m) => n + m.exercises.length, 0);
const detail = (fixture.detail ?? {}) as Record<string, Record<string, unknown>>;
const fileContents = (fixture.fileContents ?? {}) as Record<string, string>;

let sessionSeq = 1;
let cardSeq = 1;

// Dev-only terminal emulation: sessions echo typed input so the full
// EventsOn -> term:out -> xterm pipeline (and the round-trip probe) can be
// exercised without the Go pty side.
const termSessions: Record<string, { buf: string }> = {};
const listeners: Record<string, Array<(d: unknown) => void>> = {};

function b64(s: string): string {
  const bytes = new TextEncoder().encode(s);
  let bin = "";
  for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
  return btoa(bin);
}

function emit(name: string, data: unknown) {
  (listeners[name] ?? []).forEach((cb) => cb(data));
}

function checkResult(id: string) {
  const d = detail[id] ?? {};
  const methods = (d.methods as Array<{ type: string; label: string; parts: number }>) ?? [];
  return {
    id,
    pass: true,
    durationMs: 12,
    output: "[mock] all checks passed",
    methods: methods.map((m) => ({
      type: m.type,
      label: m.label,
      pass: true,
      durMs: 6,
      output: "[mock] ok",
      parts: Array.from({ length: m.parts }, (_, i) => ({
        name: `run ${i + 1}`,
        pass: true,
        detail: "",
      })),
    })),
  };
}

const handlers: Record<string, Fn> = {
  Curriculum: () => ({ modules, total }),
  GetExercise: (id: unknown) => {
    const d = detail[String(id)];
    if (!d) return Promise.reject(new Error(`exercise not found: ${id}`));
    return { ...d, lastRun: null } as Record<string, unknown>;
  },
  ListFiles: (id: unknown) => (detail[String(id)]?.files as unknown[]) ?? [],
  ReadFile: (id: unknown, rel: unknown) => fileContents[`${id}\u0000${rel}`] ?? "",
  WriteFile: () => undefined,
  RunCheck: (id: unknown) => checkResult(String(id)),
  RunCheckMethod: (id: unknown) => checkResult(String(id)),
  Hints: (id: unknown) => ({
    hints: [
      "Re-check the run command expected by the grader (see the spec).",
      "Compare a passing example against your current build output.",
      "Isolate the failing step with a small scratch run before editing sources.",
    ],
    fails: 0,
  }),
  ReadNotes: () => "",
  WriteNotes: () => undefined,
  Progress: () => ({
    passed: modules.filter((m) => m.exercises.some(() => false)).length,
    attempted: 0,
    total,
    streak: 0,
    todayRuns: 0,
    activity: [],
    minutes: [],
    weekMin: 0,
    perModule: modules.map((m) => ({ id: m.id, done: 0, total: m.exercises.length })),
  }),
  QuizDeck: () => [],
  AnswerCard: () => null,
  AddCard: (q: unknown, a: unknown) => ({
    id: cardSeq++,
    q: String(q),
    a: String(a),
    reps: 0,
    interval: 0,
    due: new Date().toISOString(),
  }),
  ListUserCards: () => [],
  UpdateCard: () => undefined,
  DeleteCard: () => undefined,
  XP: () => ({
    xp: 0,
    level: 1,
    xpIn: 0,
    xpNext: 100,
    badges: [] as unknown[],
    passed: 0,
    total,
  }),
  StartFocus: (kind: unknown, exID: unknown) => ({
    id: sessionSeq++,
    kind: String(kind ?? "focus"),
    exercise: String(exID ?? ""),
    startedAt: new Date().toISOString(),
    minutes: 0,
    open: true,
  }),
  EndFocus: () => undefined,
  Stats: () => ({ totalMinutes: 0, perExercise: {}, recent: [] }),
  GetSettings: () => ({}),
  SaveSettings: () => undefined,
  DailyPlan: () => ({
    date: new Date().toISOString().slice(0, 10),
    quests: [] as unknown[],
    dueToday: 0,
    streak: 0,
    doneAll: false,
  }),
  WorkDirs: () => [{ key: "", title: "Repository root", dir: "/tmp/opencode/mock-repo" }],
  TermStart: (id: unknown) => {
    termSessions[String(id)] = { buf: "" };
    emit(
      `term:out:${String(id)}`,
      b64("\r\n[mock shell]  type here — echo __sf_term_ok is answered\r\n$ "),
    );
    return undefined;
  },
  TermInput: (id: unknown, data: unknown) => {
    const s = termSessions[String(id)];
    if (!s) return;
    const chunk = String(data ?? "");
    s.buf += chunk;
    emit(`term:out:${String(id)}`, b64(chunk)); // echo typing back
    if (chunk.includes("\r") || chunk.includes("\n")) {
      const line = s.buf.split(/\r?\n/).pop()?.replace(/\r/g, "") ?? "";
      if (line.trim() === "echo __sf_term_ok") {
        emit(`term:out:${String(id)}`, b64("__sf_term_ok\r\n$ "));
      } else if (line.trim() === "exit") {
        emit(`term:exit:${String(id)}`, true);
      } else {
        emit(`term:out:${String(id)}`, b64("$ "));
      }
    }
  },
  TermResize: () => undefined,
  TermStop: (id: unknown) => {
    delete termSessions[String(id)];
    emit(`term:exit:${String(id)}`, true);
  },
  Log: (msg: unknown) => {
    console.log("[mock-log]", String(msg));
    return undefined;
  },
  SaveProgress: () => undefined,
};

export function installMockBackend(): void {
  if ((window as any).go?.main?.App) return; // real Wails runtime already present
  const App = new Proxy(
    {},
    {
      get: (_t, name: string) => {
        const fn = handlers[name];
        return fn ? (...args: unknown[]) => Promise.resolve(fn(...args)) : () => Promise.resolve(null);
      },
    },
  );
  (window as any).go = { main: { App } };
  (window as any).runtime = {
    EventsOn: (name: string, cb: (d: unknown) => void) => {
      listeners[name] = [...(listeners[name] ?? []), cb];
      return undefined;
    },
    EventsOff: (name: string) => {
      delete listeners[name];
      return undefined;
    },
    EventsEmit: (name: string, data: unknown) => {
      emit(name, data);
      return undefined;
    },
    Log: (msg: unknown) => Promise.resolve(console.log("[mock-log]", String(msg))),
  };
}