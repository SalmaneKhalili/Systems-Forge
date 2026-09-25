<script lang="ts">
  import type { DayActivity } from "../lib/types";

  let { days = [] }: { days?: DayActivity[] } = $props();

  // GitHub-style heatmap: 90 days arranged into a row of 7-day columns.
  const cols = $derived(Math.ceil(days.length / 7));
  const byDay = $derived(
    new Map(days.map((d) => [d.day, d.count])),
  );

  function startOfWeek(): Date {
    const d = new Date();
    const day = d.getDay(); // 0 = Sunday
    d.setHours(0, 0, 0, 0);
    return d; // first column starts on the week of today
  }

  const cells = $derived.by(() => {
    const list: { count: number; title: string; future: boolean }[] = [];
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    // Render `cols` columns × 7 rows, aligned so the rightmost column ends today.
    const grid = new Array<{ count: number; title: string; future: boolean } | null>(cols * 7).fill(null);
    for (let c = 0; c < cols; c++) {
      for (let row = 0; row < 7; row++) {
        const idx = c * 7 + row;
        const day = new Date(today);
        day.setDate(day.getDate() - (cols - 1 - c) * 7 + (row - (today.getDay())));
        const key = `${day.getFullYear()}-${String(day.getMonth() + 1).padStart(2, "0")}-${String(day.getDate()).padStart(2, "0")}`;
        const future = day.getTime() > today.getTime();
        grid[idx] = { count: byDay.get(key) ?? 0, title: `${key}: ${byDay.get(key) ?? 0} runs`, future };
      }
    }
    for (const cell of grid) if (cell) list.push(cell);
    return list;
  });

  function level(count: number): string {
    if (count === 0) return "bg-base-300";
    if (count === 1) return "bg-primary/40";
    if (count <= 3) return "bg-primary/70";
    if (count <= 6) return "bg-primary";
    return "bg-primary-content";
  }
</script>

<div class="flex gap-1 overflow-x-auto py-1">
  {#each Array(cols) as _, c}
    <div class="flex flex-col gap-1">
      {#each Array(7) as _, row}
        {@const cell = cells[c * 7 + row]}
        {#if cell}
          <div
            class="w-3 h-3 rounded-sm {level(cell.count)} {cell.future ? 'opacity-25' : ''}"
            title={cell.title}
          ></div>
        {:else}
          <div class="w-3 h-3 rounded-sm bg-base-300/30"></div>
        {/if}
      {/each}
    </div>
  {/each}
</div>