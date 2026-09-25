// Session-direction helpers: "Surprise me" and the next-up continue chain.
import { goto } from "./router.svelte";
import type { Curriculum, ExerciseInfo } from "./types";

// pendingExercises returns every exercise whose latest status isn't "pass",
// in curriculum order (M0 → M17).
export function pendingExercises(cur: Curriculum | null): ExerciseInfo[] {
  if (!cur) return [];
  const out: ExerciseInfo[] = [];
  for (const m of cur.modules) {
    for (const ex of m.exercises) {
      if (ex.status !== "pass") out.push(ex);
    }
  }
  return out;
}

// surprise jumps to a random pending exercise (falls back to the Path view
// when the curriculum is complete).
export function surprise(cur: Curriculum | null): void {
  const pending = pendingExercises(cur);
  if (pending.length === 0) {
    goto("path");
    return;
  }
  const ex = pending[Math.floor(Math.random() * pending.length)];
  goto("exercise", ex.key);
}

// Continue: due review cards first, otherwise the first pending exercise.
export function nextUp(
  cur: Curriculum | null,
  dueToday: number,
): { name: string; param?: string; label: string } {
  if (cur === null) return { name: "path", label: "Path" };
  if (dueToday > 0) return { name: "review", label: `Review ${dueToday} due` };
  const pending = pendingExercises(cur);
  if (pending.length > 0) {
    return { name: "exercise", param: pending[0].key, label: `Continue: ${pending[0].title}` };
  }
  return { name: "dashboard", label: "All caught up 🎉" };
}