<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { Terminal } from "@xterm/xterm";
  import { FitAddon } from "@xterm/addon-fit";
  import "@xterm/xterm/css/xterm.css";
  import { api } from "../lib/api";
  import type { WorkDirInfo } from "../lib/types";

  let dirs = $state<WorkDirInfo[]>([]);
  let selectedKey = $state("");
  let err = $state("");
  let running = $state(false);
  let pid = $state("");
  let container: HTMLDivElement;

  let term: Terminal | null = null;
  let fit: FitAddon | null = null;
  let ro: ResizeObserver | null = null;

  function selectedDir(): string {
    return dirs.find((d) => d.key === selectedKey)?.dir ?? "";
  }

  function onTermOut(b64: string) {
    try {
      const bin = atob(b64);
      const arr = new Uint8Array(bin.length);
      for (let i = 0; i < bin.length; i++) arr[i] = bin.charCodeAt(i);
      term?.write(arr);
    } catch {
      /* ignore malformed chunk */
    }
  }

  function registerSession(id: string) {
    const rt = (window as any).runtime;
    if (rt?.EventsOn) {
      rt.EventsOn(`term:out:${id}`, onTermOut);
      rt.EventsOn(`term:exit:${id}`, () => {
        running = false;
        term?.write("\r\n\x1b[90m[process exited]\x1b[0m\r\n");
      });
    }
  }

  function unregisterSession(id: string) {
    const rt = (window as any).runtime;
    if (rt?.EventsOff) {
      rt.EventsOff(`term:out:${id}`);
      rt.EventsOff(`term:exit:${id}`);
    }
  }

  onMount(async () => {
    try {
      dirs = await api.workDirs();
      if (dirs.length > 0) selectedKey = dirs[0].key;
    } catch (e) {
      err = String(e);
    }
    term = new Terminal({
      cursorBlink: true,
      fontSize: 13,
      fontFamily: "Menlo, Consolas, 'Fira Code', monospace",
      scrollback: 4000,
      theme: { background: "#0f1115", foreground: "#d6d3d1" },
    });
    fit = new FitAddon();
    term.loadAddon(fit);
    term.open(container);
    term.onData((data) => {
      if (pid) void api.termInput(pid, data).catch(() => {});
    });
    ro = new ResizeObserver(() => {
      fit?.fit();
      if (pid && term) void api.termResize(pid, term.cols, term.rows).catch(() => {});
    });
    ro.observe(container);
    fit.fit();
    term.focus();
  });

  onDestroy(() => {
    stopSession();
    ro?.disconnect();
    term?.dispose();
    term = null;
    fit = null;
  });

  async function startSession(kind: "shell" | "nvim") {
    if (!term) return;
    if (!selectedDir()) {
      err = "Pick a working directory first.";
      return;
    }
    if (running) stopSession();
    const id =
      "s" + Date.now().toString(36) + Math.floor(Math.random() * 1e6).toString(36);
    err = "";
    term.reset();
    registerSession(id);
    try {
      await api.termStart(id, selectedDir(), kind === "nvim" ? ["nvim"] : []);
      running = true;
      pid = id;
      if (term) void api.termResize(pid, term.cols, term.rows).catch(() => {});
    } catch (e) {
      unregisterSession(id);
      err = String(e);
    }
  }

  function stopSession() {
    if (pid) {
      const id = pid;
      pid = "";
      running = false;
      void api.termStop(id).catch(() => {});
      unregisterSession(id);
    }
  }
</script>

<div class="p-4 lg:p-6 max-w-6xl mx-auto h-full flex flex-col">
  <div class="flex items-center justify-between gap-3 mb-4 flex-wrap">
    <div>
      <h1 class="text-2xl font-bold">Terminal</h1>
      <p class="text-sm opacity-50 mt-0.5">
        A real shell in your exercise workspaces — run, poke, debug, and open nvim when you want it.
      </p>
    </div>
    <div class="flex items-center gap-2 flex-wrap">
      <select class="select select-sm select-bordered" bind:value={selectedKey}>
        {#each dirs as d}
          <option value={d.key}>{d.title}</option>
        {/each}
      </select>
      <button class="btn btn-sm btn-ghost" onclick={() => startSession("shell")} disabled={!selectedDir()}>
        Shell
      </button>
      <button class="btn btn-sm btn-primary" onclick={() => startSession("nvim")} disabled={!selectedDir()}>
        Open nvim
      </button>
      {#if running}
        <button class="btn btn-sm btn-warning" onclick={stopSession}>Stop</button>
      {/if}
    </div>
  </div>

  {#if err}
    <div class="alert alert-error mb-3 py-2">{err}</div>
  {/if}

  <div class="flex-1 min-h-0 rounded-xl overflow-hidden border border-base-300 bg-[#0f1115]">
    <div bind:this={container} class="h-full w-full"></div>
  </div>
</div>