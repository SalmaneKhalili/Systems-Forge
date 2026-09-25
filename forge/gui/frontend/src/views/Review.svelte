<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "../lib/api";
  import { appState } from "../lib/appstate.svelte";
  import type { Card, UserCardInfo } from "../lib/types";

  let deck = $state<Card[]>([]);
  let moduleFilter = $state("");
  let err = $state("");

  // personal card manager state
  let userCards = $state<UserCardInfo[]>([]);
  let showAdd = $state(false);
  let newQ = $state("");
  let newA = $state("");
  let editingId = $state<number | null>(null);
  let editQ = $state("");
  let editA = $state("");

  // session state
  let sessionActive = $state(false);
  let queue = $state<Card[]>([]);
  let answered = $state<Card[]>([]);
  let current = $state<Card | null>(null);
  let revealed = $state(false);
  let feedback = $state("");

  let modules = $state<string[]>([]);

  const ratings = [
    { v: 0, label: "Again", cls: "btn-error", hint: "relearn" },
    { v: 1, label: "Hard", cls: "btn-warning", hint: "2.4d" },
    { v: 2, label: "Good", cls: "btn-primary", hint: "3.5d" },
    { v: 3, label: "Easy", cls: "btn-success", hint: "5d" },
  ];

  const dueCount = $derived(deck.filter((c) => c.reviewable).length);
  const newCount = $derived(deck.filter((c) => c.new).length);

  onMount(async () => {
    try {
      await reload();
    } catch (e) {
      err = String(e);
    }
  });

  async function reload() {
    deck = await api.quizDeck(moduleFilter);
    modules = [...new Set(deck.map((c) => c.module))].sort();
    userCards = await api.listUserCards();
    if (userCards.length > 0 && !modules.includes("personal")) {
      modules = [...modules, "personal"].sort();
    }
  }

  function ensurePersonal() {
    if (moduleFilter !== "personal") {
      moduleFilter = "personal";
      changeFilter();
    }
  }

  async function addCard() {
    try {
      await api.addCard(newQ, newA);
      newQ = "";
      newA = "";
      showAdd = false;
      await reload();
    } catch (e) {
      err = String(e);
    }
  }

  function startEdit(uc: UserCardInfo) {
    editingId = uc.id;
    editQ = uc.q;
    editA = uc.a;
  }

  async function saveEdit() {
    if (editingId === null) return;
    try {
      await api.updateCard(editingId, editQ, editA);
      editingId = null;
      await reload();
    } catch (e) {
      err = String(e);
    }
  }

  async function deleteCard(uc: UserCardInfo) {
    if (!window.confirm(`Delete personal card?\n\n${uc.q}`)) return;
    try {
      await api.deleteCard(uc.id);
      await reload();
    } catch (e) {
      err = String(e);
    }
  }

  async function changeFilter() {
    sessionActive = false;
    current = null;
    await reload();
  }

  function startSession() {
    const limit = appState.reviewLimit();
    const due = deck.filter((c) => c.reviewable).slice(0, limit);
    const fresh = deck.filter((c) => !c.reviewable && c.new).slice(0, Math.max(0, limit - due.length));
    queue = [...due, ...fresh];
    answered = [];
    sessionActive = true;
    revealed = false;
    nextCard();
  }

  function nextCard() {
    revealed = false;
    feedback = "";
    current = queue.shift() ?? null;
    if (!current) {
      sessionActive = false;
    }
  }

  async function answer(rating: number) {
    if (!current) return;
    const key = current.key;
    try {
      const updated = await api.answerCard(key, rating);
      answered = [updated, ...answered];
      const idx = deck.findIndex((c) => c.key === key);
      if (idx >= 0) deck[idx] = updated;
      feedback = `${updated.reps} reps · ease ${updated.ease.toFixed(2)} · next ${updated.due || "today"}`;
    } catch (e) {
      err = String(e);
    }
    nextCard();
  }

  function endSession() {
    sessionActive = false;
    current = null;
    queue = [];
  }

  function fmtInterval(c: Card): string {
    if (c.new) return "new";
    if (!c.due) return "due now";
    return c.due;
  }

  function onKey(e: KeyboardEvent) {
    if (!sessionActive || !current || !revealed) {
      if (e.key === " ") {
        e.preventDefault();
        if (current) revealed = true;
      }
      return;
    }
    const n = parseInt(e.key, 10);
    if (n >= 1 && n <= 4) answer(n - 1);
  }
