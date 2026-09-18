# `utils/mailer.go`

SMTP delivery for auth emails (multipart plain-text + HTML). Never fails callers when unconfigured.

## Functions

- `SendMail(to, subject, textBody, htmlBody) error` — without `EMAIL_ADDRESS`/`EMAIL_PASSWORD` logs (`slog.Info "mailer logged-only"`) and returns nil; otherwise sends `multipart/alternative` via Gmail-compatible SMTP (`smtp.PlainAuth`).
- `VerifyLink(token)` / `ResetLink(token)` — build `APP_URL/api/auth/...` links for templates and text bodies.
- `SendVerificationEmail(to, name, token, otp)` / `SendPasswordResetEmail(to, name, token)` — thin wrappers over the builders in `utils/email_templates.go` + `SendMail`. Mostly superseded by River jobs (`jobs/jobs.go`) but kept for direct sends/tests.

Content (subjects, bodies, templates) lives in `utils/email_templates.md`. Behavior tested in `tests/mailer_test.go`.
