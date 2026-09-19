# `handlers/org_handler.go` + `routes/org_routes.go`

Org endpoints (all JWT-protected): `POST /orgs/` (creator → owner), `GET /orgs/` (mine), `POST /orgs/:id/members`, `PUT /orgs/:id/members/:uid`, `DELETE /orgs/:id/members/:uid` (members group needs `members:write`), `DELETE /orgs/:id` (needs `orgs:delete`).

`routes/org_routes.go` mounts everything under protected `/api/orgs`; an `orgScopeFromPath` step copies `:id` into `X-Org-ID` when the caller didn't scope explicitly, so member routes and org delete don't need a redundant header. Errors via `writeOrgError`: 400 invalid role/last-owner, 404 missing, 409 duplicate member, 500 sanitized.