</script>

<svelte:window onkeydown={onKey} />

<div class="p-6 lg:p-8 max-w-3xl mx-auto">
  <div class="flex items-center justify-between flex-wrap gap-4">
    <div>
      <h1 class="text-3xl font-bold">Review</h1>
      <p class="opacity-60 mt-1">Spaced repetition · comes from every exercise's Q&A.</p>
    </div>
    <div class="flex items-center gap-2">
      <select class="select select-sm select-bordered" bind:value={moduleFilter} onchange={changeFilter}>
        <option value="">All modules</option>
        {#each modules as m}
          <option value={m}>{m}</option>
        {/each}
      </select>
      {#if sessionActive}
        <button class="btn btn-sm btn-ghost" onclick={endSession}>End session</button>
      {:else}
        <button class="btn btn-sm btn-primary" onclick={startSession} disabled={deck.length === 0}>
          Start session
        </button>
      {/if}
    </div>
  </div>

  {#if err}
    <div class="alert alert-error mt-4">{err}</div>
  {/if}

  {#if !sessionActive}
    <!-- deck summary -->
    <div class="grid grid-cols-3 gap-4 mt-8">
      <div class="card bg-base-200 rounded-2xl p-5 text-center">
        <div class="text-3xl font-bold">{deck.length}</div>
        <div class="text-xs opacity-60 mt-1">total cards</div>
      </div>
      <div class="card bg-base-200 rounded-2xl p-5 text-center">
        <div class="text-3xl font-bold text-primary">{dueCount}</div>
        <div class="text-xs opacity-60 mt-1">due now</div>
      </div>
      <div class="card bg-base-200 rounded-2xl p-5 text-center">
        <div class="text-3xl font-bold text-success">{newCount}</div>
        <div class="text-xs opacity-60 mt-1">new cards</div>
      </div>
    </div>

    {#if deck.length > 0}
      <div class="card bg-base-200 rounded-2xl p-5 mt-4 text-sm opacity-70">
        Sessions review up to <b>{appState.reviewLimit()}</b> cards: cards due today first, then new ones.
        Progress is saved per card — keep a short daily habit.
      </div>
    {:else}
      <div class="card bg-base-100 border border-base-300 rounded-2xl p-8 mt-8 text-center opacity-60">
        No cards in this filter.
      </div>
    {/if}

    <!-- personal card manager -->
    <div class="card bg-base-200 rounded-2xl p-5 mt-6">
      <div class="flex items-center justify-between mb-3">
        <div>
          <h3 class="font-semibold">Personal cards</h3>
          <p class="text-xs opacity-50">Capture mistakes and insights as your own flashcards.</p>
        </div>
        <div class="flex gap-1">
          <button class="btn btn-xs btn-ghost" onclick={() => (showAdd = !showAdd)}>
            {showAdd ? "Cancel" : "+ Add"}
          </button>
          {#if moduleFilter !== "personal" && userCards.length > 0}
            <button class="btn btn-xs btn-outline" onclick={ensurePersonal}>Review personal</button>
          {/if}
        </div>
      </div>

      {#if showAdd}
        <div class="space-y-2 mb-4">
          <input
            class="input input-sm input-bordered w-full"
            placeholder="Question — what did you keep getting wrong?"
            bind:value={newQ}
          />
          <textarea
            class="textarea textarea-bordered w-full text-sm"
            placeholder="Answer — the insight that fixed it."
            rows="2"
            bind:value={newA}
          ></textarea>
          <button class="btn btn-xs btn-primary" onclick={addCard} disabled={!newQ.trim()}>Add card</button>
        </div>
      {/if}

      {#if userCards.length === 0}
        <p class="text-sm opacity-50">None yet — e.g. "Why did recv() return -1 here?" → "EINTR — retry the syscall."</p>
      {:else}
        <div class="space-y-2 max-h-72 overflow-y-auto pr-1">
          {#each userCards as uc}
            {#if editingId === uc.id}
              <div class="flex flex-col gap-2 bg-base-100 rounded-xl px-3 py-2">
                <input class="input input-xs input-bordered w-full" bind:value={editQ} />
                <textarea class="textarea textarea-bordered w-full text-sm" rows="2" bind:value={editA}></textarea>
                <div class="flex gap-1 justify-end">
                  <button class="btn btn-xs btn-primary" onclick={saveEdit} disabled={!editQ.trim()}>Save</button>
                  <button class="btn btn-xs btn-ghost" onclick={() => (editingId = null)}>Cancel</button>
                </div>
              </div>
            {:else}
              <div class="flex items-center gap-2 text-sm bg-base-100 rounded-xl px-3 py-2">
                <button class="flex-1 min-w-0 text-left" onclick={startEdit(uc)} title="Edit">
                  <div class="truncate font-medium">{uc.q}</div>
                  <div class="truncate text-xs opacity-50">{uc.a}</div>
                </button>
                <span class="text-xs badge badge-ghost badge-sm shrink-0">{uc.reps} reps{uc.due ? ` · ${uc.due}` : ""}</span>
                <button class="btn btn-xs btn-ghost text-error shrink-0" onclick={() => deleteCard(uc)} title="Delete">
                  ✕
                </button>
              </div>
            {/if}
          {/each}
        </div>
      {/if}
    </div>
  {:else}
    <!-- active session -->
    <div class="mt-8">
      <div class="flex items-center justify-between text-xs opacity-50 mb-2">
        <span>Session · {answered.length} answered</span>
        <span>{queue.length + (current ? 1 : 0)} remaining</span>
      </div>

      <div class="h-1 bg-base-300 rounded-full mb-6 overflow-hidden">
        <div
          class="h-full bg-primary transition-all"
          style="width:{(answered.length / Math.max(1, answered.length + queue.length + (current ? 1 : 0))) * 100}%"
        ></div>
      </div>

      {#if current}
        <div class="card bg-base-200 rounded-3xl shadow-lg p-8 min-h-[24rem] flex flex-col">
          <div class="flex items-center gap-2 text-xs opacity-50 mb-4">
            <span class="badge badge-ghost badge-sm">{current.module}</span>
            <span class="font-mono">{current.exID}</span>
            <span class="ml-auto">{fmtInterval(current)}</span>
          </div>

          <div class="flex-1">
            <p class="text-lg font-medium leading-relaxed">{current.q}</p>

            {#if revealed}
              <div class="divider my-4"></div>
              <div class="text-sm opacity-90 bg-base-100 rounded-xl p-4 whitespace-pre-wrap">{current.a}</div>
            {:else}
              <div class="mt-10">
                <button class="btn btn-block btn-lg btn-outline" onclick={() => (revealed = true)}>
                  Show answer
                  <span class="text-xs opacity-50">or press space</span>
                </button>
              </div>
            {/if}
          </div>

          {#if revealed}
            <div class="grid grid-cols-4 gap-2 mt-6">
              {#each ratings as r, i}
                <button class="btn {r.cls}" onclick={() => answer(r.v)}>
                  <span class="hidden sm:inline">{i + 1}· </span>{r.label}
                </button>
              {/each}
            </div>
            <div class="text-xs text-center opacity-50 mt-3 h-4">{feedback}</div>
          {:else}
            <div class="h-4"></div>
          {/if}
        </div>
      {:else}
        <div class="card bg-base-200 rounded-2xl p-10 text-center">
          <div class="text-4xl mb-3">🎉</div>
          <div class="font-semibold">Session complete</div>
          <p class="text-sm opacity-60 mt-1">{answered.length} cards reviewed.</p>
          <div class="flex justify-center gap-2 mt-4">
            <button class="btn btn-primary btn-sm" onclick={() => { answered = []; startSession(); }}>Another round</button>
            <button class="btn btn-ghost btn-sm" onclick={endSession}>Finish</button>
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- just-reviewed strip -->
  {#if answered.length > 0 && !sessionActive}
    <div class="mt-8">
      <h3 class="text-sm font-semibold opacity-60 mb-2">Recently reviewed</h3>
      <div class="space-y-1.5">
        {#each answered.slice(0, 8) as c}
          <div class="flex items-center gap-3 text-sm bg-base-200 rounded-lg px-3 py-2">
            <span class="font-mono text-xs opacity-50 shrink-0">{c.exID}</span>
            <span class="truncate flex-1">{c.q}</span>
            <span class="text-xs badge badge-ghost badge-sm">reps {c.reps}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>