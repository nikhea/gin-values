package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-learn/jobs"
	"gin-learn/routes"

	"github.com/gin-gonic/gin"
	"github.com/riverqueue/river"
)

func TestHealthOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	routes.Setup().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", w.Code)
	}
}

func TestReadyzOK(t *testing.T) {
	requireTestDB(t)
	w := doRequest(t, routes.Setup(), "GET", "/readyz", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("readyz status = %d, body %s", w.Code, w.Body.String())
	}
}

func TestCORSHeaders(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")
	gin.SetMode(gin.TestMode)
	router := routes.Setup()

	req := httptest.NewRequest(http.MethodOptions, "/api/auth/login", nil)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("allow-origin = %q, want the configured origin", got)
	}

	// Unlisted origin gets no CORS headers.
	req = httptest.NewRequest(http.MethodOptions, "/api/auth/login", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("allow-origin = %q for unlisted origin, want empty", got)
	}
}

func TestCORSDisabledByDefault(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	gin.SetMode(gin.TestMode)
	router := routes.Setup()

	req := httptest.NewRequest(http.MethodOptions, "/api/auth/login", nil)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("allow-origin = %q without configured origins, want empty", got)
	}
}

func TestCleanupWorkerWithPool(t *testing.T) {
	requireTestDB(t)
	worker := &jobs.CleanupWorker{}
	var job *river.Job[jobs.CleanupArgs]
	if err := worker.Work(context.Background(), job); err != nil {
		t.Fatalf("cleanup work: %v", err)
	}
}
