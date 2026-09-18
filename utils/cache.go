package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"grip/config"

	"github.com/redis/go-redis/v9"
)

// DefaultCacheTTL applies when callers pass 0.
const DefaultCacheTTL = 5 * time.Minute

// CacheEnabled reports whether a Redis client is available. All helpers
// below degrade gracefully when it isn't (miss / no-op).
func CacheEnabled() bool { return config.RDB != nil }

// client returns the shared client or nil when caching is disabled.
func client() *redis.Client { return config.RDB }

// SetJSON marshals value and stores it under key for ttl (DefaultCacheTTL
// when ttl <= 0). Nil client or Redis errors are returned to the caller;
// caching must never break the request path, so callers log-and-continue.
func SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	c := client()
	if c == nil {
		return errors.New("cache disabled")
	}
	if ttl <= 0 {
		ttl = DefaultCacheTTL
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal %s: %w", key, err)
	}
	return c.Set(ctx, key, raw, ttl).Err()
}

// GetJSON loads key and unmarshals into dst. Returns (false, nil) on miss
// or when caching is disabled; other errors (Redis down, corrupt data)
// are returned so callers can fall back to the database.
func GetJSON(ctx context.Context, key string, dst any) (bool, error) {
	c := client()
	if c == nil {
		return false, nil
	}
	raw, err := c.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cache get %s: %w", key, err)
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return false, fmt.Errorf("cache unmarshal %s: %w", key, err)
	}
	return true, nil
}

// Del removes one or more keys. Nil client is a no-op returning nil.
func Del(ctx context.Context, keys ...string) error {
	c := client()
	if c == nil || len(keys) == 0 {
		return nil
	}
	return c.Del(ctx, keys...).Err()
}

// DelPattern removes all keys matching a SCAN pattern (see Match* in
// cache_keys.go). Use sparingly — SCAN + DEL loops are for invalidation
// events, not the request hot path.
func DelPattern(ctx context.Context, pattern string) error {
	c := client()
	if c == nil {
		return nil
	}
	var cursor uint64
	for {
		keys, next, err := c.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("cache scan %s: %w", pattern, err)
		}
		if len(keys) > 0 {
			if err := c.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("cache del pattern %s: %w", pattern, err)
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

// Expire sets a TTL on an existing key. Nil client is a no-op.
func Expire(ctx context.Context, key string, ttl time.Duration) error {
	c := client()
	if c == nil {
		return nil
	}
	return c.Expire(ctx, key, ttl).Err()
}

// Incr atomically increments a counter, setting its TTL on first creation.
// Handy for rate counters and attempt tallies. Returns the new value.
func Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	c := client()
	if c == nil {
		return 0, errors.New("cache disabled")
	}
	n, err := c.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("cache incr %s: %w", key, err)
	}
	if n == 1 && ttl > 0 {
		_ = c.Expire(ctx, key, ttl).Err()
	}
	return n, nil
}
