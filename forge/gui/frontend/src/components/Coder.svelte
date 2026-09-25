<script lang="ts">
  import { onMount } from "svelte";
  import { EditorView, basicSetup } from "codemirror";
  import { EditorState } from "@codemirror/state";
  import { keymap } from "@codemirror/view";
  import { cpp } from "@codemirror/lang-cpp";
  import { go } from "@codemirror/lang-go";
  import { python } from "@codemirror/lang-python";
  import { javascript } from "@codemirror/lang-javascript";
  import { StreamLanguage } from "@codemirror/language";
  import { shell } from "@codemirror/legacy-modes/mode/shell";
  import { oneDark } from "@codemirror/theme-one-dark";

  let {
    value = "",
    lang = "c",
    onChange,
    onSave,
  }: {
    value?: string;
    lang?: string;
    onChange?: (v: string) => void;
    onSave?: () => void;
  } = $props();

  let host: HTMLDivElement;
  let view: EditorView;

  function language() {
    switch (lang.toLowerCase()) {
      case "cpp":
      case "cc":
      case "cxx":
      case "hpp":
      case "c":
      case "h":
        return cpp();
      case "py":
      case "python":
        return python();
      case "js":
      case "jsx":
      case "ts":
      case "tsx":
      case "javascript":
      case "typescript":
        return javascript();
      case "go":
        return go();
      case "sh":
      case "bash":
      case "shell":
      case "mk":
        return StreamLanguage.define(shell);
      default:
        return cpp();
    }
  }

  onMount(() => {
    view = new EditorView({
      state: EditorState.create({
        doc: value,
        extensions: [
          basicSetup,
          language(),
          oneDark,
          EditorView.lineWrapping,
          EditorView.updateListener.of((u) => {
            if (u.docChanged) onChange?.(u.state.doc.toString());
          }),
          keymap.of([
            {
              key: "Mod-s",
              run: () => {
                onSave?.();
                return true;
              },
              preventDefault: true,
            },
          ]),
          EditorView.theme({
            "&": { height: "100%", fontSize: "13px" },
            ".cm-scroller": {
              fontFamily: 'ui-monospace, "JetBrains Mono", "Fira Code", monospace',
              overflow: "auto",
            },
            ".cm-content": { padding: "12px 0" },
          }),
        ],
      }),
      parent: host,
    });
    return () => view.destroy();
  });
</script>

<div class="w-full h-full overflow-hidden" bind:this={host}></div>