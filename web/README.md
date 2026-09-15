# Gallery Wall — Web

This static app runs the layout engine in the browser via Go WASM. No server is required after build.

## Dev

```bash
cd web
npm install
npm run dev
```

Vite opens port 5173. The `predev` script rebuilds WASM.

## Production build

```bash
cd web
npm run build
```

Output is in `web/dist/`.

For GitHub Pages (project site under `/gallery-wall/`):

```bash
VITE_BASE=/gallery-wall/ npm run build
```

Local `vite` and default `npm run build` use `base: "/"`. Assets resolve at the site root.

## Features

- Layout engine: Go `layout` package compiled to WASM (`goLayout` global)
- Frame editor: Form tab (centerpiece and satellite batches with copy counts) or JSON tab
- Wall bounds: optional width and height constrain the cluster (0 = unbounded)
- Presets: 25-frame sample (`public/presets/twentyfive.json`), random 5–20 frames
- localStorage: autosaves params and last result between visits
- Import and export: session JSON (`gallery-wall-state.json`), wall SVG download
- Preview: SVG with shape colors, measure labels, stats panel (frames, gap, cluster size, anchor, bounds)
- Theme: Web Awesome light or dark follows system preference

## WASM rebuild only

```bash
./scripts/build-wasm.sh
```

This writes `web/public/layout.wasm` and `wasm_exec.js`. These files are gitignored. Rebuild before dev or build.

## Source layout

See [`.cursor/rules/web-app.mdc`](../.cursor/rules/web-app.mdc) for the module map and conventions.
