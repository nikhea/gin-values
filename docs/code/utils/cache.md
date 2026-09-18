# `utils/cache.go` + `utils/cache_keys.go`

Generic Redis helpers and the key-naming scheme.

## Key scheme (`cache_keys.go`)

Every key is namespaced `grip:<domain>:<parts...>` so one Redis can be shared safely:

- `Key(segments...)` — the joiner; **always build keys through it**, never by hand.
- Static domains: `KeyUsers`, `KeyProfiles`, `KeyContacts`, `KeyAuth`.
- Reusable dynamic builders (normalize inputs like email case/space):
  - `UserKey(id)`, `UserEmailKey(email)`, `ProfileKey(userID)`, `ContactKey(id)`
  - `ContactListKey(userID, type, search, page, pageSize)` — filters embedded so queries never collide
  - `OTPAttemptsKey(email)` — counter key for OTP guesses
- `MatchUsers/MatchProfiles/MatchContacts/MatchAuth()` — `prefix:*` SCAN patterns for bulk invalidation.

## Helpers (`cache.go`, all nil-client safe)

- `CacheEnabled()` — false when `config.RDB` is nil.
- `SetJSON(ctx, key, value, ttl)` / `GetJSON(ctx, key, dst)` — JSON marshal round-trip (`DefaultCacheTTL = 5m` when ttl ≤ 0); miss/disabled → `(false, nil)`; real errors returned for DB fallback.
- `Del(ctx, keys...)` / `DelPattern(ctx, pattern)` (SCAN+DEL loop, for invalidations not hot paths) / `Expire(ctx, key, ttl)` — nil-safe no-ops when disabled.
- `Incr(ctx, key, ttl)` — atomic counter, TTL set on first creation; errors when disabled.

Rule: caching must never break a request — callers log-and-continue on errors (see `service/user-service.md`). Tested in `tests/cache_test.go` (+ `requireTestRedis` helper) and `tests/cache_keys_test.go`.
