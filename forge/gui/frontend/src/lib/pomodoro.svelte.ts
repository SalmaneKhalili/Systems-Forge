// Pomodoro engine shared by the Focus view and the fullscreen focus overlay.
// Each work round is a tracked backend session: created on start, closed
// (with minutes) when the round completes or the user stops early.

import { api } from "./api";

export interface PomodoroHandle {
  readonly phase: "work" | "break";
  readonly running: boolean;
  readonly remainingSec: number;
  readonly rounds: number;
  readonly sessionId: number | null;
  start(exID: string): void;
  pause(): void;
  stop(): void;
  reset(): void;
  setLengths(workMin: number, breakMin: number): void;
  readonly workMin: number;
  readonly breakMin: number;
}

export function createPomodoro(): PomodoroHandle {
  let workMin = $state(25);
  let breakMin = $state(5);
  let phase = $state<"work" | "break">("work");
  let running = $state(false);
  let remainingSec = $state(workMin * 60);
  let rounds = $state(0);
  let sessionId = $state<number | null>(null);
  let elapsed = $state(0);
  let timer: ReturnType<typeof setInterval> | null = null;
  let currentEx = $state("");

  function clearTimer() {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  }

  function closeSession(minutes: number) {
    if (sessionId !== null) {
      api.endFocus(sessionId, minutes).catch(() => {});
      sessionId = null;
    }
  }

  function tick() {
    remainingSec -= 1;
    if (phase === "work") elapsed += 1;
    if (remainingSec > 0) return;

    if (phase === "work") {
      closeSession(Math.max(1, Math.round(elapsed / 60)));
      rounds += 1;
      phase = "break";
      remainingSec = breakMin * 60;
    } else {
      phase = "work";
      remainingSec = workMin * 60;
      elapsed = 0;
      api.startFocus("pomodoro", currentEx)
        .then((s) => (sessionId = s.id))
        .catch(() => {});
    }
  }

  const handle: PomodoroHandle = {
    get phase() {
      return phase;
    },
    get running() {
      return running;
    },
    get remainingSec() {
      return remainingSec;
    },
    get rounds() {
      return rounds;
    },
    get sessionId() {
      return sessionId;
    },
    get workMin() {
      return workMin;
    },
    get breakMin() {
      return breakMin;
    },
    start(exID: string) {
      currentEx = exID;
      if (running) return;
      running = true;
      if (!sessionId) {
        if (phase === "work") {
          api.startFocus("pomodoro", exID)
            .then((s) => (sessionId = s.id))
            .catch(() => {});
        }
      }
      clearTimer();
      timer = setInterval(tick, 1000);
    },
    pause() {
      running = false;
      clearTimer();
    },
    stop() {
      running = false;
      clearTimer();
      if (phase === "work" && elapsed > 0) {
        closeSession(Math.max(1, Math.round(elapsed / 60)));
      }
      elapsed = 0;
      phase = "work";
      remainingSec = workMin * 60;
      sessionId = null;
      rounds = 0;
    },
    reset() {
      stop();
    },
    setLengths(w: number, b: number) {
      if (w > 0) workMin = w;
      if (b > 0) breakMin = b;
      if (!running) remainingSec = phase === "work" ? workMin * 60 : breakMin * 60;
    },
  };

  return handle;
}