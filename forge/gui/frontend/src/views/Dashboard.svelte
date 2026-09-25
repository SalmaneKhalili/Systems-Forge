<script lang="ts">
  import { onMount } from "svelte";
  import { api, fmtDate, fmtMins } from "../lib/api";
  import { goto } from "../lib/router.svelte";
  import { surprise } from "../lib/quests";
  import { confetti, ding } from "../lib/celebrate";
  import type { Curriculum, DailyPlan, Progress } from "../lib/types";
  import ProgressRing from "../components/ProgressRing.svelte";
  import Heatmap from "../components/Heatmap.svelte";

  let prog = $state<Progress | null>(null);
  let cur = $state<Curriculum | null>(null);
  let plan = $state<DailyPlan | null>(null);
  let questCeleb = $state(false);
  let err = $state("");

  const weekdays = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

  onMount(async () => {
    try {
      [prog, cur, plan] = await Promise.all([
        api.progress(),
        api.curriculum(),
        api.dailyPlan(),
      ]);
    } catch (e) {
      err = String(e);
    }
  });

  // Celebrate once (per mount) when every quest is done.
  $effect(() => {
    if (plan?.doneAll && !questCeleb) {
      questCeleb = true;
      confetti(160);
      ding();
    }
  });

  const continueEx = $derived.by(() => {
    if (!cur) return null;
    for (const m of cur.modules) {
      for (const ex of m.exercises) {
        if (ex.status !== "pass") return ex;
      }
    }
    return null;
  });

  const completedModuleCount = $derived(
    cur ? cur.modules.filter((m) => m.done === m.total && m.total > 0).length : 0,
  );

  function dayLabel(iso: string): string {
    return new Date(iso).toLocaleDateString(undefined, { weekday: "short" });
  }
</script>

