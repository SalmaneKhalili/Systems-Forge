import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";

// `base: "./"` is required so Wails can serve the built assets from its
// embedded filesystem (paths are relative to the page, not the host root).
export default defineConfig({
  base: "./",
  plugins: [svelte(), tailwindcss()],
  build: {
    outDir: "dist",
    target: "es2022",
  },
  server: {
    port: 5173,
    strictPort: false,
  },
});