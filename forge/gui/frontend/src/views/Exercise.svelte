<script lang="ts">
  import { onMount } from "svelte";
  import { api, fmtDuration } from "../lib/api";
  import { confetti, ding } from "../lib/celebrate";
  import type { ExerciseDetail, CheckResult, FileInfo } from "../lib/types";
  import Markdown from "../components/Markdown.svelte";
  import Coder from "../components/Coder.svelte";

  let { id }: { id: string } = $props();

  let ex = $state<ExerciseDetail | null>(null);
  let tab = $state<"spec" | "readings" | "qa" | "history" | "hints" | "notes">("spec");
  let files = $state<FileInfo[]>([]);
  let currentFile = $state<string | null>(null);
  let content = $state("");
  let saving = $state(false);
  let savedFlash = $state("");
  let grading = $state(false);
  let result = $state<CheckResult | null>(null);
  let lastErr = $state("");
  let newFileRel = $state("");
  let hints = $state<import("../lib/types").HintsResult | null>(null);
  let revealed = $state<Set<number>>(new Set());
  let notes = $state("");
  let notesLoaded = $state(false);
  let notesSaving = $state(false);
  let notesSaved = $state("");
  let notesPreview = $state(false);

  onMount(load);

  async function load() {
    try {
      ex = await api.exercise(id);
      files = ex.files;
      if (!currentFile && files.length > 0 && !files[0].dir) {
        await openFile(files[0].path);
      }
    } catch (e) {
      lastErr = String(e);
    }
  }

  async function openFile(rel: string) {
    try {
      content = await api.readFile(id, rel);
      currentFile = rel;
    } catch (e) {
      lastErr = String(e);
    }
  }

  async function save() {
    if (!currentFile) return;
    saving = true;
    savedFlash = "";
    try {
      await api.writeFile(id, currentFile, content);
      savedFlash = "Saved";
      setTimeout(() => (savedFlash = ""), 1500);
      files = await api.files(id);
      if (ex) ex.files = files;
    } catch (e) {
      lastErr = String(e);
    } finally {
      saving = false;
    }
  }

  async function newFile() {
    const rel = newFileRel.trim().replace(/^\/+/, "");
    if (!rel) return;
    newFileRel = "";
    currentFile = rel;
    content = "";
    try {
      await api.writeFile(id, rel, content);
      files = await api.files(id);
      if (ex) ex.files = files;
    } catch (e) {
      lastErr = String(e);
    }
  }

  async function grade() {
    grading = true;
    result = null;
    lastErr = "";
    try {
      result = await api.runCheck(id);
      ex = await api.exercise(id);
      files = ex.files;
      if (result.pass) {
        confetti(120);
        ding();
      }
    } catch (e) {
      lastErr = String(e);
    } finally {
      grading = false;
    }
  }

  function langFor(rel: string): string {
    const ext = rel.split(".").pop() ?? "";
    switch (ext) {
      case "c":
      case "h":
        return "c";
      case "cpp":
      case "cc":
      case "hpp":
        return "cpp";
      case "go":
        return "go";
      case "py":
        return "python";
      case "js":
      case "ts":
        return "javascript";
      case "sh":
      case "bash":
        return "sh";
      case "mk":
        return "sh";
      default:
        return "c";
    }
  }

  function statusPill(s: string) {
    if (s === "pass") return "badge badge-success gap-1";
    if (s === "fail") return "badge badge-warning gap-1";
    return "badge badge-ghost gap-1";
  }

  function statusLabel(s: string) {
    if (s === "pass") return "Passed";
    if (s === "fail") return "In progress";
    return "Not started";
  }

  async function ensureHints() {
    if (hints) return;
    try {
      hints = await api.hints(id);
    } catch (e) {
      lastErr = String(e);
    }
  }

  function revealHint(i: number) {
    if (!hints || hints.fails < i) return;
    const n = new Set(revealed);
    if (n.has(i)) n.delete(i);
    else n.add(i);
    revealed = n;
  }

  async function ensureNotes() {
    if (notesLoaded) return;
    notesLoaded = true;
    try {
      notes = await api.readNotes(id);
    } catch (e) {
      lastErr = String(e);
    }
  }

  async function saveNotes() {
    notesSaving = true;
    notesSaved = "";
    try {
      await api.writeNotes(id, notes);
      notesSaved = "Saved";
      setTimeout(() => (notesSaved = ""), 1500);
    } catch (e) {
      lastErr = String(e);
    } finally {
      notesSaving = false;
    }
  }

  async function rerunMethod(i: number) {
    if (!result || grading) return;
    grading = true;
    try {
      const r = await api.runCheckMethod(id, i);
      if (!result) return;
      if (r.methods[0]) result.methods[i] = r.methods[0];
      result.pass = result.methods.every((x) => x.pass);
      result.durationMs = r.methods[0]?.durMs ?? result.durationMs;
    } catch (e) {
      lastErr = String(e);
    } finally {
      grading = false;
    }
  }
