<script lang="ts">
  import { onMount } from "svelte";
  import { api, fmtMins } from "../lib/api";
  import { appState } from "../lib/appstate.svelte";
  import { createPomodoro } from "../lib/pomodoro.svelte";
  import { AMBIENCE_LABELS, setAmbience, type AmbienceKind } from "../lib/ambience";
  import type { Curriculum, Progress, Stats } from "../lib/types";

  const pomo = createPomodoro();

  let cur = $state<Curriculum | null>(null);
  let prog = $state<Progress | null>(null);
  let stats = $state<Stats | null>(null);
  let err = $state("");
  let selectedEx = $state("");
  let workInput = $state(String(appState.workMin()));
  let breakInput = $state(String(appState.breakMin()));
  let freeFocus = $state(false);
  let freeSec = $state(0);
  let freeTimer: ReturnType<typeof setInterval> | null = null;
  let freeSessionId: number | null = null;
  let ambKind = $state<AmbienceKind>("none");
  let ambVol = $state(0.5);

  function fmt(sec: number): string {
    const m = Math.floor(sec / 60);
    const s = sec % 60;
    return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
  }

  onMount(async () => {
    pomo.setLengths(appState.workMin(), appState.breakMin());
    try {
      [cur, prog, stats] = await Promise.all([api.curriculum(), api.progress(), api.stats()]);
      if (cur && cur.modules.length > 0 && cur.modules[0].exercises.length > 0) {
        selectedEx = cur.modules[0].exercises[0].key;
      }
      const s = await api.settings();
      const saved = (s?.ambience ?? "none") as AmbienceKind;
      if (saved === "brown" || saved === "rain" || saved === "wind") {
        ambKind = saved;
        const v = parseFloat(s?.ambienceVol ?? "");
        if (Number.isFinite(v)) ambVol = Math.min(1, Math.max(0, v));
        setAmbience(ambKind, ambVol);
      }
    } catch (e) {
      err = String(e);
    }
  });

  function pickAmb(k: AmbienceKind) {
    ambKind = k;
    setAmbience(k, ambVol);
    api.saveSettings({ ambience: k, ambienceVol: String(ambVol) }).catch(() => {});
  }

  function changeVol() {
    setAmbience(ambKind, ambVol);
    api.saveSettings({ ambience: ambKind, ambienceVol: String(ambVol) }).catch(() => {});
  }

  function applyLengths() {
    const w = parseInt(workInput, 10) || 0;
    const b = parseInt(breakInput, 10) || 0;
    if (w > 0 && b >= 0) {
      pomo.setLengths(w, Math.max(1, b));
      appState.save({ pomodoro_work: String(w), pomodoro_break: String(Math.max(1, b)) });
    }
  }

  function toggleFree() {
    if (freeFocus) {
      freeFocus = false;
      if (freeTimer) clearInterval(freeTimer);
      freeTimer = null;
      if (freeSessionId !== null) {
        api.endFocus(freeSessionId, Math.max(1, Math.round(freeSec / 60))).catch(() => {});
        freeSessionId = null;
      }
      freeSec = 0;
      refreshStats();
      refreshProgress();
    } else {
      freeFocus = true;
      freeSec = 0;
      api.startFocus("focus", selectedEx).then((s) => (freeSessionId = s.id)).catch(() => {});
      freeTimer = setInterval(() => (freeSec += 1), 1000);
    }
  }

  async function refreshStats() {
    try {
      stats = await api.stats();
    } catch {}
  }
  async function refreshProgress() {
    try {
      prog = await api.progress();
    } catch {}
  }

  function deepFocus() {
    appState.setFocusExID(selectedEx);
    appState.setFocusOpen(true);
  }

  function moduleGroups() {
    return cur?.modules ?? [];
  }
</script>

