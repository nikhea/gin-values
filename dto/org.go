package dto

import "grip/models"

// CreateOrgRequest creates an organization; the creator becomes owner.
type CreateOrgRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100" example:"Acme"`
}

// AddMemberRequest adds (or re-roles) a member of an organization.
type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required" example:"458622d8-daba-4252-8ce1-846277353139"`

	Role string `json:"role" binding:"required,oneof=owner admin member viewer" example:"member"`
}

// UpdateMemberRequest changes a member's role.
type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=owner admin member viewer" example:"admin"`
}

type OrgEnvelope struct {
	Message string              `json:"message"`
	Org     models.Organization `json:"org"`
}

type OrgsEnvelope struct {
	Message string                `json:"message"`
	Orgs    []models.Organization `json:"orgs"`
}

type MemberEnvelope struct {
	Message string            `json:"message"`
	Member  models.Membership `json:"member"`
}
