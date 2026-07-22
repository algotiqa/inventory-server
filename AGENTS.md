# Inventory Server — AGENTS.md

## Toolchain

- Go 1.25.8 (see `go.mod`). No Makefile, no CI, no linter/formatter config, no `*_test.go` files. `go build -o inventory-server .` builds; `go test ./...` passes trivially because there is nothing to test — do not treat a green `go test` as real verification.
- Every `.go` file starts with the Elastic License 2.0 (ELv2) header block (see top of `main.go`). Add it to new files.

## Run

```bash
go build -o inventory-server .
./inventory-server
```

Hard startup failures if any of these are missing/unreachable:
- MySQL at `localhost:3400`, db `inventory`.
- RabbitMQ at `localhost:8450`.
- Keycloak OIDC at `https://algotiqa-server:8443/auth/realms/algotiqa`.
- Platform services on host `algotiqa-server`: system-adapter `:8449`, portfolio `:8445`, storage `:8452`.
- TLS client cert files `ca.crt`, `server.crt`, `server.key` are loaded by **relative path from the current working directory** (`main.go` → `req.AddDefaultClient(...)`). They are gitignored and currently sit in `config/` — running from repo root means they must be present at root too (copy or symlink).
- The `algotiqa-server` hostname must resolve (hostfile or DNS); it is not localhost.

## Config & secrets

- `config/inventory-server.yaml` **is committed with non-placeholder credentials**: MySQL/RabbitMQ `admin`/`admin`, Keycloak `clientSecret`, and a live `provider.currency.apiKey: fca_live_...` for freecurrencyapi.com. Do not "fix" these as example values, and do not commit new real secrets.
- `journal.bin` is a transient runtime messaging journal (written to CWD via `messaging.journal.directory: .`) — gitignored, never commit.
- `agent/` (built agent binaries) and `inventory-server.iml` (IntelliJ project file) are untracked local artifacts — do not commit.

## Architecture

Entry point `main.go`. Startup order — do not reorder, DB and messaging must init before service/process:

`boot.ReadConfig → boot.InitLogger → boot.InitEngine → initClients → dbms.InitDatabase → msg.InitMessaging → service.Init → process.Init → system.InitMessageListener → boot.RunHttpServer`

Packages:
- `pkg/service/` — Gin HTTP handlers. All routes registered in `service.go` under `/api/inventory/v1/`, every one secured with OIDC + role `Admin_User_Service`.
- `pkg/business/` — business logic. `pkg/business/importexport/` handles ZIP packaging for trading-system portability.
- `pkg/db/` — GORM models + DB ops. Every model sets its table name via `TableName()` as **lowercase singular snake_case** (e.g. `trading_system`, `currency_history`). Add `TableName()` to any new model.
- `pkg/platform/` — HTTP clients to system-adapter, portfolio, storage.
- `pkg/core/process/` — background goroutines started by `process.Init`: `agentscanner` (polls remote agents) and `currencyupdater` (pulls freecurrencyapi.com).
- `pkg/core/messaging/system/` — RabbitMQ listener for system-adapter messages.
- `pkg/app/config.go` — `Config` struct embedding `core.{Application,Database,Authentication,Platform,Messaging}` plus a local `Provider` (currency).

## Conventions

- All DB work in handlers goes through `dbms.RunInTransaction(func(tx *gorm.DB) error {...})`, with `c.ReturnError(err)` as the final statement. Follow this in any new handler (see `pkg/service/trading-system.go`).
- Use `types.Date` / `types.TradingSession` from `github.com/algotiqa/types` (not `time.Time`/string) for those fields — see `pkg/db/model.go`.
- Trading-system import/export is multipart: a JSON spec part (`ImportOverviewSpec` / `ImportExecutionSpec`) followed by a ZIP package part. Import is two-step: `POST /trading-systems/import/overview` then `POST /trading-systems/import/execute`.
- Rely on the shared `github.com/algotiqa/core` framework (boot, auth, dbms, msg, req, roles) rather than reimplementing.
