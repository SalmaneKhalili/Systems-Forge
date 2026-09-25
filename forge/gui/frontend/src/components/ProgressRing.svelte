<script lang="ts">
  let {
    value = 0,
    total = 1,
    size = 120,
    stroke = 8,
  }: { value?: number; total?: number; size?: number; stroke?: number } = $props();

  const pct = $derived(total > 0 ? Math.min(100, Math.round((value / total) * 100)) : 0);
  const r = $derived((size - stroke) / 2);
  const circ = $derived(2 * Math.PI * r);
  const offset = $derived(circ * (1 - pct / 100));
  const color = $derived(pct === 100 ? "var(--color-success)" : "var(--color-primary)");
</script>

<div class="relative inline-flex items-center justify-center" style="width:{size}px;height:{size}px">
  <svg width={size} height={size} class="-rotate-90">
    <circle
      cx={size / 2}
      cy={size / 2}
      r={r}
      fill="none"
      stroke="currentColor"
      stroke-opacity="0.15"
      stroke-width={stroke}
    />
    <circle
      cx={size / 2}
      cy={size / 2}
      r={r}
      fill="none"
      stroke={color}
      stroke-width={stroke}
      stroke-linecap="round"
      stroke-dasharray={circ}
      stroke-dashoffset={offset}
      style="transition: stroke-dashoffset 600ms ease;"
    />
  </svg>
  <div class="absolute text-center">
    <slot>
      <span class="text-2xl font-bold">{pct}%</span>
    </slot>
  </div>
</div>