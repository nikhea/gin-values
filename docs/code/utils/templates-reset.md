# `utils/templates/reset.html`

HTML email for password reset. Placeholders: `{{ .AppName }}`, `{{ .Name }}`, `{{ .ResetLink }}` (1-hour button link).

Same card layout/branding as the verification template (see `utils/templates-verification.md`). Rendered by `renderResetHTML`; plain-text twin built in `BuildPasswordResetEmail`.
