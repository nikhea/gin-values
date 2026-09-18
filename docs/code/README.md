# Code docs — one MD per source file

Mirrors the repo layout. Start with `main.md` + `go.mod.md`, then follow the request flow: `routes/` → `middleware/` → `handlers/` → `service/` → `repository/` → `models/`/`dto/`, with `config/`, `utils/`, `jobs/` as cross-cutting support.

## Index

- `main.md`, `go.mod.md`, `cmd/migrate-main.md`, `docs-generated.md`
- `config/`: `logger.md`, `env.md`, `database.md`, `migrate.md`, `auth.md`, `mail.md`
- `models/`: `users.md`, `profile.md`, `contact.md`
- `dto/`: `auth.md`, `user-dtos.md`, `contact.md`, `profile.md`, `pagination.md`, `responses.md`, `import.md`
- `repository/`: `user-repository.md`, `profile-repository.md`, `contact-repository.md`
- `service/`: `auth-service.md`, `user-service.md`, `profile-service.md`, `contact-service.md`, `contact-import-service.md`, `avatar-service.md`
- `handlers/`: `auth-handler.md`, `user-handler.md`, `profile-handler.md`, `contact-handler.md`, `upload-handlers.md`
- `middleware/`: `auth.md`, `request-id.md`
- `routes/`: `routes.md`, `auth-routes.md`, `user-routes.md`, `profile-routes.md`, `contact-routes.md`
- `utils/`: `jwt.md`, `password.md`, `token.md`, `otp.md`, `mailer.md`, `email-templates.md`, `templates-verification.md`, `templates-reset.md`, `file-upload.md`
- `jobs/jobs.md`
- `migrations/`: `000001-users.md` … `000006-contact-avatar.md`
- `tests/`: `setup.md`, `auth.md`, `otp.md`, `email-jobs.md`, `users.md`, `import-upload.md`, `unit.md`

Regenerate Swagger API docs separately with `swag init` (see `docs-generated.md`).
