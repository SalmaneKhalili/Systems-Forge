<script lang="ts">
  import { onMount } from "svelte";
  import { route, goto } from "./lib/router.svelte";
  import { appState } from "./lib/appstate.svelte";
  import { api } from "./lib/api";
  import type { Route } from "./lib/router.svelte";

  import Dashboard from "./views/Dashboard.svelte";
  import Path from "./views/Path.svelte";
  import Exercise from "./views/Exercise.svelte";
  import Review from "./views/Review.svelte";
  import Focus from "./views/Focus.svelte";
  import Profile from "./views/Profile.svelte";
  import Terminal from "./views/Terminal.svelte";
  import FocusOverlay from "./components/FocusOverlay.svelte";

  let modules = $state(0);
  let total = $state(0);
  let backendErr = $state("");
  let sidebarOpen = $state(true);
  const SIDEBAR_KEY = "ui.sidebar";

  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
    void api
      .saveSettings({ [SIDEBAR_KEY]: sidebarOpen ? "open" : "closed" })
      .catch(() => {});
  }

  const nav = [
    { name: "dashboard", label: "Dashboard", icon: "M3 12l9-9 9 9M5 10v10a1 1 0 001 1h4v-6h4v6h4a1 1 0 001-1V10" },
    { name: "path", label: "Path", icon: "M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l5.447-2.724A1 1 0 0015 16.382V5.618a1 1 0 00-1.447-.894L9 7m0 13V7" },
    { name: "review", label: "Review", icon: "M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" },
    { name: "focus", label: "Focus", icon: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" },
    { name: "terminal", label: "Terminal", icon: "M6 6l4 4-4 4m5-8h7" },
    { name: "profile", label: "Profile", icon: "M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" },
  ];

  onMount(() => {
    appState.load().then(async () => {
      try {
        const c = await api.curriculum();
        modules = c.modules.length;
        total = c.total;
      } catch (e) {
        backendErr = String(e);
      }
    });
    // Restore the persisted sidebar preference and wire Ctrl/Cmd+B.
    void api
      .settings()
      .then((s) => {
        if (s?.[SIDEBAR_KEY] === "closed") sidebarOpen = false;
      })
      .catch(() => {});
    const onKey = (e: KeyboardEvent) => {
      if (!(e.metaKey || e.ctrlKey) || e.key.toLowerCase() !== "b") return;
      // Never steal Ctrl+B from editors/inputs/terminal (tmux prefix etc.).
      const t = e.target as HTMLElement | null;
      if (
        t &&
        (t.tagName === "TEXTAREA" ||
          t.tagName === "INPUT" ||
          t.tagName === "SELECT" ||
          t.isContentEditable)
      )
        return;
      e.preventDefault();
      toggleSidebar();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });

  const r = $derived($route as Route);

  let mainEl: HTMLElement;

  // Each top-level page scrolls inside <main>. When navigating between pages
  // (or between exercises), the container keeps its previous scrollTop, which
  // can strand the user mid-page on a shorter/self-scrolling view. Reset it
  // on every route change.
  $effect(() => {
    void r.name;
    void r.param; // exercise → exercise navigations carry a new param
    if (mainEl) mainEl.scrollTop = 0;
    window.scrollTo(0, 0);
  });
</script>

<FocusOverlay open={appState.focusOpen} onClose={() => appState.setFocusOpen(false)} />

<div class="flex h-full relative">
  <!-- Sidebar (collapsible) -->
  <aside
    class="shrink-0 bg-base-200 flex flex-col overflow-hidden transition-[width] duration-200 ease-in-out
           {sidebarOpen
             ? 'w-16 lg:w-60 border-r border-base-300'
             : 'w-0 lg:w-0 border-r-0'}"
  >
    <div class="px-3 lg:px-5 py-4 flex items-center gap-2.5">
      <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center font-extrabold text-primary-content shadow-lg shadow-primary/30 shrink-0">
        SF
      </div>
      <div class="hidden lg:block leading-tight flex-1 min-w-0">
        <div class="font-bold truncate">systems-forge</div>
        <div class="text-[11px] opacity-50">study app</div>
      </div>
      <button
        class="btn btn-square btn-ghost btn-xs opacity-60 hover:opacity-100"
        onclick={toggleSidebar}
        title="Hide sidebar (Ctrl+B)"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M11 17l-5-5 5-5m7 10l-5-5 5-5" />
        </svg>
      </button>
    </div>

    <nav class="flex-1 px-2 lg:px-3 space-y-1 mt-2">
      {#each nav as n}
        <button
          class="w-full flex items-center gap-3 px-2 py-2 rounded-lg text-sm transition-all
                 {r.name === n.name ? 'bg-primary/15 text-primary font-semibold shadow-inner' : 'hover:bg-base-300/60 text-base-content/70'}"
          onclick={() => goto(n.name)}
          title={n.label}
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d={n.icon} />
          </svg>
          <span class="hidden lg:inline">{n.label}</span>
        </button>
      {/each}
    </nav>

    <div class="p-3 lg:px-5 pb-5 text-xs text-base-content/50 border-t border-base-300/60 hidden lg:block">
      {modules} modules · {total} exercises
      {#if backendErr}
        <div class="text-error mt-1 truncate" title={backendErr}>{backendErr}</div>
      {/if}
    </div>
  </aside>

  <!-- Content -->
  <main class="flex-1 min-w-0 overflow-y-auto bg-base-100" bind:this={mainEl}>
    {#if r.name === "dashboard"}
      <Dashboard />
    {:else if r.name === "path"}
      <Path />
    {:else if r.name === "exercise" && r.param}
      {#key r.param}
        <Exercise id={r.param} />
      {/key}
    {:else if r.name === "review"}
      <Review />
    {:else if r.name === "focus"}
      <Focus />
    {:else if r.name === "terminal"}
      <Terminal />
    {:else if r.name === "profile"}
      <Profile />
    {:else}
      <Dashboard />
    {/if}
  </main>

  {#if !sidebarOpen}
    <button
      class="absolute inset-y-0 left-0 w-5 z-40 flex items-center justify-center cursor-pointer
             bg-gradient-to-r from-base-300/40 to-transparent opacity-50 hover:opacity-100
             transition-opacity group"
      onclick={toggleSidebar}
      title="Show sidebar (Ctrl+B)"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13 17l5-5-5-5M6 17l5-5-5-5" />
      </svg>
    </button>
  {/if}
</div>