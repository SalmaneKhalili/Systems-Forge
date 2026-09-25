import MarkdownIt from "markdown-it";
import hljs from "highlight.js";
import "highlight.js/styles/github-dark.css";

const md = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: false,
  highlight(str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return (
          `<pre class="hljs"><code>` +
          hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
          `</code></pre>`
        );
      } catch {
        /* fall through to escaped output */
      }
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`;
  },
});

export function renderMarkdown(src: string): string {
  return md.render(src);
}