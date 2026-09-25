import { mount } from "svelte";
import "./app.css";
import App from "./App.svelte";

// In a plain browser tab (vite dev) there is no Wails webview, so install the
// dev-only binding mock before the app mounts. The `import.meta.env.DEV` is
// statically replaced at build time: production bundles drop this entirely.
if (import.meta.env.DEV) {
  const { installMockBackend } = await import("./lib/mock");
  installMockBackend();
}

// Forward webview errors to the backend's stderr so runtime issues are
// visible outside the window (debugging aid).
window.addEventListener("error", (e) => {
  console.error("[app] uncaught error:", e.error ?? e.message);
  const msg = String(e.error ?? e.message);
  (window as any)?.go?.main?.App?.Log?.(`[app] ${msg}`).catch?.(() => {});
});
window.addEventListener("unhandledrejection", (e) => {
  console.error("[app] unhandled rejection:", e.reason);
  (window as any)?.go?.main?.App?.Log?.(
    `[app] unhandled rejection: ${e.reason instanceof Error ? e.reason.message : String(e.reason)}`,
  ).catch?.(() => {});
});

const app = mount(App, {
  target: document.getElementById("app")!,
});

export default app;