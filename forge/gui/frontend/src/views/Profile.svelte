<script lang="ts">
  import { onMount } from "svelte";
  import { api, fmtMins } from "../lib/api";
  import { appState } from "../lib/appstate.svelte";
  import { confetti } from "../lib/celebrate";
  import type { Curriculum, Progress, Stats, XPInfo } from "../lib/types";

  let cur = $state<Curriculum | null>(null);
  let prog = $state<Progress | null>(null);
  let stats = $state<Stats | null>(null);
  let xp = $state<XPInfo | null>(null);
  let err = $state("");
  let savedMsg = $state("");

  let theme = $state("business");
  let workMin = $state("25");
  let breakMin = $state("5");
  let reviewLimit = $state("20");

  const themes = ["business", "dark", "dracula", "night", "synthwave", "winter", "corporate", "emerald"];

  onMount(async () => {
    try {
      [cur, prog, stats, xp] = await Promise.all([
        api.curriculum(),
        api.progress(),
        api.stats(),
        api.xp(),
      ]);
      celebrateNewBadges();
    } catch (e) {
      err = String(e);
    }
    theme = appState.settings.theme ?? "business";
    workMin = String(appState.workMin());
    breakMin = String(appState.breakMin());
    reviewLimit = String(appState.reviewLimit());
  });

  // One-shot confetti when badges we've never seen before are earned.
  function celebrateNewBadges() {
    if (!xp) return;
    const key = "sf_badges_seen";
    let seen: string[] = [];
    try {
      seen = JSON.parse(localStorage.getItem(key) ?? "[]");
    } catch {
      seen = [];
    }
    const earned = xp.badges.filter((b) => b.earned).map((b) => b.id);
    const fresh = earned.filter((id) => !seen.includes(id));
    if (fresh.length > 0) {
      confetti(140);
      localStorage.setItem(key, JSON.stringify(earned));
    }
  }

  function titleFor(id: string): string {
    const m = cur?.modules.find((x) => x.id === id);
    return m?.title ?? id;
  }

  function saveSettings() {
    appState.save({
      theme,
      pomodoro_work: workMin,
      pomodoro_break: breakMin,
      review_limit: reviewLimit,
    });
    savedMsg = "Settings saved";
    setTimeout(() => (savedMsg = ""), 1500);
  }

  function topExercises(): { id: string; minutes: number; title: string }[] {
    if (!cur || !stats) return [];
    const exTitle = new Map<string, string>();
    for (const m of cur.modules) for (const e of m.exercises) exTitle.set(e.id, e.title);
    return Object.entries(stats.perExercise)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 10)
      .map(([id, minutes]) => ({ id, minutes, title: exTitle.get(id) ?? id }));
  }
</script>

