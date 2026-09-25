<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "../lib/api";
  import { goto } from "../lib/router.svelte";
  import type { Curriculum, ModuleInfo } from "../lib/types";

  let cur = $state<Curriculum | null>(null);
  let expanded = $state<Set<string>>(new Set());
  let err = $state("");

  onMount(async () => {
    try {
      cur = await api.curriculum();
      // Default: expand modules that still have work, collapse completed ones.
      const open = new Set<string>();
      for (const m of cur.modules) {
        if (m.done < m.total) open.add(m.id);
      }
      expanded = open;
    } catch (e) {
      err = String(e);
    }
  });

  function toggle(m: ModuleInfo) {
    const next = new Set(expanded);
    if (next.has(m.id)) next.delete(m.id);
    else next.add(m.id);
    expanded = next;
  }

  function dot(status: string): string {
    switch (status) {
      case "pass":
        return "bg-success";
      case "fail":
        return "bg-warning";
      default:
        return "bg-base-400 bg-opacity-40";
    }
  }
</script>

<div class="p-6 lg:p-8 max-w-4xl mx-auto">
  <h1 class="text-3xl font-bold">Path</h1>
  <p class="opacity-60 mt-1">18 modules · {cur?.total ?? "—"} exercises.</p>

  {#if err}
    <div class="alert alert-error mt-4">{err}</div>
  {/if}

  {#if cur}
    <!-- overall progress -->
    <div class="card bg-base-200 rounded-2xl p-5 mt-6">
      <div class="flex items-center justify-between text-sm mb-2">
        <span class="font-semibold">Curriculum</span>
        <span class="opacity-60">
          {cur.modules.filter((m) => m.done === m.total && m.total > 0).length}/{cur.modules.length} modules complete
          · {cur.modules.reduce((s, m) => s + m.done, 0)}/{cur.total} passed
        </span>
      </div>
      <div class="h-3 rounded-full bg-base-300 overflow-hidden">
        <div
          class="h-full rounded-full bg-primary transition-all"
          style="width:{cur.total ? (cur.modules.reduce((s, m) => s + m.done, 0) / cur.total) * 100 : 0}%"
        ></div>
      </div>
    </div>

    <!-- module list -->
    <div class="mt-8 space-y-3">
      {#each cur.modules as mod, mi}
        {@const open = expanded.has(mod.id)}
        <div class="card bg-base-200 rounded-2xl overflow-hidden">
          <button
            class="w-full text-left px-5 py-4 flex items-center gap-4 hover:bg-base-300/40 transition-colors"
            onclick={() => toggle(mod)}
          >
            <span class="flex items-center gap-2">
              <span class="badge badge-outline badge-sm">{mi}</span>
              <span class="font-semibold">{mod.title}</span>
            </span>
            <span class="text-xs opacity-50 hidden sm:inline flex-1">{mod.skill}</span>
            <span class="text-xs opacity-60">{mod.done}/{mod.total}</span>
            <div class="w-24 h-2 rounded-full bg-base-300 overflow-hidden hidden sm:block">
              <div
                class="h-full {mod.done === mod.total ? 'bg-success' : 'bg-primary'}"
                style="width:{mod.total ? (mod.done / mod.total) * 100 : 0}%"
              ></div>
            </div>
            <svg
              class="w-4 h-4 opacity-40 transition-transform {open ? 'rotate-180' : ''}"
              fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          {#if open}
            <div class="border-t border-base-300/60 px-2 py-2">
              {#each mod.exercises as ex}
                <button
                  class="w-full text-left flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-300/50 transition-colors"
                  onclick={() => goto("exercise", ex.key)}
                >
                  <span class="w-2.5 h-2.5 rounded-full shrink-0 {dot(ex.status)}"></span>
                  <span class="text-sm flex-1">
                    <span class="text-xs opacity-50 mr-2 font-mono">{ex.id}</span>
                    {ex.title}
                  </span>
                  <span class="text-[10px] uppercase tracking-wide opacity-40 hidden sm:inline">{ex.lang}</span>
                  {#if ex.isGate}
                    <span class="badge badge-primary badge-sm">gate</span>
                  {/if}
                  {#if ex.status === "pass"}
                    <svg class="w-4 h-4 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                    </svg>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {:else if !err}
    <div class="flex justify-center py-20">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>
  {/if}
</div>