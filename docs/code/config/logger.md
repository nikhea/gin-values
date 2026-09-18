# `config/logger.go`

Central logging setup. Installs the process-wide structured logger.

## Exports

- `InitLogger()` — builds `slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))` and registers it via `slog.SetDefault`. Emits one JSON object per line on stdout, including `Debug` records.

## Usage

Call once at startup before anything logs — `main()` and `cmd/migrate` both do this first. Afterwards, log with the `log/slog` package directly:

```go
slog.Info("Server listening", "addr", ":8080")
slog.Warn("Mailer disabled", "hint", "Set EMAIL_ADDRESS/EMAIL_PASSWORD")
slog.Error("Database connection failed", "error", err)
```

## Notes

- Replaces all `log.Println`/`log.Fatal` usage; fatal paths now do `slog.Error(...)` + `os.Exit(1)` at the call site.
- Gin's own `[GIN-debug]` request logs are separate (Gin's internal logger) and stay plain-text.
