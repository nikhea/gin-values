# `config/database.go`

Opens the primary Postgres connection with GORM and exposes it as a package global.

## Exports

- `DB *gorm.DB` — shared handle used by every `repository` function. Must be set by `ConnectDatabase()` before serving (tests overwrite it with a test-database handle).
- `ConnectDatabase()` — reads `DATABASE_URL`, exits if unset/connection fails; pool tuning: 10 idle, 100 open, 1h lifetime; logs `slog.Info("Database connected")`.

## Depends on

`gorm.io/driver/postgres` (pgx-based), `gorm.io/gorm`. Fatal paths log with `slog.Error` and `os.Exit(1)`.

## Notes

- River (`jobs`) uses its own separate pgx pool, not this handle.
