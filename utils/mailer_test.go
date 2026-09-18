package utils

import (
	"strings"
	"testing"
)

// With no SMTP credentials the mailer must log-and-succeed so auth
// flows never fail in dev/test environments.
func TestMailerDisabledLogsAndSucceeds(t *testing.T) {
	t.Setenv("EMAIL_ADDRESS", "")
	t.Setenv("EMAIL_PASSWORD", "")

	if err := SendVerificationEmail("user@example.com", "User", "tok123"); err != nil {
		t.Fatalf("SendVerificationEmail disabled: %v", err)
	}
	if err := SendPasswordResetEmail("user@example.com", "User", "tok123"); err != nil {
		t.Fatalf("SendPasswordResetEmail disabled: %v", err)
	}
}

func TestAuthLinks(t *testing.T) {
	t.Setenv("APP_URL", "http://example.test:8080")

	verify := VerifyLink("abc")
	if !strings.Contains(verify, "/api/auth/verify?token=abc") {
		t.Fatalf("unexpected verify link: %q", verify)
	}
	reset := ResetLink("xyz")
	if !strings.Contains(reset, "/api/auth/reset-password?token=xyz") {
		t.Fatalf("unexpected reset link: %q", reset)
	}
}
