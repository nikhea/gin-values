# `service/org_service.go` + `repository/org_repository.go`

Organization and membership management. Route middleware owns authorization; services own invariants.

## Service functions

- `CreateOrg(creatorID, req)` — creates org + owner membership, seeds Casbin policies, grants owner; rolls back (delete org) if any step fails.
- `ListOrgs(userID)` — membership join, newest first.
- `AddMember(orgID, userID, role)` — validates role/user/org; `ErrAlreadyMember` on dupes; grants Casbin grouping (row removed again if grant fails).
- `UpdateMemberRole` / `RemoveMember` — `ErrLastOwner` protects the sole owner from demotion/removal.
- `DeleteOrg(orgID)` — row delete (memberships cascade, team contacts detach via `ON DELETE SET NULL`) + Casbin policy removal.

Sentinels: `ErrInvalidRole`, `ErrAlreadyMember`, `ErrLastOwner`. Repository adds org/membership CRUD plus `CountOwnersByRole` and `ListUserOrgIDs`.
