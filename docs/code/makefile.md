# `Makefile`

Developer workflow targets (run `make help` for the list).

## Targets

- `fmt` / `fmt-check` — format / verify formatting of hand-written Go files (generated `docs/` excluded via `git ls-files`).
- `vet` — `go vet ./...`.
- `lint` — `golangci-lint run ./...` (auto-installs v2.13.2 if missing). Config: `.golangci.yml`.
- `test` — `go test -count=1 ./...` (needs Postgres; see `tests/setup.md`).
- `build` — `go build ./...`.
- `swagger` — `swag init`.
- `migrate` — `go run ./cmd/migrate $(ARGS)`, e.g. `make migrate ARGS=up`.
- `check` — the full gate: `fmt-check` + `vet` + `lint` + `test`. Same steps run in CI (`.github/workflows/ci.md`).
