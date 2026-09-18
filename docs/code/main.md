# `main.go`

Application entry point. Wires config, database, background jobs and HTTP routes together, then serves the API with graceful shutdown.

## Startup sequence

1. `config.InitLogger()` — installs the JSON `slog` logger (must run first so all later logs are structured).
2. `config.LoadEnv()` — loads `.env` into the environment.
3. `config.RunMigrations()` — applies pending SQL migrations from `migrations/`.
4. `config.ConnectDatabase()` — opens the GORM Postgres connection (`config.DB`).
5. `jobs.Setup(ctx, DATABASE_URL)` — migrates River tables, registers workers, starts the job client. Fatal if it fails.
6. `routes.Setup()` — builds the Gin engine.
7. Serves on `:8080` via `http.Server`; on `SIGINT`/`SIGTERM` shuts down the HTTP server and the River client (15s timeout).

## Swagger annotations

The comment block above `main()` holds the global API docs consumed by `swag init`: title, version, host, `BasePath /api`, and the `BearerAuth` security definition (header `Authorization`, value `Bearer <token>`).

## Depends on

`config`, `jobs`, `routes`, stdlib (`context`, `log/slog`, `net/http`, `os`, `os/signal`, `syscall`, `time`).
