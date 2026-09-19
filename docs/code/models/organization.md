# `models/organization.go` + contact/org DTOs

## Models

- `Organization{ID, Name, CreatedAt, UpdatedAt}`.
- `Membership{UserID, OrgID (composite PK), Role, CreatedAt}`.
- Role constants `RoleViewer/RoleMember/RoleAdmin/RoleOwner` + `ValidRole(r)`.
- `Contact.OrgID *string` (nullable, indexed, `omitempty`): nil = personal contact, owner-only.

## DTOs (`dto/org.go`)

`CreateOrgRequest{Name}`, `AddMemberRequest{UserID, Role (oneof)}`, `UpdateMemberRequest{Role}`, plus `OrgEnvelope` / `OrgsEnvelope` / `MemberEnvelope` for Swagger. `dto.ContactFilter` gained `org_id` (query or `X-Org-ID` header); `CreateContactRequest` gained optional `org_id` (handlers overwrite `user_id` with the caller — ownership can't be spoofed).
