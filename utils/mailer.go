package utils

import (
	"fmt"
	"log"
	"net/smtp"

	"gin-learn/config"
)

// sendMail delivers a plain-text email via SMTP. When credentials are
// missing (local dev), it logs the email and returns nil so auth flows
// never fail just because SMTP isn't configured.
func sendMail(to, subject, body string) error {
	cfg := config.MailConfigFromEnv()
	if cfg.Address == "" || cfg.Password == "" {
		log.Printf("[mailer:logged-only] to=%s subject=%q\n%s\n", to, subject, body)
		return nil
	}

	msg := []byte(fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.FromName, cfg.Address, to, subject, body))

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

// SendVerificationEmail sends the verify-your-email message for auth.
func SendVerificationEmail(to, name, token string) error {
	subject := "Verify your email"
	body := fmt.Sprintf("Hi %s,\n\nThanks for registering. Please verify your email by visiting:\n\n%s\n\nIf you didn't create this account, you can ignore this email.",
		name, VerifyLink(token))
	return sendMail(to, subject, body)
}

// SendPasswordResetEmail sends the password-reset message for auth.
func SendPasswordResetEmail(to, name, token string) error {
	subject := "Reset your password"
	body := fmt.Sprintf("Hi %s,\n\nYou requested a password reset. Set a new password here (valid for 1 hour):\n\n%s\n\nIf you didn't request this, you can ignore this email.",
		name, ResetLink(token))
	return sendMail(to, subject, body)
}
