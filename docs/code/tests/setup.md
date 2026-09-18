# `tests/setup_test.go` — integration test harness

Shared setup for all DB-backed tests (`package tests`).

## Helpers

- `testDatabaseURL(t)` — `TEST_DATABASE_URL` or default local `gin_app_test` DSN.
- `requireTestDB(t)` — skips the test if Postgres is unreachable (keeps CI green without a DB); otherwise creates the DB if missing, `AutoMigrate`s `User`/`Profile`/`Contact`, points `config.DB` at it, starts the River client (`jobs.Setup`) against the test DB, truncates tables, sets test JWT secret + blank email creds, enables Gin test mode, and registers cleanup (truncate + River shutdown).
- `truncateAll(t)` — `TRUNCATE users, profiles, contacts, river_job ... CASCADE`.
- `tryCreateTestDatabase` / `databaseName` — create the test DB via the `postgres` maintenance DB.

## Notes

- Unit tests (no DB) live in the same `tests/` directory by project convention — see `tests/` docs below.
