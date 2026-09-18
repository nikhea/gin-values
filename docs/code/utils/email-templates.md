# `utils/email_templates.go`

Embedded HTML email templates and the builders shared by direct sends and River jobs.

## Mechanism

`//go:embed templates/*.html` + `template.Must(ParseFS(...))` — templates compile at startup (a broken template fails fast). Data structs `verificationData{Name, OTP, VerifyLink, AppName}` and `resetData` feed `renderVerificationHTML` / `renderResetHTML`.

## Builders

- `BuildVerificationEmail(name, token, otp) (subject, text, html, err)` — subject `"Verify your email"`; text includes the OTP code + link; HTML renders `templates/verification.html`.
- `BuildPasswordResetEmail(name, token) (subject, text, html, err)` — 1-hour link; renders `templates/reset.html`.

Consumed by `utils/mailer.go` (`Send*Email`) and `jobs/jobs.go` (`Enqueue*Email`). Brand name is the `emailAppName = "Grip"` constant.
