# Code docs — one MD per source file

Mirrors the repo layout. Start with `main.md` + `go.mod.md`, then follow the request flow: `routes/` → `middleware/` → `handlers/` → `service/` → `repository/` → `models/`/`dto/`, with `config/`, `utils/`, `jobs/` as cross-cutting support.

## Index

- `main.md`, `go.mod.md`, `makefile.md`, `ci-lint.md`, `Dockerfile.md`, `docker-compose.md`, `cmd/migrate-main.md`, `docs-generated.md`
- `PRODUCTION.md` (top-level `docs/`) — production runbook: secrets, TLS, migrations, backups, jobs, observability
- `config/`: `logger.md`, `env.md`, `database.md`, `migrate.md`, `auth.md`, `mail.md`, `redis.md`
- `models/`: `users.md`, `profile.md`, `contact.md`, `organization.md`
- `dto/`: `auth.md`, `user-dtos.md`, `contact.md`, `profile.md`, `pagination.md`, `responses.md`, `import.md`
- `repository/`: `user-repository.md`, `profile-repository.md`, `contact-repository.md`
- `service/`: `auth-service.md`, `refresh-service.md`, `user-service.md`, `profile-service.md`, `contact-service.md`, `contact-import-service.md`, `avatar-service.md`, `org-service.md`
- `handlers/`: `auth-handler.md`, `user-handler.md`, `profile-handler.md`, `contact-handler.md`, `upload-handlers.md`, `health-handler.md`, `errors-probes.md`, `org-handler.md`
- `middleware/`: `auth.md`, `request-id.md`, `rate-limit.md`, `authorize.md`
- `routes/`: `routes.md`, `auth-routes.md`, `user-routes.md`, `profile-routes.md`, `contact-routes.md` (org routes covered in `handlers/org-handler.md`)
- `utils/`: `jwt.md`, `password.md`, `token.md`, `otp.md`, `mailer.md`, `email-templates.md`, `templates-verification.md`, `templates-reset.md`, `file-upload.md`, `cache.md`
- `jobs/jobs.md`
- `migrations/`: `000001-users.md` … `0011-notifications.md`
- `authz/authz.md` — Casbin domains model, role policies, version-pairing rule
- `audit/audit.md` — append-only trail, writer helpers, admin listing
- `notifications/notifications.md` — inbox model, River worker, emitters, endpoints
- `tests/`: `setup.md`, `auth.md`, `otp.md`, `email-jobs.md`, `users.md`, `users-pagination.md`, `import-upload.md`, `orgs-rbac.md`, `audit.md`, `notifications.md`, `unit.md`

Regenerate Swagger API docs separately with `swag init` (see `docs-generated.md`).