<div class="p-6 lg:p-8 max-w-5xl mx-auto">
  <h1 class="text-3xl font-bold">Profile</h1>
  <p class="opacity-60 mt-1">Skills, stats and settings.</p>

  {#if err}
    <div class="alert alert-error mt-4">{err}</div>
  {/if}

  {#if prog && cur}
    <!-- stats row -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-8">
      <div class="card bg-base-200 rounded-2xl p-5">
        <div class="text-2xl font-bold">{prog.passed}/{prog.total}</div>
        <div class="text-xs opacity-60 mt-1">exercises passed</div>
      </div>
      <div class="card bg-base-200 rounded-2xl p-5">
        <div class="text-2xl font-bold">{prog.total - prog.passed}</div>
        <div class="text-xs opacity-60 mt-1">remaining</div>
      </div>
      <div class="card bg-base-200 rounded-2xl p-5">
        <div class="text-2xl font-bold">🔥 {prog.streak}</div>
        <div class="text-xs opacity-60 mt-1">day streak</div>
      </div>
      <div class="card bg-base-200 rounded-2xl p-5">
        <div class="text-2xl font-bold">{fmtMins(stats?.totalMinutes ?? 0)}</div>
        <div class="text-xs opacity-60 mt-1">focused total</div>
      </div>
    </div>

    <!-- XP + badges -->
    {#if xp}
      <div class="card bg-gradient-to-br from-primary/10 to-secondary/10 border border-primary/20 rounded-2xl p-6 mt-6">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-3">
              <div class="badge badge-primary badge-lg">Level {xp.level}</div>
              <div class="text-sm opacity-70">
                {xp.xpIn}/{xp.xpNext} XP to level {xp.level + 1}
              </div>
            </div>
            <div class="w-72 max-w-full h-2.5 rounded-full bg-base-300 overflow-hidden mt-3">
              <div
                class="h-full rounded-full bg-primary transition-all"
                style="width:{Math.min(100, (xp.xpIn / Math.max(1, xp.xpNext)) * 100)}%"
              ></div>
            </div>
          </div>
          <div class="text-3xl font-bold text-primary">{xp.xp} <span class="text-sm font-normal opacity-50">XP</span></div>
        </div>

        <div class="mt-5">
          <h3 class="text-xs font-semibold uppercase tracking-wider opacity-50 mb-3">Badges</h3>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3">
            {#each xp.badges as b}
              <div
                class="rounded-xl px-3 py-3 text-center {b.earned ? 'bg-base-200 border border-primary/30' : 'bg-base-200/50 opacity-45 border border-transparent'}"
                title={`${b.name} — ${b.desc}`}
              >
                <div class="text-2xl">{b.earned ? b.icon : "🔒"}</div>
                <div class="text-xs font-medium mt-1 truncate">{b.name}</div>
                <div class="text-[10px] opacity-50 truncate">{b.earned ? (b.unlockedAt ? new Date(b.unlockedAt).toLocaleDateString() : "earned") : "locked"}</div>
              </div>
            {/each}
          </div>
        </div>
      </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-6">
      <!-- skills -->
      <div class="card bg-base-200 rounded-2xl p-6">
        <h2 class="font-semibold mb-4">Skills</h2>
        <div class="space-y-3">
          {#each prog.perModule as mb}
            {@const mod = cur.modules.find((m) => m.id === mb.id)}
            <div>
              <div class="flex items-center justify-between text-sm mb-1">
                <span class="truncate">
                  <span class="opacity-50 font-mono text-xs mr-2">{mb.id}</span>
                  <span class="opacity-80">{mod?.skill ?? mod?.title ?? mb.id}</span>
                </span>
                <span class="text-xs opacity-50">{mb.done}/{mb.total}</span>
              </div>
              <div class="h-2 rounded-full bg-base-300 overflow-hidden">
                <div
                  class="h-full rounded-full transition-all {mb.done === mb.total ? 'bg-success' : 'bg-primary'}"
                  style="width:{mb.total ? (mb.done / mb.total) * 100 : 0}%"
                ></div>
              </div>
            </div>
          {/each}
        </div>
      </div>

      <div class="space-y-6">
        <!-- settings -->
        <div class="card bg-base-200 rounded-2xl p-6">
          <h2 class="font-semibold mb-4">Preferences</h2>
          <div class="form-control">
            <label class="label"><span class="label-text text-xs">Theme</span></label>
            <select class="select select-sm select-bordered" bind:value={theme}>
              {#each themes as t}
                <option value={t}>{t}</option>
              {/each}
            </select>
          </div>
          <div class="grid grid-cols-3 gap-3 mt-4">
            <div class="form-control">
              <label class="label"><span class="label-text text-xs">Work (min)</span></label>
              <input class="input input-sm input-bordered" type="number" min="1" max="90" bind:value={workMin} />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text text-xs">Break (min)</span></label>
              <input class="input input-sm input-bordered" type="number" min="1" max="30" bind:value={breakMin} />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text text-xs">Review/session</span></label>
              <input class="input input-sm input-bordered" type="number" min="1" max="200" bind:value={reviewLimit} />
            </div>
          </div>
          <div class="flex items-center gap-3 mt-4">
            <button class="btn btn-sm btn-primary" onclick={saveSettings}>Save settings</button>
            {#if savedMsg}
              <span class="text-xs text-success">{savedMsg}</span>
            {/if}
          </div>
        </div>

        <!-- time -->
        <div class="card bg-base-200 rounded-2xl p-6">
          <h2 class="font-semibold mb-3">Most studied exercises</h2>
          {#if topExercises().length === 0}
            <p class="text-sm opacity-50">Nothing tracked yet — start a focus session.</p>
          {/if}
          <div class="space-y-2">
            {#each topExercises() as t}
              <div class="flex items-center gap-3 text-sm">
                <span class="font-mono text-xs opacity-50 shrink-0">{t.id}</span>
                <span class="truncate flex-1">{t.title}</span>
                <span class="text-xs opacity-60">{fmtMins(t.minutes)}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
  {:else if !err}
    <div class="flex justify-center py-20">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>
  {/if}
</div>