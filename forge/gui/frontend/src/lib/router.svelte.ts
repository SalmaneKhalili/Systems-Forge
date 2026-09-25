// Minimal hash router for a single-window desktop app.
// Routes:
//   #/              dashboard
//   #/path          curriculum path
//   #/ex/<id>       exercise (spec + editor + grader)
//   #/review        flashcards
//   #/focus         pomodoro / focus mode
//   #/terminal      terminal (shell / nvim)
//   #/profile       profile & skills

import { get } from "svelte/store";
import { readable } from "svelte/store";

export interface Route {
  name: string;
  param?: string;
}

function parse(hash: string): Route {
  const h = hash.replace(/^#\/?/, "");
  const parts = h.split("/").filter(Boolean);
  if (parts.length === 0) return { name: "dashboard" };
  if (parts[0] === "ex" && parts[1]) {
    return { name: "exercise", param: decodeParam(parts[1]) };
  }
  return { name: parts[0] };
}

// Decode the percent-encoded route param (keys contain "/" which we encode
// as %2F in goto()).
function decodeParam(s: string): string {
  try {
    return decodeURIComponent(s);
  } catch {
    return s;
  }
}

function currentRoute(): Route {
  return parse(window.location.hash);
}

export const route = readable<Route>(currentRoute(), (set) => {
  const onChange = () => set(currentRoute());
  window.addEventListener("hashchange", onChange);
  return () => window.removeEventListener("hashchange", onChange);
});

export function goto(name: string, param?: string) {
  const hash = param ? `#/${name}/${encodeURIComponent(param)}` : `#/${name}`;
  if (window.location.hash === hash) {
    // Re-emit so stateful views can react to re-navigation.
    window.dispatchEvent(new HashChangeEvent("hashchange"));
    return;
  }
  window.location.hash = hash;
}

export function current(): Route {
  return get(route);
}