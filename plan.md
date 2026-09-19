# Grip — Feature Plan

Seven items, ordered by dependency.

## Status

- [x] 1. Refresh-token rotation — **done** (migration `000007`, hash-only storage, rotation + reuse-kills-family, `/refresh`, `/logout`, `/logout-all`, 15m access / 30d refresh, tests + docs + live-verified)
- [x] 2. Soft delete — **done** (migration `000008`, `DeletedAt` on User/Profile/Contact, partial unique indexes, `HardDelete*` repo funcs, `soft_deleted` DELETE flag, tests + docs + live-verified)
- [x] 3. Teams/organizations + Casbin RBAC — **done** (casbin/v3 + gorm-adapter/v3, migration `000009`, domains model, owner/admin/member/viewer, org + member endpoints, contact scoping, tests + docs + live-verified)
- [x] 4. Audit log — **done** (migration `0010`, append-only table, handler hooks on all mutations, org-admin listing, tests + docs + live-verified)
- [x] 5. In-app notifications — **done** (migration `0011`, River `notify` worker, welcome/password/import/org emitters, self-scoped inbox endpoints, tests + docs + live-verified)
- [x] 6. Pagination + filtering parity — **done** (`dto.UserFilter`, `repository.ListUsers`, `meta` envelope on `GET /users/`, legacy `GetUsers` kept, tests + docs + live-verified)

## Execution order

`1 → 2 → 3 → 4+5+6` (4/5/6 are independent slices once 2+3 land; 7 runs throughout)

---

## 1. Refresh-token rotation

**Design:** short-lived access JWT (15m via new `ACCESS_TTL_MINUTES`, default keeps working) + opaque rotating refresh tokens. New `refresh_tokens` table (`id`, `user_id` FK cascade, `token_hash` SHA-256 unique, `expires_at` 30d via `REFRESH_TTL_DAYS`, `revoked_at`, `replaced_by`, `created_at`). Rotation on every use; **reuse detection**: presenting an already-rotated token revokes the whole family (theft response).

- **Files:** `migrations/000007_refresh_tokens`, `models/refresh_token.go` (json:`-` everywhere), `service/refresh_service.go` (`IssuePair`, `Rotate`, `Revoke`, `RevokeAll`), `repository/refresh_repository.go`
- **Endpoints:** `POST /auth/refresh {refresh_token}` → new pair (200) / reuse → 401 + family revoked; `POST /auth/logout` (revoke one, protected); `POST /auth/logout-all` (protected). `Login` response gains `refresh_token`.
- **Notes:** reuse existing `GenerateSecureToken` (32B) + SHA-256 pattern from OTP code; hash-only storage means a DB leak can't mint sessions.
- **Tests:** rotate-happy-path, reuse-triggers-family-revoke, expired rejected, logout/logout-all.

## 2. Soft delete

- **Models:** add `DeletedAt gorm.DeletedTime gorm:"index"` to `User`, `Profile`, `Contact` (GORM then auto-filters + turns `Delete` soft). **Gotcha:** `users.email` unique index + soft delete breaks re-registration → migration replaces it with partial unique index `WHERE deleted_at IS NULL` (same for `profiles.user_id`). New migration `000008_soft_delete`.
- **Repo:** `DeleteUser/Profile/Contact` stay as-is (now soft); add `HardDelete*` for admin/purge use + `FindIncludingDeleted` where needed (login on deleted account 401s as "not found" via the existing `ErrRecordNotFound` path — stays correct).
- **API:** `DELETE` responses gain a `soft_deleted: true` note; `?hard=true` admin-only param only if Casbin (item 3) lands — otherwise skip hard-delete endpoint.
- **Tests:** delete → get 404s but row exists with `deleted_at`; re-register same email works; hard delete removes row.

## 3. Teams/organizations + Casbin RBAC

- **Deps:** `github.com/casbin/casbin/v2` + `github.com/casbin/gorm-adapter/v3` (new). Casbin model embedded as Go string (no file dependency): RBAC **with domains**, domains = org IDs:
  ```
  [request_definition] r = sub, dom, obj, act
  [policy_definition]  p = sub, dom, obj, act
  [role_definition]    g = _, _, _
  [policy_effect]      e = some(where (p.eft == allow))
  [matchers]           m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
  ```
  Roles: `owner` (all), `admin` (all minus org delete), `member` (contacts CRUD), `viewer` (read). Seeded per org on creation.
