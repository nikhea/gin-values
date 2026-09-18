# `config/redis.go`

Redis connection for caching. Deliberately non-fatal: the API serves from Postgres when Redis is down.

## Exports

- `RDB *redis.Client` — shared go-redis client; **nil when Redis is unreachable** (all cache helpers handle nil).
- `ConnectRedis()` — dials `REDIS_ADDR` (default `localhost:6379`), `REDIS_PASSWORD`, `REDIS_DB` (default 0) with 5s dial / 3s IO timeouts; pings, warns + leaves `RDB` nil on failure, logs connect on success. Called in `main()` after the database.
- `CloseRedis() error` — nil-safe release, called on shutdown.

## Notes

- Compose provides a `redis` service (`REDIS_ADDR=redis:6379`); local dev uses `localhost:6379` or nothing at all.
