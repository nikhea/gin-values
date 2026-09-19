# `tests/orgs_rbac_test.go` — organizations + RBAC tests

- `TestOrgLifecycle` — create (creator is owner in DB + Casbin), list visibility (stranger sees nothing), owner HTTP delete.
- `TestMemberManagement` — add / duplicate / invalid role; member 403 on member-admin and outsider 403 on org delete via HTTP; last-owner demotion/removal guards.
- `TestContactRBACMatrix` — owner creates team contact (caller-spoof attempt lands on caller); read matrix owner/member/viewer 200, stranger 404; write matrix viewer 403, stranger 404, member 200; org-scoped list member 200 / stranger 403; unscoped list hides others' team contacts.
