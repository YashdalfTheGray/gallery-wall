import { defineConfig } from "vite";

// Default "/" for local `vite` / `vite preview`.
// GitHub Pages sets VITE_BASE=/gallery-wall/ in CI.
export default defineConfig({
  base: process.env.VITE_BASE ?? "/",
  root: ".",
  publicDir: "public",
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