<div class="p-6 lg:p-8 max-w-6xl mx-auto">
  <div class="flex items-center justify-between flex-wrap gap-4">
    <div>
      <h1 class="text-3xl font-bold">Focus</h1>
      <p class="opacity-60 mt-1">Pomodoro + time tracking. Every focus block is logged.</p>
    </div>
    <button class="btn btn-primary" onclick={deepFocus}>Enter deep focus</button>
  </div>

  {#if err}
    <div class="alert alert-error mt-4">{err}</div>
  {/if}

  <div class="grid grid-cols-1 lg:grid-cols-5 gap-6 mt-8">
    <!-- timer card -->
    <div class="lg:col-span-3 card bg-base-200 rounded-2xl p-8 flex flex-col items-center">
      <div
        class="text-[6rem] font-light tabular-nums {pomo.phase === 'break' ? 'text-success' : 'text-base-content'}"
      >
        {fmt(pomo.remainingSec)}
      </div>
      <div class="flex items-center gap-2 mt-2">
        <span class="badge {pomo.phase === 'work' ? 'badge-primary' : 'badge-success'}">
          {pomo.phase === "work" ? `Focus · round ${pomo.rounds + 1}` : "Break"}
        </span>
        {#if pomo.running && pomo.sessionId}
          <span class="text-xs opacity-40">logging…</span>
        {/if}
      </div>

      <div class="flex items-center gap-3 mt-8">
        {#if pomo.running}
          <button class="btn btn-lg btn-circle" onclick={() => pomo.pause()}>
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M6 5h4v14H6zM14 5h4v14h-4z" /></svg>
          </button>
        {:else}
          <button class="btn btn-lg btn-primary btn-circle" onclick={() => pomo.start(selectedEx)}>
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
          </button>
        {/if}
        <button class="btn btn-lg btn-circle btn-ghost" onclick={() => pomo.reset()} title="Reset">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v6h6M20 20v-6h-6M20 9A8 8 0 004.3 6.7L4 6m16 9l-.3 2.3A8 8 0 014 15" />
          </svg>
        </button>
      </div>

      <div class="w-full max-w-sm mt-8">
        <label class="label text-xs opacity-60">Working on</label>
        <select class="select select-sm select-bordered w-full" bind:value={selectedEx}>
          {#each moduleGroups() as m}
            <optgroup label={m.title}>
              {#each m.exercises as ex}
                <option value={ex.key}>{ex.id} — {ex.title}</option>
              {/each}
            </optgroup>
          {/each}
        </select>

        <div class="grid grid-cols-2 gap-3 mt-4">
          <label class="label">
            <span class="text-xs opacity-60">Work min</span>
            <input class="input input-xs input-bordered w-20" type="number" min="1" max="90" bind:value={workInput} />
          </label>
          <label class="label">
            <span class="text-xs opacity-60">Break min</span>
            <input class="input input-xs input-bordered w-20" type="number" min="1" max="30" bind:value={breakInput} />
          </label>
        </div>
        <button class="btn btn-xs btn-ghost w-full" onclick={applyLengths}>Apply lengths</button>

        <div class="divider my-3 text-xs opacity-40">or</div>

        <button class="btn {freeFocus ? 'btn-warning' : 'btn-outline'} btn-sm w-full" onclick={toggleFree}>
          {#if freeFocus}
            Stop free focus · {fmt(freeSec)}
          {:else}
            Start free focus (no timer)
          {/if}
        </button>

        <div class="mt-6">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs opacity-60 font-semibold">Ambience</span>
            <span class="flex gap-1">
              {#each AMBIENCE_LABELS as a}
                <button
                  class="btn btn-xs {ambKind === a.kind ? 'btn-primary' : 'btn-ghost'}"
                  onclick={() => pickAmb(a.kind)}
                >
                  {a.label}
                </button>
              {/each}
            </span>
          </div>
          {#if ambKind !== "none"}
            <input
              class="range range-primary range-xs w-full"
              type="range"
              min="0"
              max="1"
              step="0.05"
              bind:value={ambVol}
              oninput={changeVol}
              title={`Volume ${Math.round(ambVol * 100)}%`}
            />
          {/if}
        </div>
      </div>
    </div>

    <!-- stats card -->
    <div class="lg:col-span-2 space-y-6">
      <div class="card bg-base-200 rounded-2xl p-6">
        <h2 class="font-semibold mb-4">This week</h2>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <div class="text-2xl font-bold text-primary">{fmtMins(prog?.weekMin ?? 0)}</div>
            <div class="text-xs opacity-60">focused</div>
          </div>
          <div>
            <div class="text-2xl font-bold">{stats?.totalMinutes ? fmtMins(stats.totalMinutes) : "0m"}</div>
            <div class="text-xs opacity-60">all time</div>
          </div>
        </div>
        <div class="flex items-end gap-2 h-20 mt-4">
          {#each prog?.minutes ?? [] as d}
            {@const max = Math.max(1, ...(prog!.minutes.map((x) => x.count) || [1]))}
            <div class="flex-1 flex flex-col items-center gap-1">
              <div class="w-full rounded-t-md bg-primary/70" style="height:{Math.max(4, (d.count / max) * 64)}px" title={`${d.count} min`}></div>
              <span class="text-[10px] opacity-50">{d.day.slice(8)}</span>
            </div>
          {/each}
        </div>
      </div>

      <div class="card bg-base-200 rounded-2xl p-6">
        <h2 class="font-semibold mb-3">Recent sessions</h2>
        {#if stats && (stats.recent ?? []).length === 0}
          <p class="text-sm opacity-50">No sessions yet.</p>
        {/if}
        <div class="space-y-1.5">
          {#each stats?.recent.slice(0, 8) ?? [] as s}
            <div class="flex items-center gap-3 text-sm">
              <span class={`w-2 h-2 rounded-full ${s.open ? "bg-warning" : "bg-primary"}`}></span>
              <span class="text-xs opacity-60">{s.kind}</span>
              {#if s.exercise}
                <span class="font-mono text-xs opacity-60 truncate">{s.exercise}</span>
              {:else}
                <span class="opacity-40 text-xs">—</span>
              {/if}
              <span class="ml-auto text-xs opacity-60">{s.minutes}m</span>
              <span class="text-[10px] opacity-40">{new Date(s.startedAt).toLocaleDateString()}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>