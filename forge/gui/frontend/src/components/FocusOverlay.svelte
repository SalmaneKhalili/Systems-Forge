<script lang="ts">
  import { onMount } from "svelte";
  import { appState } from "../lib/appstate.svelte";
  import { createPomodoro } from "../lib/pomodoro.svelte";
  import { goto } from "../lib/router.svelte";
  import { api } from "../lib/api";
  import type { ExerciseInfo } from "../lib/types";

  let {
    open = false,
    onClose,
  }: { open?: boolean; onClose: () => void } = $props();

  const pomo = createPomodoro();
  let ex: ExerciseInfo | null = $state(null);

  function fmt(sec: number): string {
    const m = Math.floor(sec / 60);
    const s = sec % 60;
    return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
  }

  onMount(() => {
    pomo.setLengths(appState.workMin(), appState.breakMin());
    loadEx();
  });

  async function loadEx() {
    const exID = appState.focusExID;
    if (!exID) return;
    try {
      const d = await api.exercise(exID);
      ex = { id: exID, title: d.title, lang: d.lang, status: d.status, attempts: 0, isGate: false };
    } catch {
      /* fine */
    }
  }
</script>

{#if open}
  <div class="fixed inset-0 z-50 bg-base-100 flex flex-col">
    <!-- top bar -->
    <div class="flex items-center justify-between px-6 py-3 border-b border-base-300">
      <div class="text-sm opacity-60">
        Focus mode
        {#if ex}
          — <button class="link link-primary" onclick={() => goto("exercise", ex!.id)}>{ex.title}</button>
        {/if}
      </div>
      <button class="btn btn-sm btn-ghost" onclick={onClose}>
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
        Exit focus
      </button>
    </div>

    <!-- centered timer -->
    <div class="flex-1 flex flex-col items-center justify-center gap-10">
      <div
        class="text-[9rem] font-light tracking-tight tabular-nums {pomo.phase === 'break' ? 'text-success' : 'text-base-content'}"
      >
        {fmt(pomo.remainingSec)}
      </div>

      <div class="flex items-center gap-3 text-sm opacity-70">
        <span class="badge {pomo.phase === 'work' ? 'badge-primary' : 'badge-success'} badge-lg">
          {pomo.phase === "work" ? "Focus" : "Break"}
        </span>
        <span>Round {pomo.rounds + 1}</span>
      </div>

      <div class="flex items-center gap-3">
        {#if pomo.running}
          <button class="btn btn-lg btn-circle" onclick={() => pomo.pause()} title="Pause">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
              <path d="M6 5h4v14H6zM14 5h4v14h-4z" />
            </svg>
          </button>
        {:else}
          <button class="btn btn-lg btn-primary btn-circle" onclick={() => pomo.start(appState.focusExID)} title="Start">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z" />
            </svg>
          </button>
        {/if}
        <button class="btn btn-lg btn-circle btn-ghost" onclick={() => pomo.reset()} title="Reset">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v6h6M20 20v-6h-6M20 9A8 8 0 004.3 6.7L4 6m16 9l-.3 2.3A8 8 0 014 15" />
          </svg>
        </button>
      </div>

      <p class="text-xs opacity-40 max-w-sm text-center">
        Time in focus is logged automatically. Breathe, commit, <span class="opacity-70">build something.</span>
      </p>
    </div>
  </div>
{/if}