package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

// JWTSecret returns the HMAC secret used to sign access tokens.
// It fail-closes: missing, known-dev, or short (<32 byte) secrets exit
// the process unless ALLOW_INSECURE_JWT=true explicitly opts into them
// (local development only).
func JWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	insecureAllowed := os.Getenv("ALLOW_INSECURE_JWT") == "true"
	if secret == "" || secret == "dev-only-insecure-secret-change-me" {
		if insecureAllowed {
			slog.Warn("JWT_SECRET not securely set, using insecure dev fallback", "hint", "Set a 32+ byte JWT_SECRET in production")
			return []byte("dev-only-insecure-secret-change-me")
		}
		slog.Error("JWT_SECRET must be set to 32+ random bytes (or set ALLOW_INSECURE_JWT=true for local dev only)")
		os.Exit(1)
		return nil
	}
	if len(secret) < 32 {
		if insecureAllowed {
			slog.Warn("JWT_SECRET shorter than 32 bytes", "hint", "Use 32+ random bytes in production")
			return []byte(secret)
		}
		slog.Error("JWT_SECRET must be 32+ bytes", "hint", "Generate with: openssl rand -hex 32")
		os.Exit(1)
		return nil
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
