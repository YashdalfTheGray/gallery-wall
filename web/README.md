# Gallery Wall — Web

Static app: layout runs in the browser via Go WASM. No server required after build.

## Dev

```bash
cd web
npm install
npm run dev
```

Opens Vite on port 5173. WASM is rebuilt automatically via `predev`.

## Production build

```bash
cd web
npm run build
```

Output in `web/dist/`.

For GitHub Pages (project site under `/gallery-wall/`):

```bash
VITE_BASE=/gallery-wall/ npm run build
```

Local `vite` / default `npm run build` use `base: "/"` so assets resolve at the site root.

## Features

- **Layout engine** — Go `layout` package compiled to WASM (`goLayout` global)
- **Frame editor** — Form tab (centerpiece + satellite batches with copy counts) or raw JSON tab
- **Wall bounds** — optional width/height constrain the cluster (0 = unbounded)
- **Presets** — 25-frame sample (`public/presets/twentyfive.json`), random 5–20 frames
- **localStorage** — autosaves params + last result between visits
- **Import / export** — session JSON (`gallery-wall-state.json`), wall SVG download
- **Preview** — SVG with shape-aware colors, measure labels, stats panel (frames, gap, cluster size, anchor, bounds)
- **Theme** — Web Awesome light/dark follows system preference

## WASM rebuild only

```bash
./scripts/build-wasm.sh
```

Writes `web/public/layout.wasm` and `wasm_exec.js` (gitignored; required before dev/build).

## Source layout

See [`.cursor/rules/web-app.mdc`](../.cursor/rules/web-app.mdc) for module map and conventions.
