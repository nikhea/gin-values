# `utils/templates/verification.html`

HTML email for address verification. Placeholders: `{{ .AppName }}` (header/footer brand), `{{ .Name }}` (greeting), `{{ .OTP }}` (large letter-spaced code block), `{{ .VerifyLink }}` (fallback button).

Layout: indigo (`#4f46e5`) header, 480px centered card, dashed code box, button, footer note. Styling is fully inline for email-client compatibility. Rendered by `renderVerificationHTML` (see `utils/email_templates.md`); plain-text twin is built alongside it in `BuildVerificationEmail`.
