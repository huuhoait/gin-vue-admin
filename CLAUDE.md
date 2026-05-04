# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository layout

Two independently deployable halves, built together but run separately:

- `server/` — Go 1.25, Gin, GORM. Go module is `github.com/huuhoait/gin-vue-admin/server` (note: **not** the upstream `huuhoaitvn` path). Entrypoint `server/main.go` → `core.RunServer()`. Listens on the port from `server/config.yaml` (`system.addr`, default `8888`).
- `web/` — Vue 3 + Vite 6 + Pinia + Element Plus + UnoCSS. Package manager is **pnpm** (lockfile `pnpm-lock.yaml`); the Makefile still invokes `yarn` inside Docker builds.
- `server/plugin/`, `web/src/plugin/` — paired GVA plugins (see plugin model below).
- `deploy/`, `docs/`, `scripts/` — deployment templates, runbooks, and a CJK-extraction toolchain used by the i18n migration.

`external-frontend-integration.md` at the repo root is the authoritative contract between this admin BFF and an external Vue frontend that consumes SkyAgent Core/Order via this server's proxy layer.

## Common commands

Frontend (run from `web/`):

```bash
pnpm install
pnpm dev           # or: pnpm serve  — Vite dev server on VITE_CLI_PORT (8081)
pnpm build         # production build → web/dist
pnpm preview
```

Backend (run from `server/`):

```bash
go mod tidy
go run .           # starts admin server on :8888 (see config.yaml system.addr)
go build -o server # produces ./server binary
go test ./...
go test ./middleware/... -run TestCasbin   # single package / single test pattern
swag init          # regenerate server/docs/swagger.{json,yaml,go}
```

Combined (from repo root; Docker-backed):

```bash
make build         # builds web + server in containers, produces ./build/
make build-local   # builds both on host (skips Docker)
make doc           # swag init
make plugin PLUGIN=email   # packages server/plugin/<name> + web/src/plugin/<name> into <name>.zip
```

**No root-level lint/test target exists.** Lint via `cd web && pnpm lint` (flat config in `eslint.config.mjs`) and Go vet/test per package. A custom i18n lint (`make admin-i18n-lint`) is referenced in `I18N.md` but **not wired into this Makefile** — use `node scripts/extract-cjk.mjs --scope=all` directly if needed.

Vite dev proxies `VITE_BASE_API` (default `/api`) to `VITE_BASE_PATH:VITE_SERVER_PORT` (default `http://127.0.0.1:8888`) — see `web/vite.config.js` and `web/.env.development`. The proxy **strips** the `/api` prefix before forwarding.

## Architecture: GVA layered backend

This project follows the gin-vue-admin (GVA) layered architecture strictly. See `.claude/rules/project_rules.md` (also symlinked via `.aone_copilot/rules/project_rules.md`) for the full contract — it's the most important file to read before making changes. Key invariants:

- **Layer direction is one-way**: `router → api → service → model`. The API layer must not touch the DB; the service layer must not touch `gin.Context`.
- **`enter.go` group pattern**: each of `api/v1/`, `service/`, `router/` exposes an `ApiGroup` / `ServiceGroup` / `RouterGroup` singleton. Cross-module references go through these globals (`service.ServiceGroupApp.XxxService`, `api.ApiGroupApp.XxxApi`). This is how circular imports are avoided — do not short-circuit it.
- **Swagger annotations are mandatory** on every exposed handler. The FE contract depends on `swag init` producing a complete `docs/swagger.json`.
- **Response envelope** is `{ code, data, msg }` via `server/model/common/response`. All handlers must go through `response.OkWithDetailed` / `response.FailWithMessage` etc. — do not write raw JSON.
- **i18n (forced)**: every user-visible string must go through bundles — no hardcoded UI copy or API `msg` literals in new or touched code. Backend: `server/global/i18n` + TOML in `server/resource/i18n/` (`response.FailWithCode` / message codes, not ad-hoc strings). Frontend: vue-i18n (`web/src/i18n/`, `web/src/i18n/locales/*.json`). Backend logs stay English. Key format: `admin.<domain>.<section>.<label>`. Default locale `vi-VN`, fallback `en-US`. Before commit, run the CJK gate from `I18N.md` (`node scripts/extract-cjk.mjs --scope=all` or `make admin-i18n-lint` if your Makefile defines it).

### Plugin model

Plugins are self-contained vertical slices living in `server/plugin/<name>/` (mirroring GVA directories: `api/ config/ initialize/ model/ router/ service/` plus `plugin.go`) and optionally `web/src/plugin/<name>/`. They register via `init()` in the plugin's `plugin.go` and are activated by anonymous import in `server/plugin/register.go`. Reference implementation: `server/plugin/announcement/`. The `make plugin PLUGIN=<name>` target zips both halves together for distribution.

### SkyAgent proxy layer (project-specific)

This admin is deployed as the **BFF for SkyAgent** — a separate microservice system (Core on `:8001`, Order on `:8002`). The BFF pattern lives in:

- `server/service/proxy/` — HTTP client (`client.go`), per-service wrappers (`core_proxy.go`, `order_proxy.go`), GVA-envelope parser (`response.go`), maker/checker injection from JWT (`auth.go`).
- `server/api/v1/proxy/skyagent_api.go` — thin handlers that forward to Core/Order and surface the envelope unchanged.
- `server/router/proxy/skyagent_router.go` — routes mounted under `/admin-api/v1/` (JWT + Casbin protected).
- `web/src/api/skyagent/` + `web/src/view/{agent,onboarding,catalog,order}/` — the frontend that consumes these proxied endpoints.

**The BFF does not retry mutations.** Timeout is 10s. FE is responsible for `X-Idempotency-Key`. The authoritative contract for FE↔BFF↔Core/Order lives in `external-frontend-integration.md` at the repo root (envelope shapes, error codes, DTOs, pagination modes, idempotency rules).

## Frontend notes

- `src/utils/request.js` is the **only** sanctioned axios wrapper — it handles JWT, token refresh, global loading, and envelope unwrapping. Never `import axios` directly in feature code.
- Common formatting helpers (date, dictionary, PII masks, UUIDs, permission checks, event bus) are pre-built under `src/utils/`; `.claude/rules/project_rules.md` §"Mandatory Usage Requirements" lists the mapping of scenario → utility. Grep `src/utils/` before writing a new one.
- Pinia stores use the Composition API form (`defineStore('name', () => { ... })`) with `useStorage` from `@vueuse/core` for persisted state.
- Element Plus drawers/dialogs **must** set `destroy-on-close` to avoid state leaks.

## Go module path gotcha

The `go.mod` declares `github.com/huuhoait/gin-vue-admin/server` but upstream GVA uses `github.com/huuhoaitvn/gin-vue-admin/server`. Imports inside this repo **must** use `huuhoait`; copy-pasted snippets from upstream docs will compile-fail until the prefix is swapped.

## Swagger / API doc regeneration

Regenerate whenever handlers or DTOs change:

```bash
cd server && swag init
```

Outputs land in `server/docs/`. Swagger UI is served at `/swagger/index.html` when `gin.Mode() != release` or `system.enable-swagger: true` in `config.yaml`.
