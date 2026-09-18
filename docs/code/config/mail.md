# `config/mail.go`

Resolves SMTP settings from the environment for auth emails.

## Exports

- `MailConfig{Address, Password, Host, Port, FromName}` — credential + server bundle.
- `MailConfigFromEnv() MailConfig` — maps `EMAIL_SERVICE=Gmail` (or `Google`) to `smtp.gmail.com:587`; otherwise honors `EMAIL_HOST`/`EMAIL_PORT`. Trims the address; `FromName` is `"Gin Learn"`.
- `MailEnabled() bool` — true only when address + password are both set; otherwise warns that emails will be logged only.

## Depends on

Environment only (see `config/env.md`). Consumed by `utils/mailer.go`.

## Notes

- Missing credentials never break auth flows — the mailer logs instead of sending.