</script>

<div class="h-full flex flex-col">
  {#if ex}
    <!-- header -->
    <div class="px-6 py-4 border-b border-base-300 flex items-center gap-4 flex-wrap">
      <div>
        <div class="text-xs opacity-50 font-mono mb-0.5">{ex.module} / {ex.id}</div>
        <h1 class="text-xl font-bold leading-tight">{ex.title}</h1>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <span class="badge badge-outline badge-sm">{ex.lang}</span>
        <span class={statusPill(ex.status)}>
          {#if ex.status === "pass"}✓{/if}{statusLabel(ex.status)}
        </span>
      </div>
    </div>

    <div class="flex-1 min-h-0 grid grid-cols-1 lg:grid-cols-5">
      <!-- left: spec -->
      <div class="lg:col-span-3 min-h-0 flex flex-col border-r border-base-300">
        <div class="tabs px-4 pt-3 gap-1">
          <button class="tab {tab === 'spec' ? 'tab-active' : ''}" onclick={() => (tab = "spec")}>Spec</button>
          <button class="tab {tab === 'readings' ? 'tab-active' : ''}" onclick={() => (tab = "readings")}>
            Readings ({ex.readings.length})
          </button>
          <button class="tab {tab === 'qa' ? 'tab-active' : ''}" onclick={() => (tab = "qa")}>
            Q&A ({ex.qa.length})
          </button>
          <button class="tab {tab === 'history' ? 'tab-active' : ''}" onclick={() => (tab = "history")}>
            History ({ex.history.length})
          </button>
          <button class="tab {tab === 'hints' ? 'tab-active' : ''}" onclick={() => { tab = "hints"; ensureHints(); }}>
            Hints
          </button>
          <button class="tab {tab === 'notes' ? 'tab-active' : ''}" onclick={() => { tab = "notes"; ensureNotes(); }}>
            Notes
          </button>
        </div>

        <div class="flex-1 overflow-y-auto px-6 py-4">
          {#if tab === "spec"}
            {#if ex.spec}
              <Markdown source={ex.spec} />
            {:else}
              <p class="opacity-50 text-sm">No spec markdown found for this exercise.</p>
            {/if}
          {:else if tab === "readings"}
            {#if ex.readings.length === 0}
              <p class="opacity-50 text-sm">No readings for this exercise.</p>
            {/if}
            <div class="space-y-3">
              {#each ex.readings as r}
                <div class="card bg-base-200 rounded-xl p-4">
                  <div class="font-semibold text-sm">{r.source}</div>
                  <div class="text-xs opacity-60 mt-1">
                    {r.author}{r.chapter ? ` · ${r.chapter}` : ""}{r.section ? ` · ${r.section}` : ""}
                  </div>
                  {#if r.url}
                    <a class="text-xs text-primary mt-2 inline-block" href={r.url} target="_blank" rel="noreferrer">
                      {r.url} ↗
                    </a>
                  {/if}
                </div>
              {/each}
            </div>
          {:else if tab === "qa"}
            {#if ex.qa.length === 0}
              <p class="opacity-50 text-sm">No Q&A metadata for this exercise.</p>
            {/if}
            <div class="space-y-2">
              {#each ex.qa as qa, i}
                <details class="collapse collapse-arrow bg-base-200 rounded-xl">
                  <summary class="collapse-title text-sm font-medium py-3">
                    <span class="text-xs opacity-40 font-mono mr-2">{i + 1}.</span>{qa.q}
                  </summary>
                  <div class="collapse-content text-sm opacity-80">{qa.a}</div>
                </details>
              {/each}
            </div>
          {:else if tab === "hints"}
            {#if !hints}
              <p class="opacity-50 text-sm">Loading hints…</p>
            {:else if hints.hints.length === 0}
              <p class="opacity-50 text-sm">No hints available for this exercise.</p>
            {:else}
              <div class="space-y-3">
                <p class="text-xs opacity-50">
                  Staged hints — each unlocks as you keep trying. You've had {hints.fails} failed run{hints.fails === 1 ? "" : "s"}.
                </p>
                {#each hints.hints as h, i}
                  {@const unlocked = hints.fails >= i}
                  {@const open = revealed.has(i)}
                  <div class="card bg-base-200 rounded-xl p-4">
                    <div class="flex items-center justify-between gap-3 mb-2">
                      <span class="text-xs font-semibold uppercase tracking-wider opacity-60">Hint {i + 1}</span>
                      {#if unlocked}
                        <button class="btn btn-xs btn-ghost" onclick={() => revealHint(i)}>
                          {open ? "Hide" : "Reveal"}
                        </button>
                      {:else}
                        <span class="badge badge-ghost badge-sm text-[10px]">unlocks after {i} failed run{i === 1 ? "" : "s"}</span>
                      {/if}
                    </div>
                    {#if open}
                      <div class="text-sm whitespace-pre-wrap opacity-85">{h}</div>
                    {:else}
                      <div class="text-sm opacity-25 select-none" aria-hidden="true">
                        {#if unlocked}Click reveal to read this hint.{:else}Keep failing the grader to unlock…{/if}
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
          {:else if tab === "notes"}
            <div class="flex flex-col h-full gap-3">
              <div class="flex items-center gap-2">
                <button class="btn btn-xs btn-ghost" onclick={() => (notesPreview = !notesPreview)}>
                  {notesPreview ? "✏️ Edit" : "👁 Preview"}
                </button>
                <button class="btn btn-xs btn-primary" onclick={saveNotes} disabled={notesSaving}>
                  {#if notesSaving}<span class="loading loading-spinner loading-xs"></span>{:else}Save{/if}
                </button>
                {#if notesSaved}
                  <span class="text-xs text-success">{notesSaved}</span>
                {/if}
              </div>
              {#if notesPreview}
                <div class="flex-1 overflow-y-auto card bg-base-200 rounded-xl p-4">
                  <Markdown source={notes || "_Nothing yet — capture approach notes, mistakes, and insights here._"} />
                </div>
              {:else}
                <textarea
                  class="textarea textarea-bordered bg-base-200 flex-1 font-mono text-sm leading-relaxed resize-none"
                  placeholder="Private notes for this exercise — saved to your answers workspace as notes.md"
                  bind:value={notes}
                ></textarea>
              {/if}
            </div>
          {:else}
            {#if ex.history.length === 0}
              <p class="opacity-50 text-sm">No runs yet — write some code and grade it.</p>
            {/if}
            <div class="space-y-2">
              {#each ex.history as run}
                <div class="card bg-base-200 rounded-xl px-4 py-3 flex items-center gap-3">
                  <span class={`w-2 h-2 rounded-full ${run.status === "pass" ? "bg-success" : run.status === "fail" ? "bg-warning" : "bg-error"}`}></span>
                  <span class="text-sm font-medium">{run.status === "pass" ? "Passed" : "Failed"}</span>
                  <span class="text-xs opacity-50 ml-auto">{fmtDuration(run.durationMs)}</span>
                  <span class="text-xs opacity-50">{new Date(run.at).toLocaleString()}</span>
                </div>
              {/each}
            </div>
          {/if}

          <!-- grader result -->
          {#if result}
            <div class="mt-6">
              <div class={`alert ${result.pass ? "alert-success" : "alert-error"} rounded-xl`}>
                <span class="font-semibold">
                  {result.pass ? "🎉 All checks passed!" : "✗ Some checks failed"}
                </span>
                <span class="text-xs opacity-70">{fmtDuration(result.durationMs)}</span>
              </div>

              <div class="space-y-3 mt-4">
                {#each result.methods as m, i}
                  <div class="border border-base-300 rounded-xl overflow-hidden">
                    <div class="flex items-center gap-3 px-4 py-2.5 bg-base-200">
                      <span class={`w-2 h-2 rounded-full ${m.pass ? "bg-success" : "bg-error"}`}></span>
                      <span class="text-sm font-medium">{m.label}</span>
                      <span class="text-[10px] px-2 py-0.5 rounded bg-base-300 opacity-70">{m.type}</span>
                      <span class="text-xs opacity-50 ml-auto">{fmtDuration(m.durMs)}</span>
                      <button
                        class="btn btn-xs btn-ghost"
                title={"Re-run this method only (no history record)"}
                        disabled={grading}
                        onclick={() => rerunMethod(i)}
                      >
                        ↻
                      </button>
                    </div>
                    {#if m.parts.length}
                      <div class="px-4 py-3 space-y-1.5">
                        {#each m.parts as p}
                          <div class="flex items-start gap-2 text-sm">
                            {#if p.pass}
                              <span class="text-success font-bold">✓</span>
                            {:else}
                              <span class="text-error font-bold">✗</span>
                            {/if}
                            <div class="min-w-0">
                              <span class="font-mono text-xs opacity-80">{p.name}</span>
                              {#if p.detail}
                                <div class="text-xs opacity-60 whitespace-pre-wrap">{p.detail}</div>
                              {/if}
                            </div>
                          </div>
                        {/each}
                      </div>
                    {/if}
                    {#if m.output}
                      <pre class="text-xs overflow-x-auto px-4 pb-3 opacity-70 bg-base-200/50 mx-3 mb-3 rounded-lg p-3">{m.output}</pre>
                    {/if}
                  </div>
                {/each}
              </div>

              {#if result.output && !result.pass}
                <pre class="text-xs overflow-x-auto mt-4 border border-dashed border-error/40 rounded-xl p-3 opacity-80">{result.output}</pre>
              {/if}
            </div>
          {/if}
        </div>
      </div>

      <!-- right: workbench -->
      <div class="lg:col-span-2 min-h-0 flex flex-col">
        <div class="flex items-center gap-2 px-4 py-3 border-b border-base-300">
          <button class="btn btn-primary btn-sm" onclick={grade} disabled={grading}>
            {#if grading}
              <span class="loading loading-spinner loading-xs"></span> Grading…
            {:else}
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Run grader
            {/if}
          </button>
          <button class="btn btn-sm" onclick={save} disabled={saving || !currentFile}>
            {#if saving}<span class="loading loading-spinner loading-xs"></span>{:else}Save{/if}
          </button>
          {#if savedFlash}
            <span class="text-xs text-success">{savedFlash}</span>
          {/if}
        </div>

        <div class="flex-1 min-h-0 flex">
          <!-- file list -->
          <div class="w-44 shrink-0 border-r border-base-300 bg-base-200/40 overflow-y-auto py-2">
            <div class="px-3 pb-2 text-[10px] uppercase tracking-wider opacity-40 font-semibold">Files</div>
            {#each files.filter((f) => !f.dir) as f}
              <button
                class="w-full text-left px-3 py-1.5 text-xs font-mono truncate hover:bg-base-300/60 transition-colors {currentFile === f.path ? 'bg-primary/20 border-l-2 border-primary' : 'border-l-2 border-transparent'}"
                onclick={() => openFile(f.path)}
                title={f.path}
              >
                {f.name}
              </button>
            {/each}
            <div class="px-3 pt-3">
              <div class="flex gap-1">
                <input
                  class="input input-xs input-bordered w-full font-mono"
                  placeholder="new.txt"
                  bind:value={newFileRel}
                  onkeydown={(e) => e.key === "Enter" && newFile()}
                />
                <button class="btn btn-xs btn-ghost" onclick={newFile} title="Create file">+</button>
              </div>
            </div>
          </div>

          <!-- editor -->
          <div class="flex-1 min-w-0">
            {#if currentFile}
              {#key currentFile}
                <div class="h-full">
                  <div class="h-8 px-3 flex items-center justify-between border-b border-base-300 bg-base-200/40">
                    <span class="text-xs font-mono opacity-70">{currentFile}</span>
                    <span class="text-[10px] uppercase opacity-40">{ex.lang}</span>
                  </div>
                  <div class="h-[calc(100%-2rem)]">
                    <Coder value={content} lang={langFor(currentFile)} onChange={(v) => (content = v)} onSave={save} />
                  </div>
                </div>
              {/key}
            {:else}
              <div class="h-full flex items-center justify-center text-sm opacity-40">
                Select a file to start writing
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  {:else if lastErr}
    <div class="p-8"><div class="alert alert-error">{lastErr}</div></div>
  {:else}
    <div class="flex-1 flex items-center justify-center">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>
  {/if}
</div>

{#if lastErr}
  <div class="toast toast-end">
    <div class="alert alert-error">{lastErr}</div>
  </div>
{/if}