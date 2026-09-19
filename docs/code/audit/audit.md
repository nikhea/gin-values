# `audit/audit.go` + `repository/audit_repository.go`

Append-only audit trail of state-changing API calls.

## Writing

- `audit.Log(c, action, resourceType, resourceID, meta)` — actor/org/request-ID/IP from context; for anonymous routes actor is nil.
- `audit.LogFor(c, actorID, ...)` — same but with an explicit actor (register/login/verify, which run outside `AuthRequired`); still captures IP + request ID.
- Both swallowErrors (slog only) — auditing can never break a request.

Handlers call these after every successful mutation: auth events, user/profile/contact CRUD, imports (with counts), avatars, org + member changes. Action names are `models.Audit*` constants (`users.register`, `contacts.create`, …) so the log stays queryable.

## Reading

`repository.CreateAuditLog` (no update/delete exists on purpose), `ListAuditLogs(AuditFilter)` (actor/action/resource/org/from/to + pagination, newest first, self-normalizing), `GetAuditLogByID(id, orgID)` (org-scoped, 404 outside scope).

Listing endpoints (`handlers/audit_handler.go`, `routes/audit_routes.go`): `GET /audit-logs/` + `GET /audit-logs/:id`, both behind `RequireOrgAccess(members, write)` — org admins only.
