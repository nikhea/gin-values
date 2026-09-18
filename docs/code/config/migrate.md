# `config/migrate.go`

Applies versioned SQL migrations from `migrations/` using `golang-migrate`, plus a version inspector.

## Exports

- `RunMigrations()` — `migrate.New("file://migrations", DATABASE_URL)` then `Up()`. `ErrNoChange` logs "already up to date"; any other error is fatal. Called on every boot in `main()`.
- `MigrationVersion() (uint, bool)` — returns current version + dirty flag for debugging; logs instead of printing.
- `fatal(msg, err)` (private) — `slog.Error` + `os.Exit(1)`, mirroring `log.Fatal`.

## Depends on

`golang-migrate` postgres driver + file source. `DATABASE_URL` must be set.

## Notes

- River's own tables are migrated separately by `rivermigrate` inside `jobs.Setup`, not here.
- "Dirty" means a previous migration failed halfway — fix with `go run ./cmd/migrate force <version>` (see `cmd/migrate-main.md`).
