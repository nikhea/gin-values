package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"grip/dto"
	"grip/models"
	"grip/repository"
	"grip/routes"
	services "grip/service"
	"grip/utils"
)

func loginPair(t *testing.T, email, password string) (*models.User, string, string) {
	t.Helper()
	user, access, refresh, err := services.Login(dto.LoginRequest{Email: email, Password: password})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatal("expected access + refresh tokens")
	}
	return user, access, refresh
}

func TestRotateHappyPath(t *testing.T) {
	requireTestDB(t)
	registerVerifiedUser(t, "Refresh User", "refresh@example.com", "supersecret123")
	_, _, refresh1 := loginPair(t, "refresh@example.com", "supersecret123")

	access2, refresh2, user, err := services.Rotate(refresh1)
	if err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if access2 == "" || refresh2 == "" || refresh2 == refresh1 {
		t.Fatal("expected fresh access + refresh tokens")
	}
	if user.Email != "refresh@example.com" {
		t.Fatalf("user = %q", user.Email)
	}
	if _, err := utils.ParseToken(access2); err != nil {
		t.Fatalf("new access token invalid: %v", err)
	}

	// Old token is single-use: presenting it again signals reuse and
	// revokes the whole family. The killed replacement is then plain
	// invalid (it was never rotated itself).
	if _, _, _, err := services.Rotate(refresh1); err != services.ErrRefreshReuse {
		t.Fatalf("reuse err = %v, want ErrRefreshReuse", err)
	}
	if _, _, _, err := services.Rotate(refresh2); err != services.ErrInvalidRefreshToken {
		t.Fatalf("family member err = %v, want ErrInvalidRefreshToken", err)
	}
}

func TestRotateRejectsBadTokens(t *testing.T) {
	requireTestDB(t)
	registerVerifiedUser(t, "Refresh Bad", "refreshbad@example.com", "supersecret123")

	if _, _, _, err := services.Rotate(""); err != services.ErrInvalidRefreshToken {
		t.Fatalf("empty err = %v, want ErrInvalidRefreshToken", err)
	}
	if _, _, _, err := services.Rotate("no-such-token"); err != services.ErrInvalidRefreshToken {
		t.Fatalf("unknown err = %v, want ErrInvalidRefreshToken", err)
	}

	// Expired token.
	_, _, refresh := loginPair(t, "refreshbad@example.com", "supersecret123")
	row, err := repository.GetRefreshTokenByHash(utils.SHA256Hex(refresh))
	if err != nil {
		t.Fatalf("load row: %v", err)
	}
	past := time.Now().Add(-time.Hour)
	row.ExpiresAt = past
	if err := repository.UpdateRefreshToken(row); err != nil {
		t.Fatalf("expire row: %v", err)
	}
	if _, _, _, err := services.Rotate(refresh); err != services.ErrInvalidRefreshToken {
		t.Fatalf("expired err = %v, want ErrInvalidRefreshToken", err)
	}
}

func TestLogoutAndLogoutAll(t *testing.T) {
	requireTestDB(t)
	registerVerifiedUser(t, "Logout User", "logout@example.com", "supersecret123")
	_, _, refreshA := loginPair(t, "logout@example.com", "supersecret123")
	_, _, refreshB := loginPair(t, "logout@example.com", "supersecret123")

	// Logout kills one session; the other survives.
	if err := services.Revoke(refreshA); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, _, _, err := services.Rotate(refreshA); err != services.ErrInvalidRefreshToken {
		t.Fatalf("logged-out rotate err = %v, want ErrInvalidRefreshToken", err)
	}
	if _, _, _, err := services.Rotate(refreshB); err != nil {
		t.Fatalf("surviving session rotate: %v", err)
	}
	if err := services.Revoke("unknown-token"); err != nil {
		t.Fatalf("revoke unknown must succeed silently: %v", err)
	}

	// Logout-all kills everything.
	_, _, refreshC := loginPair(t, "logout@example.com", "supersecret123")
	user, err := repository.GetUserByEmail("logout@example.com")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if err := services.RevokeAll(user.ID); err != nil {
		t.Fatalf("revoke all: %v", err)
	}
	if _, _, _, err := services.Rotate(refreshB); err != services.ErrInvalidRefreshToken {
		t.Fatalf("post-logout-all err = %v, want ErrInvalidRefreshToken", err)
	}
	if _, _, _, err := services.Rotate(refreshC); err != services.ErrInvalidRefreshToken {
		t.Fatalf("post-logout-all err = %v, want ErrInvalidRefreshToken", err)
	}
}

func TestRefreshHTTP(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	registerVerifiedUser(t, "Refresh HTTP", "refreshhttp@example.com", "supersecret123")

	w := doRequest(t, router, "POST", "/api/auth/login",
		dto.LoginRequest{Email: "refreshhttp@example.com", Password: "supersecret123"}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("login status = %d", w.Code)
	}
	var loginBody struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginBody); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	if loginBody.RefreshToken == "" {
		t.Fatalf("login response missing refresh_token: %s", w.Body.String())
	}

	w = doRequest(t, router, "POST", "/api/auth/refresh",
		dto.RefreshRequest{RefreshToken: loginBody.RefreshToken}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, body %s", w.Code, w.Body.String())
	}

	w = doRequest(t, router, "POST", "/api/auth/refresh",
		dto.RefreshRequest{RefreshToken: "bogus"}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("bogus refresh status = %d, want 401", w.Code)
	}

	w = doRequest(t, router, "POST", "/api/auth/logout",
		dto.RefreshRequest{RefreshToken: loginBody.RefreshToken}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("logout status = %d, body %s", w.Code, w.Body.String())
	}

	w = doRequest(t, router, "POST", "/api/auth/logout-all", nil, loginBody.Token)
	if w.Code != http.StatusOK {
		t.Fatalf("logout-all status = %d, body %s", w.Code, w.Body.String())
	}
	w = doRequest(t, router, "POST", "/api/auth/logout-all", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("logout-all without token status = %d, want 401", w.Code)
	}
}
