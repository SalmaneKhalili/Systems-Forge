<script lang="ts">
  import { renderMarkdown } from "../lib/markdown";

  let { source = "" }: { source?: string } = $props();
  let html = $derived(renderMarkdown(source));
</script>

<div class="md-body text-sm" class:md-body={true}>{@html html}</div>