package utils

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*.html
var emailTemplates embed.FS

var (
	verificationTmpl = template.Must(template.ParseFS(emailTemplates, "templates/verification.html"))
	resetTmpl        = template.Must(template.ParseFS(emailTemplates, "templates/reset.html"))
)

const emailAppName = "Gin Learn"

type verificationData struct {
	AppName    string
	Name       string
	OTP        string
	VerifyLink string
}

type resetData struct {
	AppName   string
	Name      string
	ResetLink string
}

func renderVerificationHTML(name, token, otp string) (string, error) {
	var buf bytes.Buffer
	err := verificationTmpl.Execute(&buf, verificationData{
		AppName:    emailAppName,
		Name:       name,
		OTP:        otp,
		VerifyLink: VerifyLink(token),
	})
	return buf.String(), err
}

func renderResetHTML(name, token string) (string, error) {
	var buf bytes.Buffer
	err := resetTmpl.Execute(&buf, resetData{
		AppName:   emailAppName,
		Name:      name,
		ResetLink: ResetLink(token),
	})
	return buf.String(), err
}

// BuildVerificationEmail composes the verify-your-email subject, plain-text
// fallback and HTML body. Used both by direct sends and River email jobs.
func BuildVerificationEmail(name, token, otp string) (subject, text, html string, err error) {
	subject = "Verify your email"
	text = fmt.Sprintf("Hi %s,\n\nThanks for registering. Your verification code is %s (valid for 10 minutes).\n\nOr verify here:\n\n%s\n\nIf you didn't create this account, you can ignore this email.",
		name, otp, VerifyLink(token))
	html, err = renderVerificationHTML(name, token, otp)
	return subject, text, html, err
}

// BuildPasswordResetEmail composes the password-reset subject, plain-text
// fallback and HTML body. Used both by direct sends and River email jobs.
func BuildPasswordResetEmail(name, token string) (subject, text, html string, err error) {
	subject = "Reset your password"
	text = fmt.Sprintf("Hi %s,\n\nYou requested a password reset. Set a new password here (valid for 1 hour):\n\n%s\n\nIf you didn't request this, you can ignore this email.",
		name, ResetLink(token))
	html, err = renderResetHTML(name, token)
	return subject, text, html, err
}
