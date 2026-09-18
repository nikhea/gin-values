package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"grip/middleware"
	"grip/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", middleware.AuthRequired(), func(c *gin.Context) {
		id, _ := middleware.GetUserID(c)
		email, _ := middleware.GetUserEmail(c)
		c.JSON(http.StatusOK, gin.H{"userID": id, "email": email})
	})
	return r
}

func TestAuthRequiredTable(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-middleware-tests")

	valid, err := utils.GenerateToken("user-abc", "u@example.com")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	expiredClaims := utils.Claims{
		Email: "u@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-abc",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).
		SignedString([]byte("test-secret-for-middleware-tests"))
	if err != nil {
		t.Fatalf("sign expired: %v", err)
	}

	wrongSecret, err := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.Claims{
		Email:            "u@example.com",
		RegisteredClaims: jwt.RegisteredClaims{Subject: "user-abc", ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	}).SignedString([]byte("other-secret"))
	if err != nil {
		t.Fatalf("sign wrong secret: %v", err)
	}

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"empty bearer", "Bearer ", http.StatusUnauthorized},
		{"wrong scheme", "Token " + valid, http.StatusUnauthorized},
		{"garbage token", "Bearer garbage", http.StatusUnauthorized},
		{"wrong secret", "Bearer " + wrongSecret, http.StatusUnauthorized},
		{"expired token", "Bearer " + expired, http.StatusUnauthorized},
		{"valid token", "Bearer " + valid, http.StatusOK},
		{"valid bare token without scheme", valid, http.StatusOK},
	}

	router := setupAuthTestRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestAuthRequiredSetsContext(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-middleware-tests")

	token, err := utils.GenerateToken("user-abc", "u@example.com")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	setupAuthTestRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"userID":"user-abc"`) || !strings.Contains(body, `"email":"u@example.com"`) {
		t.Fatalf("context not propagated, body: %s", body)
	}
}

func TestGetUserIDWithoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if _, ok := middleware.GetUserID(c); ok {
		t.Fatal("expected ok=false without AuthRequired")
	}
	if _, ok := middleware.GetUserEmail(c); ok {
		t.Fatal("expected ok=false without AuthRequired")
	}
}
