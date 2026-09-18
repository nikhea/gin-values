package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Rate-limit budgets. Auth endpoints get a strict budget (brute-force
// protection); the general API gets a generous one.
const (
	// AuthRateLimit allows 10 requests per minute per client IP on
	// login/register/verify/OTP/password-reset endpoints.
	AuthRateLimit = 10
	// APIRateLimit allows 300 requests per minute per client IP on
	// authenticated API routes.
	APIRateLimit = 300
	// RateWindow is the fixed window all budgets apply to.
	RateWindow = time.Minute
)

// clientBucket is a fixed-window counter. Buckets are keyed by
// client IP (+ handler scope) and expire after RateWindow.
type clientBucket struct {
	count     int
	expiresAt time.Time
}

var (
	buckets   = sync.Map{} // map[string]*clientBucket guarded by bucketMu
	bucketMu  sync.Mutex
	janitorOn sync.Once
)

// RateLimit returns a fixed-window rate limiter allowing n requests per
// RateWindow per client IP. Scope separates budgets (e.g. "auth" vs "api").
// The budget can be overridden with RATE_LIMIT_<SCOPE>_PER_MINUTE
// (e.g. RATE_LIMIT_AUTH_PER_MINUTE); tests use this to disable limiting.
// Exceeded requests are rejected with 429 + Retry-After.
func RateLimit(scope string, n int) gin.HandlerFunc {
	envKey := "RATE_LIMIT_" + strings.ToUpper(strings.ReplaceAll(scope, "-", "_")) + "_PER_MINUTE"
	if v, err := strconv.Atoi(os.Getenv(envKey)); err == nil && v > 0 {
		n = v
	}
	janitorOn.Do(startJanitor)
	return func(c *gin.Context) {
		key := scope + "|" + c.ClientIP()
		if !take(key, n) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded, try again later",
			})
			return
		}
		c.Next()
	}
}

// take consumes one token from the bucket, refilling per window.
func take(key string, n int) bool {
	now := time.Now()
	bucketMu.Lock()
	defer bucketMu.Unlock()

	v, ok := buckets.Load(key)
	if !ok || v.(*clientBucket).expiresAt.Before(now) {
		buckets.Store(key, &clientBucket{count: 1, expiresAt: now.Add(RateWindow)})
		return true
	}
	b := v.(*clientBucket)
	if b.count >= n {
		return false
	}
	b.count++
	return true
}

// startJanitor periodically drops expired buckets so the map can't grow
// unboundedly from rotating client IPs.
func startJanitor() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			bucketMu.Lock()
			buckets.Range(func(k, v any) bool {
				if v.(*clientBucket).expiresAt.Before(now) {
					buckets.Delete(k)
				}
				return true
			})
			bucketMu.Unlock()
		}
	}()
}
