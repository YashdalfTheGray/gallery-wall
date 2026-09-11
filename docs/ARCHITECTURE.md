# Architecture

This document describes the structure of the gallery-wall repository.

## Components

```
gallery-wall/                    (root module: web, wasm, cmd)
├── layout/                      (nested Go module: publishable library)
│   ├── go.mod                   github.com/yashdalfthegray/gallery-wall/layout
│   └── …                        placement algorithm and tests
├── wasm/                        GOOS=js bridge, goLayout()
├── web/                         Vite, TypeScript, Web Awesome SPA
├── cmd/gallery-svg/             result JSON to SVG (dev diagram tool)
├── scripts/build-wasm.sh
├── go.mod + go.work             monorepo wiring (replace → ./layout)
└── docs/
```

## Data flow (web app)

1. The user edits params in the frame editor (Form or JSON tab).
2. `runLayout()` sends JSON to WASM `goLayout()`.
3. Go `layout.Layout()` returns positions or typed errors.
4. The preview renders SVG. Stats and session data go to `localStorage`.
5. Export writes session JSON or a standalone SVG.

The deployed static site makes no network calls at runtime.

## layout library

| Item | Value |
|------|-------|
| Import path | `github.com/yashdalfthegray/gallery-wall/layout` |
| Install | `go get github.com/yashdalfthegray/gallery-wall/layout@v1.0.0` |
| Release tags | `layout/v1.0.0` on this repository |
| API | `Layout(params) (Result, error)`, `Validate(params) error` |
| Wall constraint | Optional `wallWidth` and `wallHeight`, centered on anchor `(0, 0)` |

See [layout/README.md](../layout/README.md) and [layout/ALGORITHM.md](../layout/ALGORITHM.md).

## Not implemented

- HTTP API server
- Frame rotation
- Unit conversion (the library uses abstract units)
- Drag-and-drop editing on the preview canvas

## Related docs

- [DESIGN.md](./DESIGN.md) — original v1 design spec
- [.cursor/rules/](../.cursor/rules/) — Cursor agent steering
