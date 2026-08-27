# gallery-wall

Lay out picture frames of mixed sizes and shapes around a centerpiece.

## Web app

Browser-based editor and preview — runs the layout engine locally via WASM.
See [web/README.md](./web/README.md).

Live site (GitHub Pages): https://yashdalfthegray.github.io/gallery-wall/

## Go library

The [`layout`](./layout/) package is a standalone Go module (v1.0.0). Import it in your
own project:

```bash
go get github.com/yashdalfthegray/gallery-wall/layout@v1.0.0
```

```go
import "github.com/yashdalfthegray/gallery-wall/layout"
```

See [layout/README.md](./layout/README.md) for usage and API details.

## CI / Docs

- PRs: Go tests + web format check + production build (`.github/workflows/pr.yml`)
- `main`: same checks, then deploy `web/dist` to GitHub Pages (`.github/workflows/pages.yml`)
- Dependabot: weekly npm, Go modules, and Actions updates (`.github/dependabot.yml`)
- [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) — repo structure
- [AGENTS.md](./AGENTS.md) — Cursor agent guidance
