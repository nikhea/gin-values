# `middleware/authorize.go`

Route-level Casbin enforcement for organization resources.

## Exports

- `OrgIDFromRequest(c)` — org scope from `X-Org-ID` header (preferred) or `?org_id=`.
- `GetOrgID(c)` — ID stored by the middleware (the `req.orgId` equivalent).
- `RequireOrgAccess(obj, act)` — runs after `AuthRequired`: 401 without identity, 400 without org scope, 403 on deny or enforcer error (fail-closed), stores org ID on success.

Wired in `routes/org_routes.go` (`members` → write, org delete → delete). Tested via HTTP in `tests/orgs_rbac_test.go` (member 403, outsider 403, owner OK).
