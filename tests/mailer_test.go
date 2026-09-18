package tests

import (
	"strings"
	"testing"

	"gin-learn/utils"
)

// With no SMTP credentials the mailer must log-and-succeed so auth
// flows never fail in dev/test environments.
func TestMailerDisabledLogsAndSucceeds(t *testing.T) {
	t.Setenv("EMAIL_ADDRESS", "")
	t.Setenv("EMAIL_PASSWORD", "")

	if err := utils.SendVerificationEmail("user@example.com", "User", "tok123", "123456"); err != nil {
		t.Fatalf("SendVerificationEmail disabled: %v", err)
	}
	if err := utils.SendPasswordResetEmail("user@example.com", "User", "tok123"); err != nil {
		t.Fatalf("SendPasswordResetEmail disabled: %v", err)
	}
}

func TestAuthLinks(t *testing.T) {
	t.Setenv("APP_URL", "http://example.test:8080")

	verify := utils.VerifyLink("abc")
	if !strings.Contains(verify, "/api/auth/verify?token=abc") {
		t.Fatalf("unexpected verify link: %q", verify)
	}
	reset := utils.ResetLink("xyz")
	if !strings.Contains(reset, "/api/auth/reset-password?token=xyz") {
		t.Fatalf("unexpected reset link: %q", reset)
	}
}

func TestBuildEmails(t *testing.T) {
	subject, text, html, err := utils.BuildVerificationEmail("Kaige", "tok123", "482914")
	if err != nil {
		t.Fatalf("build verification: %v", err)
	}
	if subject == "" || !strings.Contains(text, "482914") || !strings.Contains(text, "tok123") {
		t.Fatalf("bad verification text: %q %q", subject, text)
	}
	if !strings.Contains(html, "482914") || !strings.Contains(html, "tok123") || !strings.Contains(html, "Kaige") {
		t.Fatal("verification HTML missing code, link or name")
	}

	subject, text, html, err = utils.BuildPasswordResetEmail("Kaige", "tok456")
	if err != nil {
		t.Fatalf("build reset: %v", err)
	}
	if subject == "" || !strings.Contains(text, "tok456") {
		t.Fatalf("bad reset text: %q %q", subject, text)
	}
	if !strings.Contains(html, "tok456") || !strings.Contains(html, "Kaige") {
		t.Fatal("reset HTML missing link or name")
	}
}

func TestSendMailMultipartDisabled(t *testing.T) {
	t.Setenv("EMAIL_ADDRESS", "")
	t.Setenv("EMAIL_PASSWORD", "")
	if err := utils.SendMail("u@example.com", "Subj", "text", "<p>html</p>"); err != nil {
		t.Fatalf("SendMail disabled: %v", err)
	}
	if err := utils.SendVerificationEmail("u@example.com", "Kaige", "tok", "123456"); err != nil {
		t.Fatalf("SendVerificationEmail disabled: %v", err)
	}
	if err := utils.SendPasswordResetEmail("u@example.com", "Kaige", "tok"); err != nil {
		t.Fatalf("SendPasswordResetEmail disabled: %v", err)
	}
}
