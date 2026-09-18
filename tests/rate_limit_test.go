package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin-learn/middleware"

	"github.com/gin-gonic/gin"
)

func TestRateLimitAllowsThenRejects(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", middleware.RateLimit("unittest-basic", 2), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200", i+1, w.Code)
		}
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("request 3 status = %d, want 429", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on 429")
	}
}

func TestRateLimitEnvOverride(t *testing.T) {
	t.Setenv("RATE_LIMIT_UNITTEST_ENV_PER_MINUTE", "100000")
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", middleware.RateLimit("unittest-env", 1), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200 (env override)", i+1, w.Code)
		}
	}
}

func TestRateLimitWindowRefills(t *testing.T) {
	// Separate scope so other tests' buckets can't interfere.
	scope := fmt.Sprintf("unittest-refill-%d", time.Now().UnixNano())
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", middleware.RateLimit(scope, 1), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("first status = %d, want 200", w.Code)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want 429", w.Code)
	}
}
