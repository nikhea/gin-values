# `cmd/migrate/main.go`

Standalone CLI for running database migrations manually, without booting the API server.

## Usage

```bash
go run ./cmd/migrate up       # apply pending migrations
go run ./cmd/migrate down     # roll back one step
go run ./cmd/migrate version  # print current version + dirty flag
```

## How it works

1. `config.InitLogger()` + `config.LoadEnv()` for structured logs and `DATABASE_URL`.
2. Opens a `golang-migrate` instance with source `file://migrations` (relative to the repo root — run from there).
3. Dispatches on `os.Args[1]`; `ErrNoChange`/`ErrNilVersion` are reported as friendly messages instead of errors.
4. Fatal problems go through the local `fatal(msg, err)` helper: `slog.Error` + `os.Exit(1)` (same behavior as the old `log.Fatal`).

## Depends on

`config`, `golang-migrate` (postgres driver + file source). Human-readable output uses `fmt.Println`; errors use `slog`.
