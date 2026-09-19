package models

import "time"

// Organization roles, weakest to strongest.
const (
	RoleViewer = "viewer"
	RoleMember = "member"
	RoleAdmin  = "admin"
	RoleOwner  = "owner"
)

// Organization groups users and contacts for team-based access control.
// Casbin treats the org ID as a policy domain (see authz package).
type Organization struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id" example:"b3c4d5e6-f7a8-49b0-c1d2-e3f4a5b6c7d8"`

	Name string `json:"name" example:"Acme"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

// Membership binds a user to an organization with a role.
type Membership struct {
	UserID string `gorm:"primaryKey;type:varchar(36)" json:"user_id"`

	OrgID string `gorm:"primaryKey;type:varchar(36)" json:"org_id"`

	Role string `json:"role" example:"admin"`

	CreatedAt time.Time `json:"created_at"`
}

// ValidRole reports whether r is a known organization role.
func ValidRole(r string) bool {
	switch r {
	case RoleViewer, RoleMember, RoleAdmin, RoleOwner:
		return true
	}
	return false
}
