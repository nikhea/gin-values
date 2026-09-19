package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"grip/config"
	"grip/dto"
	"grip/models"
	"grip/repository"
	"grip/routes"
	services "grip/service"
	"grip/utils"

	"github.com/gin-gonic/gin"
)

func doRequest(t *testing.T, router *gin.Engine, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func registerVerifiedUser(t *testing.T, name, email, password string) (*models.User, string) {
	t.Helper()
	user, err := services.Register(dto.RegisterRequest{
		Name: name, Email: email, Password: password, Age: 30,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	stored, err := repository.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("load verification token: %v", err)
	}
	if _, err := services.VerifyEmail(stored.VerificationToken); err != nil {
		t.Fatalf("verify: %v", err)
	}
	loggedIn, token, _, err := services.Login(dto.LoginRequest{Email: email, Password: password})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	_ = loggedIn
	_ = user
	return stored, token
}

func TestRegisterVerifyLoginFlow(t *testing.T) {
	requireTestDB(t)

	user, err := services.Register(dto.RegisterRequest{
		Name: "Integ Test", Email: "integ@example.com", Password: "supersecret123", Age: 28,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if user.PasswordHash == "" || user.EmailVerified {
		t.Fatal("expected password hash set and email unverified")
	}

	// Duplicate email.
	if _, err := services.Register(dto.RegisterRequest{
		Name: "Integ Test", Email: "integ@example.com", Password: "supersecret123", Age: 28,
	}); err != services.ErrEmailTaken {
		t.Fatalf("duplicate register err = %v, want ErrEmailTaken", err)
	}

	// Login before verification is forbidden.
	if _, _, _, err := services.Login(dto.LoginRequest{Email: "integ@example.com", Password: "supersecret123"}); err != services.ErrEmailNotVerified {
		t.Fatalf("pre-verify login err = %v, want ErrEmailNotVerified", err)
	}

	// Wrong password.
	if _, _, _, err := services.Login(dto.LoginRequest{Email: "integ@example.com", Password: "wrongpass1"}); err != services.ErrInvalidCredentials {
		t.Fatalf("wrong password err = %v, want ErrInvalidCredentials", err)
	}

	// Verify with bad token.
	if _, err := services.VerifyEmail("nope"); err != services.ErrInvalidToken {
		t.Fatalf("bad token verify err = %v, want ErrInvalidToken", err)
	}

	stored, err := repository.GetUserByEmail("integ@example.com")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if _, err := services.VerifyEmail(stored.VerificationToken); err != nil {
		t.Fatalf("verify: %v", err)
	}

	loggedIn, token, _, err := services.Login(dto.LoginRequest{Email: "integ@example.com", Password: "supersecret123"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if !loggedIn.EmailVerified || token == "" {
		t.Fatal("expected verified user and token")
	}
	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Subject != loggedIn.ID {
		t.Fatalf("token subject = %q, want %q", claims.Subject, loggedIn.ID)
	}
}

func TestPasswordResetFlow(t *testing.T) {
	requireTestDB(t)
	registerVerifiedUser(t, "Reset Me", "reset@example.com", "oldpassword1")

	if err := services.ForgotPassword("unknown@example.com"); err != nil {
		t.Fatalf("forgot unknown email must succeed: %v", err)
	}
	if err := services.ForgotPassword("reset@example.com"); err != nil {
		t.Fatalf("forgot: %v", err)
	}

	stored, err := repository.GetUserByEmail("reset@example.com")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if stored.ResetToken == "" {
		t.Fatal("expected reset token to be set")
	}

	if _, err := services.ResetPassword("bad-token", "newpassword1"); err != services.ErrInvalidToken {
		t.Fatalf("bad token reset err = %v, want ErrInvalidToken", err)
	}
	if _, err := services.ResetPassword(stored.ResetToken, "newpassword1"); err != nil {
		t.Fatalf("reset: %v", err)
	}

	if _, _, _, err := services.Login(dto.LoginRequest{Email: "reset@example.com", Password: "oldpassword1"}); err != services.ErrInvalidCredentials {
		t.Fatalf("old password err = %v, want ErrInvalidCredentials", err)
	}
	if _, _, _, err := services.Login(dto.LoginRequest{Email: "reset@example.com", Password: "newpassword1"}); err != nil {
		t.Fatalf("login with new password: %v", err)
	}

	// Token is single-use.
	if _, err := services.ResetPassword(stored.ResetToken, "anotherpass1"); err != services.ErrInvalidToken {
		t.Fatalf("reused token err = %v, want ErrInvalidToken", err)
	}
}

func TestAuthHTTPAndProtectedRoutes(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()

	// Public: register over HTTP.
	w := doRequest(t, router, "POST", "/api/auth/register", dto.RegisterRequest{
		Name: "HTTP User", Email: "http@example.com", Password: "supersecret123", Age: 30,
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register status = %d, body %s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), "password_hash") || strings.Contains(w.Body.String(), "verification_token") {
		t.Fatalf("response leaks secrets: %s", w.Body.String())
	}

	// Protected without token.
	w = doRequest(t, router, "GET", "/api/users/", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("no-token status = %d, want 401", w.Code)
	}

	// Verify + login to get a token.
	stored, err := repository.GetUserByEmail("http@example.com")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if _, err := services.VerifyEmail(stored.VerificationToken); err != nil {
		t.Fatalf("verify: %v", err)
	}
	w = doRequest(t, router, "POST", "/api/auth/login",
		dto.LoginRequest{Email: "http@example.com", Password: "supersecret123"}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("login status = %d, body %s", w.Code, w.Body.String())
	}
	var loginBody struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginBody); err != nil || loginBody.Token == "" {
		t.Fatalf("login response has no token: %s", w.Body.String())
	}

	// /me with token.
	w = doRequest(t, router, "GET", "/api/auth/me", nil, loginBody.Token)
	if w.Code != http.StatusOK {
		t.Fatalf("/me status = %d, body %s", w.Code, w.Body.String())
	}

	// Protected resource with token.
	w = doRequest(t, router, "GET", "/api/users/", nil, loginBody.Token)
	if w.Code != http.StatusOK {
		t.Fatalf("protected users status = %d, body %s", w.Code, w.Body.String())
	}

	// Tampered token.
	w = doRequest(t, router, "GET", "/api/users/", nil, loginBody.Token+"x")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("tampered token status = %d, want 401", w.Code)
	}
}

func TestLegacyUserWithoutPasswordCannotLogin(t *testing.T) {
	requireTestDB(t)
	// Simulate a row created before auth (no password hash).
	legacy := &models.User{ID: "legacy-1", Name: "Legacy", Email: "legacy@example.com", Age: 40}
	if err := config.DB.Create(legacy).Error; err != nil {
		t.Fatalf("create legacy user: %v", err)
	}
	if _, _, _, err := services.Login(dto.LoginRequest{Email: "legacy@example.com", Password: "whatever12"}); err != services.ErrInvalidCredentials {
		t.Fatalf("legacy login err = %v, want ErrInvalidCredentials", err)
	}
}
