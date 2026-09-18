package utils

import (
	"fmt"
	"log/slog"
	"net/smtp"

	"gin-learn/config"
)

// SendMail delivers a multipart (plain-text + HTML) email via SMTP.
// When credentials are missing (local dev), it logs the email and
// returns nil so auth flows never fail just because SMTP isn't configured.
func SendMail(to, subject, textBody, htmlBody string) error {
	cfg := config.MailConfigFromEnv()
	if cfg.Address == "" || cfg.Password == "" {
		slog.Info("mailer logged-only", "to", to, "subject", subject, "body", textBody)
		return nil
	}

	boundary := "gin-learn-boundary"
	header := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=%s\r\n\r\n",
		cfg.FromName, cfg.Address, to, subject, boundary)
	partText := fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n", boundary, textBody)
	partHTML := fmt.Sprintf("--%s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n", boundary, htmlBody)
	msg := []byte(header + partText + partHTML + "--" + boundary + "--\r\n")

	auth := smtp.PlainAuth("", cfg.Address, cfg.Password, cfg.Host)
	if err := smtp.SendMail(cfg.Host+":"+cfg.Port, auth, cfg.Address, []string{to}, msg); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

// VerifyLink builds the frontend/API link a user clicks to verify email.
func VerifyLink(token string) string {
	return fmt.Sprintf("%s/api/auth/verify?token=%s", config.AppURL(), token)
}

// ResetLink builds the link a user clicks to set a new password.
func ResetLink(token string) string {
	return fmt.Sprintf("%s/api/auth/reset-password?token=%s", config.AppURL(), token)
}

// SendVerificationEmail sends the verify-your-email message for auth,
// including the OTP code and a verification link.
func SendVerificationEmail(to, name, token, otp string) error {
	subject, text, html, err := BuildVerificationEmail(name, token, otp)
	if err != nil {
		return err
	}
	return SendMail(to, subject, text, html)
}

// SendPasswordResetEmail sends the password-reset message for auth.
func SendPasswordResetEmail(to, name, token string) error {
	subject, text, html, err := BuildPasswordResetEmail(name, token)
	if err != nil {
		return err
	}
	return SendMail(to, subject, text, html)
}
