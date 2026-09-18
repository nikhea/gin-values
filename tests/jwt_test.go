package tests

import (
	"testing"
	"time"

	"gin-learn/utils"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-jwt-unit-tests-01")

	token, err := utils.GenerateToken("user-123", "a@b.c")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.Subject != "user-123" {
		t.Fatalf("subject = %q, want %q", claims.Subject, "user-123")
	}
	if claims.Email != "a@b.c" {
		t.Fatalf("email = %q, want %q", claims.Email, "a@b.c")
	}
	if claims.ExpiresAt == nil || time.Until(claims.ExpiresAt.Time) <= 0 {
		t.Fatal("expected future expiry")
	}
}

func TestParseTokenTampered(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-jwt-unit-tests-01")

	token, err := utils.GenerateToken("user-123", "a@b.c")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if _, err := utils.ParseToken(token + "tampered"); err == nil {
		t.Fatal("expected error for tampered token")
	}
	if _, err := utils.ParseToken("not.a.token"); err == nil {
		t.Fatal("expected error for malformed token")
	}
	if _, err := utils.ParseToken(""); err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "correct-secret-32bytes-000000001")

	token, err := utils.GenerateToken("user-123", "a@b.c")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	t.Setenv("JWT_SECRET", "different-secret-32bytes-00000002")
	if _, err := utils.ParseToken(token); err == nil {
		t.Fatal("expected error when secret differs")
	}
}

func TestParseTokenExpired(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-for-jwt-unit-tests-01")

	claims := utils.Claims{
		Email: "a@b.c",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte("test-secret-for-jwt-unit-tests-01"))
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}
	if _, err := utils.ParseToken(expired); err == nil {
		t.Fatal("expected error for expired token")
	}
}