- **Models/migration `000009_orgs`:** `organizations(id, name, created_at)`, `memberships(user_id, org_id, role, unique(user_id,org_id))`; `contacts` gains nullable `org_id` FK (NULL = personal contact, owner-only).
- **Enforcement:** new `middleware/authorize.go` — `RequireOrgAccess(obj, act)` resolves org from `X-Org-ID` header (or query), loads `middleware.GetUserID`, calls `authz.Enforce(userID, orgID, resource, action)` → 403 on deny; personal resources keep the existing owner check (`:id == userID` or contact.UserID match). New `authz/` package: `Init(db)` singleton enforcer + gorm-adapter (auto-creates `casbin_rule`), helpers `GrantRole`, `RevokeRole`, `AddOrgPolicies`.
- **Endpoints** (`handlers/org_handler.go`, protected): `POST /orgs` (creator → owner), `GET /orgs` (mine), `POST /orgs/:id/members {user_id, role}` (admin+), `PUT /orgs/:id/members/:uid` (admin+), `DELETE /orgs/:id/members/:uid` (admin+, can't remove last owner), `DELETE /orgs/:id` (owner).
- **Existing routes change:** contacts handlers take optional org scope; `GET /contacts` filters by org membership when `X-Org-ID` present.
- **Tests:** owner/admin/member/viewer matrix on contacts + member management guards.

## 4. Audit log

- **Model/migration `000010_audit_logs`:** `audit_logs(id, actor_id index, action, resource_type, resource_id, org_id nullable index, ip, request_id, metadata JSONB, created_at index)`. No updates/deletes ever (append-only; no repository delete function on purpose).
- **Mechanism:** `audit/Log(c, action, resource...)` helper reading `GetUserID`, `ClientIP`, `requestID` from context; called in services on every mutation (auth events, user/profile/contact CRUD, org membership changes, imports). Synchronous single-row insert (ordering guarantees over async).
- **Endpoint:** `GET /audit-logs` (admin-only via Casbin item 3, filters: `actor_id`, `action`, `resource_type`, `from/to`, paginated) + `GET /audit-logs/:id`.
- **Tests:** action creates entry with correct actor/resource/request_id; list filters; non-admin 403.

## 5. In-app notifications

- **Model/migration `000011_notifications`:** `notifications(id, user_id index, type, title, body, data JSONB, read_at nullable index, created_at)`.
- **Producers:** reuse River — new `NotifyArgs{UserID, Type, Title, Body, Data}` + `NotifyWorker` in `jobs/` (same pattern as `SendEmailWorker`). Emit on: welcome (post-verify), password-changed, import finished (with counts), added-to-org, role-changed.
- **Endpoints** (`handlers/notification_handler.go`, protected, self-scoped by `GetUserID`): `GET /notifications` (paginated, `?unread_only`), `GET /notifications/unread-count`, `PATCH /notifications/:id/read`, `POST /notifications/read-all`.
- **Tests:** job creates row; unread count math; users can't read others' (404/empty); read-all idempotent.

## 6. Pagination + filtering parity (users list)

- Mirror the contacts pattern: `dto.UserFilter{search (name/email ILIKE), page, page_size}` + `repository.ListUsers` returning `(users, total)` + service change + `GET /users/` returns `{message, user, meta}` (`UsersEnvelope` gains `Meta PageMeta` — additive, backward compatible).
- **Tests:** search match/no-match, page boundaries, meta totals.

## Cross-cutting for all items

- **Migrations** `000007`–`000011` (up/down pairs, golang-migrate convention); boot applies them, River untouched.
- **Swagger:** annotations + `example` tags on every new DTO/endpoint, `swag init`, BearerAuth on protected.
- **Tests** (`tests/` convention): `refresh_test.go`, `soft_delete_test.go`, `orgs_rbac_test.go`, `audit_test.go`, `notifications_test.go`, `users_pagination_test.go` (+ unit tests for key builders/validators where pure).
- **Docs:** per-file `docs/code/*.md` + index updates (repo convention), README endpoint table + `.env.example` additions (`ACCESS_TTL_MINUTES`, `REFRESH_TTL_DAYS`).
- **Gate:** `make check` must stay green throughout; each item ends with a live smoke (boot → curl flow) before moving on.

## Open questions (unresolved)

1. **Refresh storage:** Postgres table (recommended — transactional rotation, matches existing patterns) vs Redis (faster, but revocation/family logic gets harder)?
2. **Org scoping for existing data:** contacts with `org_id = NULL` stay personal/owner-only — acceptable, or do you want a default org per user?
3. **Logout semantics:** revoke presented token only, or also "logout all sessions" from the start? (Priced in; say the word to cut it.)
4. **Order:** refresh → soft-delete → teams/casbin → audit → notifications → pagination. Confirm or re-sequence (notifications is the most user-visible if you want it earlier).
