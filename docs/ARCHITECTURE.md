# Architecture

Current structure of the gallery-wall repository.

## Components

```
gallery-wall/                    (root module — web, wasm, cmd)
├── layout/                      (nested Go module — publishable library)
│   ├── go.mod                   github.com/yashdalfthegray/gallery-wall/layout
│   └── …                        placement algorithm + tests
├── wasm/                        GOOS=js bridge → goLayout()
├── web/                         Vite + TypeScript + Web Awesome SPA
├── cmd/gallery-svg/             result JSON → SVG (dev/diagram tool)
├── scripts/build-wasm.sh
├── go.mod + go.work             monorepo wiring (replace → ./layout)
└── docs/
```

## Data flow (web app)

1. User edits params in the frame editor (Form or JSON tab).
2. `runLayout()` sends JSON to WASM `goLayout()`.
3. Go `layout.Layout()` returns positions or typed errors.
4. Preview renders SVG; stats + session saved to `localStorage`.
5. Export: session JSON or standalone SVG.

No network calls at runtime after static deploy.

## layout library

- **Import:** `github.com/yashdalfthegray/gallery-wall/layout`
- **Install:** `go get github.com/yashdalfthegray/gallery-wall/layout@v1.0.0`
- **Release tags:** `layout/v1.0.0` on this repository
- **API:** `Layout(params) (Result, error)`, `Validate(params) error`
- **Optional wall constraint:** `wallWidth` + `wallHeight` center on anchor `(0, 0)`

See [`layout/README.md`](../layout/README.md) and [`layout/ALGORITHM.md`](../layout/ALGORITHM.md).

## Not implemented

- HTTP API server
- Frame rotation
- Unit conversion (library uses abstract units throughout)
- Drag-and-drop editing on the preview canvas

## Related docs

- [`DESIGN.md`](./DESIGN.md) — original v1 design spec
- [`.cursor/rules/`](../.cursor/rules/) — Cursor agent steering (preferred for up-to-date context)