<div class="p-6 lg:p-8 max-w-6xl mx-auto">
  {#if err}
    <div class="alert alert-error">{err}</div>
  {/if}

  <!-- hero -->
  <div class="flex flex-wrap items-center justify-between gap-6">
    <div>
      <h1 class="text-3xl font-bold">Dashboard</h1>
      <p class="opacity-60 mt-1">
        {new Date().toLocaleDateString(undefined, { weekday: "long", month: "long", day: "numeric" })} —
        {#if prog}you've done <b class="text-primary">{prog.todayRuns}</b> grading run{prog.todayRuns === 1 ? "" : "s"} today.{/if}
      </p>
    </div>
    <div class="flex gap-2">
      {#if prog}
        <div class="stat bg-base-200 rounded-xl px-5 py-3">
          <div class="stat-title text-xs">Streak</div>
          <div class="stat-value text-2xl flex items-center gap-2">
            🔥 {prog.streak}
            <span class="text-sm font-normal opacity-50">days</span>
          </div>
        </div>
        <div class="stat bg-base-200 rounded-xl px-5 py-3">
          <div class="stat-title text-xs">Focus this week</div>
          <div class="stat-value text-2xl">{fmtMins(prog.weekMin)}</div>
        </div>
      {/if}
    </div>
  </div>

  {#if plan}
    <!-- today: daily quests -->
    <div class="card bg-gradient-to-br from-primary/10 to-secondary/10 border border-primary/20 rounded-2xl p-6 mt-8">
      <div class="flex items-center justify-between mb-4 flex-wrap gap-3">
        <div>
          <h2 class="font-semibold text-lg">Today {plan.doneAll ? "🎉" : ""}</h2>
          <p class="text-xs opacity-60">Your daily quests · {plan.date}</p>
        </div>
        <div class="flex items-center gap-2">
          {#if plan.dueToday > 0}
            <button class="btn btn-sm btn-outline" onclick={() => goto("review")}>
              🃏 {plan.dueToday} card{plan.dueToday === 1 ? "" : "s"} due
            </button>
          {/if}
          <button class="btn btn-sm btn-primary" onclick={() => surprise(cur)}>🎲 Surprise me</button>
        </div>
      </div>
      <div class="grid sm:grid-cols-3 gap-3">
        {#each plan.quests as q}
          {@const done = q.done >= q.target}
          <div class="flex items-center gap-3 bg-base-200/70 rounded-xl px-4 py-3">
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-sm shrink-0 {done ? 'bg-success/20 text-success' : 'bg-primary/15 text-primary'}"
            >
              {done ? "✓" : q.key === "solve" ? "⚔️" : q.key === "review" ? "🃏" : "🎯"}
            </div>
            <div class="min-w-0 flex-1">
              <div class="text-sm font-medium truncate">{q.label}</div>
              <div class="flex items-center gap-2 mt-1">
                <div class="flex-1 h-1.5 rounded-full bg-base-300 overflow-hidden">
                  <div
                    class="h-full rounded-full {done ? 'bg-success' : 'bg-primary'} transition-all"
                    style="width:{Math.min(100, (q.done / Math.max(1, q.target)) * 100)}%"
                  ></div>
                </div>
                <span class="text-xs opacity-60 shrink-0">{q.done}/{q.target} {q.unit}</span>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  {#if prog && cur}
    <!-- main grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-8">
      <!-- progress ring + continue -->
      <div class="card bg-base-200 rounded-2xl p-6 flex flex-col items-center justify-center gap-4">
        <ProgressRing value={prog.passed} total={prog.total} size={150} stroke={12}>
          <div class="text-center">
            <div class="text-3xl font-bold">{prog.passed}<span class="text-lg opacity-50">/{prog.total}</span></div>
            <div class="text-xs opacity-60">exercises passed</div>
          </div>
        </ProgressRing>
        {#if continueEx}
          <button class="btn btn-primary w-full" onclick={() => goto("exercise", continueEx.key)}>
            Continue: {continueEx.title}
          </button>
        {:else}
          <div class="badge badge-success badge-lg gap-1">🎉 Curriculum complete</div>
        {/if}
      </div>

      <!-- activity heatmap -->
      <div class="card bg-base-200 rounded-2xl p-6 lg:col-span-2">
        <div class="flex items-center justify-between mb-3">
          <h2 class="font-semibold">Activity</h2>
          <div class="flex items-center gap-1 text-xs opacity-60">
            <span>less</span>
            <div class="w-3 h-3 rounded-sm bg-base-300"></div>
            <div class="w-3 h-3 rounded-sm bg-primary/40"></div>
            <div class="w-3 h-3 rounded-sm bg-primary/70"></div>
            <div class="w-3 h-3 rounded-sm bg-primary"></div>
            <div class="w-3 h-3 rounded-sm bg-primary-content"></div>
            <span>more</span>
          </div>
        </div>
        <Heatmap days={prog.activity} />
        <div class="flex justify-end text-xs opacity-50 mt-2">
          runs per day · last 90 days
        </div>

        <!-- weekly minutes -->
        <div class="mt-6">
          <h3 class="font-semibold text-sm mb-2">Focus minutes</h3>
          <div class="flex items-end gap-2 h-20">
            {#each prog.minutes as d, i}
              {@const max = Math.max(1, ...prog!.minutes.map((x) => x.count))}
              <div class="flex-1 flex flex-col items-center gap-1">
                <div class="w-full rounded-t-md bg-primary/70" style="height:{Math.max(4, (d.count / max) * 72)}px" title={`${d.count} min`}></div>
                <span class="text-[10px] opacity-50">{d.day.slice(8)}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>

    <!-- per-module -->
    <div class="card bg-base-200 rounded-2xl p-6 mt-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="font-semibold">Modules</h2>
        <button class="btn btn-sm btn-ghost" onclick={() => goto("path")}>Open path →</button>
      </div>
      <div class="space-y-3">
        {#each prog.perModule as mb}
          {@const mod = cur.modules.find((m) => m.id === mb.id)}
          <div>
            <div class="flex items-center justify-between text-sm mb-1">
              <span class="opacity-80 truncate">{mod?.title ?? mb.id}</span>
              <span class="opacity-50 text-xs">{mb.done}/{mb.total}</span>
            </div>
            <div class="h-2 rounded-full bg-base-300 overflow-hidden">
              <div
                class="h-full rounded-full {mb.done === mb.total ? 'bg-success' : 'bg-primary'} transition-all"
                style="width:{mb.total ? (mb.done / mb.total) * 100 : 0}%"
              ></div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {:else if !err}
    <div class="flex justify-center py-20">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>
  {/if}
</div>