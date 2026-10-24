# gallery-wall

Place picture frames of mixed sizes and shapes around a centerpiece.

## Web app

The browser editor runs the layout engine locally via WASM. See [web/README.md](./web/README.md).

Live site: https://yashdalfthegray.github.io/gallery-wall/

## Go library

The [layout](./layout/) package is a standalone Go module (v1.0.0).

```bash
go get github.com/yashdalfthegray/gallery-wall/layout@v1.0.0
```

```go
import "github.com/yashdalfthegray/gallery-wall/layout"
```

See [layout/README.md](./layout/README.md) for usage and API details.

## CI and docs

PR workflow runs Go tests, a web format check, and a production build. See `.github/workflows/pr.yml`.

The main branch runs the same checks and deploys `web/dist` to GitHub Pages. See `.github/workflows/pages.yml`.

Dependabot opens weekly updates for npm, Go modules, and Actions. See `.github/dependabot.yml`.

[docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) describes repo structure.

[AGENTS.md](./AGENTS.md) describes Cursor agent guidance.
