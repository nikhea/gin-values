// Package authz is the Casbin authorization layer: RBAC with domains,
// where each organization ID is a policy domain. Policies live in the
// casbin_rule table (gorm-adapter); the model is embedded so no external
// file is needed at runtime.
package authz

import (
	"fmt"
	"log/slog"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Resources and actions enforced by middleware.RequireOrgAccess.
const (
	ResourceContacts = "contacts"
	ResourceMembers  = "members"
	ResourceOrgs     = "orgs"

	ActionRead   = "read"
	ActionWrite  = "write"
	ActionDelete = "delete"
)

// Enforcer is the shared instance, set by Init. Casbin enforcers are safe
// for concurrent Enforce calls.
var Enforcer *casbin.Enforcer

const modelText = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`

// Init builds the enforcer on the shared GORM connection. The adapter
// creates casbin_rule if missing (migrations/000009 also declares it for
// migrate-driven environments).
func Init(db *gorm.DB) error {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return fmt.Errorf("authz: adapter: %w", err)
	}
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return fmt.Errorf("authz: model: %w", err)
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return fmt.Errorf("authz: enforcer: %w", err)
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("authz: load policy: %w", err)
	}
	Enforcer = enforcer
	slog.Info("Casbin enforcer ready")
	return nil
}

// Enforce reports whether sub may perform act on obj within domain dom.
// Errors fail closed (deny) and are logged by the caller.
func Enforce(sub, dom, obj, act string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("authz: enforcer not initialized")
	}
	return Enforcer.Enforce(sub, dom, obj, act)
}

// SeedOrgPolicies installs the role policies for a new organization.
// owner: everything; admin: everything but org delete; member: contacts
// read+write; viewer: contacts read.
//
// NOTE: keyMatch is prefix matching, not regex — wildcard resources must
// be enumerated explicitly (".*" never matches under keyMatch).
func SeedOrgPolicies(orgID string) error {
	if Enforcer == nil {
		return fmt.Errorf("authz: enforcer not initialized")
	}
	policies := [][]string{
		{"owner", orgID, ResourceContacts, "(read|write|delete)"},
		{"owner", orgID, ResourceMembers, "(read|write|delete)"},
		{"owner", orgID, ResourceOrgs, "(read|write|delete)"},
		{"admin", orgID, ResourceContacts, "(read|write)"},
		{"admin", orgID, ResourceMembers, "(read|write)"},
		{"admin", orgID, ResourceOrgs, ActionRead},
		{"member", orgID, ResourceContacts, "(read|write)"},
		{"viewer", orgID, ResourceContacts, ActionRead},
	}
	if _, err := Enforcer.AddPolicies(policies); err != nil {
		return fmt.Errorf("authz: seed policies: %w", err)
	}
	return nil
}

// RemoveOrgPolicies drops every policy and grouping scoped to the org.
func RemoveOrgPolicies(orgID string) error {
	if Enforcer == nil {
		return nil
	}
	if _, err := Enforcer.RemoveFilteredPolicy(1, orgID); err != nil {
		return fmt.Errorf("authz: remove policies: %w", err)
	}
	if _, err := Enforcer.RemoveFilteredGroupingPolicy(1, orgID); err != nil {
		return fmt.Errorf("authz: remove groupings: %w", err)
	}
	return nil
}

// GrantRole assigns role to userID within orgID.
func GrantRole(userID, orgID, role string) error {
	if Enforcer == nil {
		return fmt.Errorf("authz: enforcer not initialized")
	}
	if _, err := Enforcer.AddGroupingPolicy(userID, role, orgID); err != nil {
		return fmt.Errorf("authz: grant role: %w", err)
	}
	return nil
}

// RevokeRole removes all of the user's groupings in the org,
// regardless of current role.
func RevokeRole(userID, orgID string) error {
	if Enforcer == nil {
		return fmt.Errorf("authz: enforcer not initialized")
	}
	if _, err := Enforcer.RemoveFilteredGroupingPolicy(0, userID, orgID); err != nil {
		return fmt.Errorf("authz: revoke role: %w", err)
	}
	return nil
}
