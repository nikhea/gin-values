# Compose files: base + environment overrides

Three files; always pass two `-f` flags (or rely on the base defaulting to the production target).

## `docker-compose.yml` (base)

Shared `db` (`postgres:16-alpine`, `pgdata` volume, `pg_isready` healthcheck) and the `api` skeleton: `depends_on` healthy db, `dns: [127.0.0.11]` (embedded DNS — required where host forwarders are broken, otherwise `db` won't resolve), env defaults (`DATABASE_URL` → `db:5432`, dev `JWT_SECRET` fallback, `EMAIL_*` passthrough), port `8080`, `uploads` volume, `/health` healthcheck. Build `target: production` unless overridden.

## `docker-compose.dev.yml`

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Sets `target: development` (Air live reload), bind-mounts the source (`.:/app`, `uploads` stays a named volume), isolated `gocache` volume, `GIN_MODE=debug`, longer healthcheck grace for Air's first build. Edit → Air rebuilds → serving again (~30–60s cold). Config: `.air.toml` (watches `.go`/templates, ignores `docs/`, `uploads/`, tests).

## `docker-compose.prod.yml`

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

`target: production`, `restart: unless-stopped`, `GIN_MODE=release`. Set a real `JWT_SECRET` (`openssl rand -hex 32`) via shell env or `.env.production`.

## `.env` interaction

Compose auto-loads the repo `.env` for `${VAR}` interpolation — local `EMAIL_*`/`JWT_SECRET` flow into containers unless overridden in the shell (shell wins). `DATABASE_URL` is set explicitly, so the local one is ignored.

Stop with `docker compose down` (add `-f` flags to match; `-v` also drops `pgdata`/`uploads`). The `migrate` binary inside the image handles manual migrations: `docker compose run --rm api /app/migrate <up|down|version>`.
