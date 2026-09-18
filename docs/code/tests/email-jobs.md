# `tests/email_jobs_test.go` — River enqueue tests

Asserts services actually queue emails (River client runs against the test DB via `requireTestDB`).

- `TestRegisterEnqueuesVerificationEmail` — one `send_email` row appears in `river_job`, args JSON contains the recipient.
- `TestForgotPasswordEnqueuesResetEmail` — same for the reset flow.
- Helpers `riverJobCount(kind)` and `latestRiverJobArgs(kind)` query `river_job` through `config.DB`.
