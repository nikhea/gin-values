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
- `vendor` — re-runs `go mod vendor` **and** re-sanitizes the known dummy Slack webhook in `vendor/github.com/go-openapi/spec/appveyor.yml`, then fails if any secret pattern remains. Always use this instead of raw `go mod vendor`, otherwise the dummy secret returns and the next push gets blocked by GitHub push protection.
- `check` — the full gate: `fmt-check` + `vet` + `lint` + `test`. Same steps run in CI (`.github/workflows/ci.md`).
