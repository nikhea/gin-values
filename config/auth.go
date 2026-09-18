package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

// JWTSecret returns the HMAC secret used to sign access tokens.
// JWT_SECRET must be set in production; a dev-only fallback keeps
// local development working without silently securing prod traffic.
func JWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		slog.Warn("JWT_SECRET is not set, using insecure dev fallback", "hint", "Set JWT_SECRET in .env")
		secret = "dev-only-insecure-secret-change-me"
	}
	return []byte(secret)
}

// JWTTTL returns how long issued tokens stay valid (default 24h).
func JWTTTL() time.Duration {
	raw := os.Getenv("JWT_TTL_HOURS")
	if raw == "" {
		return 24 * time.Hour
	}
	hours, err := strconv.Atoi(raw)
	if err != nil || hours <= 0 {
		slog.Warn("Invalid JWT_TTL_HOURS, falling back to 24h", "value", raw)
		return 24 * time.Hour
	}
	return time.Duration(hours) * time.Hour
}

// AppURL is the public base URL used to build email links.
func AppURL() string {
	if url := os.Getenv("APP_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}
