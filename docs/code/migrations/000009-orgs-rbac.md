# Migration `000009` — orgs, memberships, contact scoping, Casbin table

- **Up:** `organizations(id PK, name, timestamps)`; `memberships(user_id, org_id, role, PK(user_id,org_id))` with FKs cascading off users/orgs + index on `org_id`; `contacts.org_id` nullable FK → organizations `ON DELETE SET NULL` (team contacts detach to personal on org delete) + index; `casbin_rule` policy table (+ index) for gorm-adapter (also auto-created by the adapter).
- **Down:** drops in reverse order.
