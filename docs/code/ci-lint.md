# `.golangci.yml` + `.github/workflows/ci.yml`

Lint configuration and continuous integration (golangci-lint **v2** — v1 can't typecheck modern Go).

## `.golangci.yml`

Enabled linters: `errcheck`, `govet`, `staticcheck`, `unused`, `ineffassign`, `misspell`, `revive`, `unconvert`, `whitespace`; formatters `gofmt` + `goimports` (local prefix `gin-learn`). Generated `docs/` excluded everywhere; `_test.go` files exempt from `errcheck`/`revive`; a few `revive` style rules (`exported`, `package-comments`, `var-naming`) disabled project-wide.

## `.github/workflows/ci.yml`

Runs on every push/PR: Postgres 16 service → `make fmt-check` → `make vet` → pinned `golangci-lint@v2.13.2` install + run → `make test` (against `gin_app_test` via `TEST_DATABASE_URL`).

## `.github/workflows/docker.yml`

Runs on every push/PR: builds both Dockerfile targets (`development`, `production`), then boots the prod compose stack, waits for `/health`, runs a register smoke test, and tears everything down (`down -v`).
