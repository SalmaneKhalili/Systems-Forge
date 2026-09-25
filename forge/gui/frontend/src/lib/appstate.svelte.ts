// Shared app state: persisted settings + focus-mode flags.

import { api } from "./api";

export const appState = (() => {
  let settings = $state<Record<string, string>>({});
  let ready = $state(false);
  let focusOpen = $state(false);
  let focusExID = $state("");

  async function load() {
    try {
      settings = (await api.settings()) ?? {};
    } catch {
      /* backend not reachable (e.g. vite dev outside wails) */
    }
    ready = true;
    applyTheme();
  }

  function applyTheme() {
    const t = settings.theme || "business";
    document.documentElement.dataset.theme = t;
  }

  function save(kv: Record<string, string>) {
    settings = { ...settings, ...kv };
    applyTheme();
    api.saveSettings(settings).catch(() => {});
  }

  function number(key: string, fallback: number): number {
    const n = parseInt(settings[key] ?? "", 10);
    return Number.isFinite(n) && n > 0 ? n : fallback;
  }

  return {
    get ready() {
      return ready;
    },
    get settings() {
      return settings;
    },
    get focusOpen() {
      return focusOpen;
    },
    get focusExID() {
      return focusExID;
    },
    setFocusOpen(v: boolean) {
      focusOpen = v;
    },
    setFocusExID(id: string) {
      focusExID = id;
    },
    load,
    save,
    number,
    workMin: () => number("pomodoro_work", 25),
    breakMin: () => number("pomodoro_break", 5),
    reviewLimit: () => number("review_limit", 20),
  };
})();