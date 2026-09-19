# `authz/authz.go`

Casbin authorization layer: RBAC with domains, one domain per organization ID. Policies persist in `casbin_rule` (gorm-adapter); the model is embedded — no external file at runtime.

## Model

`r = sub, dom, obj, act` / `g = _, _, _` / allow-override. Resources: `contacts`, `members`, `orgs`; actions: `read`, `write`, `delete` (`regexMatch` on actions, `keyMatch` on objects).

## Role policies (seeded per org)

owner: everything on all three resources; admin: contacts+members read/write, orgs read; member: contacts read/write; viewer: contacts read.

⚠️ `keyMatch` is **prefix** matching, not regex — wildcard resources must be enumerated per resource (`".*"` never matches under `keyMatch`; only actions use regex).

## Exports

- `Init(db)` — adapter + enforcer + `LoadPolicy`; called in `main()` and `tests/setup_test.go`. Enforcer is safe for concurrent `Enforce`.
- `Enforce(sub, dom, obj, act) (bool, error)` — errors fail closed.
- `SeedOrgPolicies(orgID)` / `RemoveOrgPolicies(orgID)` / `GrantRole` / `RevokeRole` (filtered grouping removal, role-agnostic).

Deps: `casbin/casbin/v3` + `gorm-adapter/v3` — the adapter major must track the Casbin major (v3↔v3); mixing v2 enforcer with v3 adapter panics at startup.
