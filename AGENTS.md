# Agent guidance

Read these before you write code.

| Rule | When |
|------|------|
| [`.cursor/rules/documentation-ste100.mdc`](.cursor/rules/documentation-ste100.mdc) | Always. STE-100 documentation standard. |
| [`.cursor/rules/gallery-wall-project.mdc`](.cursor/rules/gallery-wall-project.mdc) | Always. Repo overview, stack, locked decisions. |
| [`.cursor/rules/layout-package.mdc`](.cursor/rules/layout-package.mdc) | Go `layout/**`. API, file map, tests. |
| [`.cursor/rules/web-app.mdc`](.cursor/rules/web-app.mdc) | `web/**`, `wasm/**`, `scripts/**`. WASM and frontend. |

Human-readable docs:

| Doc | Purpose |
|-----|---------|
| [`docs/DESIGN.md`](docs/DESIGN.md) | Algorithm spec. Wall bounds and web app are implemented. |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Repo architecture summary. |
| [`layout/ALGORITHM.md`](layout/ALGORITHM.md) | Placement algorithm walkthrough. Aligned with code. |
| [`layout/README.md`](layout/README.md) | Standalone Go module install and usage (v1.0.0). |
| [`web/README.md`](web/README.md) | Web app dev and build. |

## Workflow

Match conventions in the relevant rule file. Keep diffs minimal. Run `go test ./layout/...` for Go changes. Run `npm run build` in `web/` for UI changes. Do not commit unless asked.
