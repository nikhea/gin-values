# `tests/audit_test.go` — audit trail tests

- `TestMutationsWriteAuditRows` — HTTP register writes a row with actor + request ID; contact create/update write rows with resource IDs; import writes a summary row with counts in metadata.
- `TestAuditListingGatedByOrgAdmin` — owner lists 200, member/outsider 403; rows never leak across orgs (direct repository check + cross-org detail 404, own-org detail 200).
- `doRequestOrg` helper adds the `X-Org-ID` scope header.
