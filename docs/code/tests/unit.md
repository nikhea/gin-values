# Unit tests in `tests/` (no DB)

By project convention all tests live in `tests/` (`package tests`), importing the tested package. Run with `go test ./tests/ -run <Name> -v`.

- `password_test.go` — bcrypt round-trip, wrong password, salting, empty hash (`utils/password.md`).
- `jwt_test.go` — token round-trip + claims, tampered/malformed/empty, wrong secret, expired (`utils/jwt.md`).
- `token_test.go` — length, hex validity, uniqueness (`utils/token.md`).
- `otp_utils_test.go` — OTP digit format, uniqueness, hash/check (`utils/otp.md`). Template rendering itself is covered via `TestBuildEmails` (unexported render funcs aren't reachable from `tests`).
- `mailer_test.go` — log-only mode without SMTP creds, link formats, builders (subject/text/HTML), multipart disabled send (`utils/mailer.md`, `utils/email_templates.md`).
- `middleware_test.go` — `AuthRequired` table test (missing/empty/wrong-scheme/garbage/wrong-secret/expired → 401; `Bearer` and bare token → 200 with context propagated) + helpers without middleware (`middleware/auth.md`).
- `contact_test.go` — every `ValidateContactValue` branch incl. `*ValidationError` type (`dto/contact.md`).
- `jobs_test.go` — `Kind()` string, worker delivery (no DB), enqueue without setup errors (`jobs/jobs.md`).
