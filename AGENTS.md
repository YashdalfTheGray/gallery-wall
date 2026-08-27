# Agent guidance

Read these before writing code:

| Rule | When |
|------|------|
| [`.cursor/rules/gallery-wall-project.mdc`](.cursor/rules/gallery-wall-project.mdc) | Always — repo overview, stack, locked decisions |
| [`.cursor/rules/layout-package.mdc`](.cursor/rules/layout-package.mdc) | Go `layout/**` — API, file map, tests |
| [`.cursor/rules/web-app.mdc`](.cursor/rules/web-app.mdc) | `web/**`, `wasm/**`, `scripts/**` — WASM + frontend |

Human-readable docs:

| Doc | Purpose |
|-----|---------|
| [`docs/DESIGN.md`](docs/DESIGN.md) | Algorithm spec (original design; wall bounds & web app now implemented) |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Current repo architecture summary |
| [`layout/ALGORITHM.md`](layout/ALGORITHM.md) | Placement algorithm walkthrough (aligned with code) |
| [`layout/README.md`](layout/README.md) | Standalone Go module install/usage (v1.0.0) |
| [`web/README.md`](web/README.md) | Web app dev/build |

## Workflow

- Match conventions in the relevant rule file; keep diffs minimal.
- Run `go test ./layout/...` for Go changes; `npm run build` in `web/` for UI changes.
- Do not commit unless asked.
