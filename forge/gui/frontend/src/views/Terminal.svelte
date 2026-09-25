<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { Terminal } from "@xterm/xterm";
  import { FitAddon } from "@xterm/addon-fit";
  import "@xterm/xterm/css/xterm.css";
  import { api, logError } from "../lib/api";
  import type { WorkDirInfo } from "../lib/types";

  const PROBE = "__sf_term_ok"; // echoed back by the interactive shell

  let dirs = $state<WorkDirInfo[]>([]);
  let selectedKey = $state("");
  let err = $state("");
  let running = $state(false);
  let sessionKind = $state<"shell" | "nvim" | null>(null);
  let pid = $state("");
  let container: HTMLDivElement;

  // Round-trip probe: after starting, one echoed marker proves keystrokes are
  // reaching the pty (shell echo would not appear otherwise).
  let connected = $state<"unknown" | "ok" | "fail">("unknown");
  let probePending = $state(false);
  let probeTimer: ReturnType<typeof setTimeout> | undefined;
  let outBuf = "";

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
      if (probePending) {
        outBuf += bin;
        if (outBuf.length > 8192) outBuf = outBuf.slice(-4096);
        if (outBuf.includes(PROBE)) {
          probePending = false;
          connected = "ok";
          if (probeTimer) {
            clearTimeout(probeTimer);
            probeTimer = undefined;
          }
        }
      }
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
        sessionKind = null;
        probePending = false;
        connected = "unknown";
        if (probeTimer) {
          clearTimeout(probeTimer);
          probeTimer = undefined;
        }
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
      logError("WorkDirs", e);
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
      sessionKind = kind;
      if (term) void api.termResize(pid, term.cols, term.rows).catch(() => {});
      // Verify the keystroke round trip once per session.
      connected = "unknown";
      probePending = true;
      outBuf = "";
      void api.termInput(id, `echo ${PROBE}\r`).catch(() => {});
      probeTimer = setTimeout(() => {
        if (probePending) {
          probePending = false;
          connected = "fail";
        }
      }, 6000);
      term.focus();
    } catch (e) {
      logError("TermStart", e);
      unregisterSession(id);
      err = String(e);
    }
  }

  function stopSession() {
    if (pid) {
      const id = pid;
      pid = "";
      running = false;
      sessionKind = null;
      connected = "unknown";
      probePending = false;
      if (probeTimer) {
        clearTimeout(probeTimer);
        probeTimer = undefined;
      }
      void api.termStop(id).catch(() => {});
      unregisterSession(id);
    }
  }

  function connPill() {
    if (connected === "ok")
      return { cls: "badge badge-success badge-sm gap-1", txt: "connected" };
    if (connected === "fail")
      return {
        cls: "badge badge-warning badge-sm gap-1",
        txt: "no echo — typing is not reaching the shell",
      };
    return { cls: "badge badge-ghost badge-sm gap-1", txt: "idle" };
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
      <select class="select select-sm select-bordered max-w-56" bind:value={selectedKey} title={selectedDir()}>
        {#each dirs as d}
          <option value={d.key}>{d.title}</option>
        {/each}
      </select>
      <button class="btn btn-sm btn-ghost" onclick={() => startSession("shell")} disabled={!selectedDir() && dirs.length === 0}>
        ⟩ Shell
      </button>
      <button class="btn btn-sm btn-primary" onclick={() => startSession("nvim")} disabled={!selectedDir() && dirs.length === 0}>
        nvim
      </button>
    </div>
  </div>

  {#if err}
    <div class="alert alert-error mb-3 py-2">{err}</div>
  {/if}

  <div class="flex-1 min-h-0 rounded-xl overflow-hidden border border-base-300 flex flex-col">
    <!-- session strip -->
    <div class="flex items-center gap-2 px-3 py-2 border-b border-base-300 bg-base-200/60 shrink-0">
      <span class={`w-2 h-2 rounded-full ${running ? "bg-success animate-pulse" : "bg-base-300"}`}></span>
      <span class="text-xs font-semibold">{running ? (sessionKind === "nvim" ? "nvim" : "shell") : "not running"}</span>
      <span class="text-[11px] opacity-50 font-mono truncate hidden md:inline max-w-64" title={selectedDir()}>
        {selectedDir() || "pick a working directory"}
      </span>
      <div class="ml-auto flex items-center gap-2">
        <span class={connPill().cls} title="A marker is echoed back through the shell to confirm keystrokes reach it.">
          <span class="w-1.5 h-1.5 rounded-full {connected === 'ok' ? 'bg-success' : connected === 'fail' ? 'bg-warning' : 'bg-base-300'}"></span>
          {connPill().txt}
        </span>
        {#if running}
          <button class="btn btn-xs btn-ghost text-error" onclick={stopSession} title="Kill the process and close the session">
            Stop
          </button>
        {/if}
      </div>
    </div>

    <!-- terminal surface -->
    <div class="flex-1 min-h-0 bg-[#0f1115] p-2">
      <div bind:this={container} class="h-full w-full"></div>
    </div>
  </div>
</div>