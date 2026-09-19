package services

import (
	"errors"

	"grip/authz"
	"grip/dto"
	"grip/models"
	"grip/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidRole   = errors.New("invalid role (owner, admin, member, viewer)")
	ErrAlreadyMember = errors.New("user is already a member")
	ErrLastOwner     = errors.New("cannot remove the last owner")
)

// CreateOrg creates an organization, makes the creator its owner, and
// seeds + grants the Casbin policies.
func CreateOrg(creatorID string, req dto.CreateOrgRequest) (*models.Organization, error) {
	if _, err := repository.GetUserByID(creatorID); err != nil {
		return nil, err
	}

	org := &models.Organization{ID: uuid.New().String(), Name: req.Name}
	if err := repository.CreateOrganization(org); err != nil {
		return nil, err
	}
	if err := repository.CreateMembership(&models.Membership{
		UserID: creatorID, OrgID: org.ID, Role: models.RoleOwner,
	}); err != nil {
		_ = repository.DeleteOrganization(org.ID)
		return nil, err
	}
	if err := authz.SeedOrgPolicies(org.ID); err != nil {
		_ = repository.DeleteMembership(creatorID, org.ID)
		_ = repository.DeleteOrganization(org.ID)
		return nil, err
	}
	if err := authz.GrantRole(creatorID, org.ID, models.RoleOwner); err != nil {
		_ = repository.DeleteMembership(creatorID, org.ID)
		_ = repository.DeleteOrganization(org.ID)
		return nil, err
	}
	return org, nil
}

// ListOrgs returns organizations the user belongs to.
func ListOrgs(userID string) ([]models.Organization, error) {
	return repository.ListOrganizationsForUser(userID)
}

// AddMember adds a user to the org (or updates their role if present).
// Route middleware guarantees the caller may manage members.
func AddMember(orgID, targetUserID, role string) (*models.Membership, error) {
	if !models.ValidRole(role) {
		return nil, ErrInvalidRole
	}
	if _, err := repository.GetOrganizationByID(orgID); err != nil {
		return nil, err
	}
	if _, err := repository.GetUserByID(targetUserID); err != nil {
		return nil, err
	}

	if _, err := repository.GetMembership(targetUserID, orgID); err == nil {
		return nil, ErrAlreadyMember
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	m := &models.Membership{UserID: targetUserID, OrgID: orgID, Role: role}
	if err := repository.CreateMembership(m); err != nil {
		return nil, err
	}
	if err := authz.GrantRole(targetUserID, orgID, role); err != nil {
		_ = repository.DeleteMembership(targetUserID, orgID)
		return nil, err
	}
	return m, nil
}

// UpdateMemberRole changes a member's role. Demoting the last owner fails.
func UpdateMemberRole(orgID, targetUserID, role string) (*models.Membership, error) {
	if !models.ValidRole(role) {
		return nil, ErrInvalidRole
	}
	m, err := repository.GetMembership(targetUserID, orgID)
	if err != nil {
		return nil, err
	}
	if m.Role == models.RoleOwner && role != models.RoleOwner {
		n, err := repository.CountOwnersByRole(orgID, models.RoleOwner)
		if err != nil {
			return nil, err
		}
		if n <= 1 {
			return nil, ErrLastOwner
		}
	}

	m.Role = role
	if err := repository.UpsertMembership(m); err != nil {
		return nil, err
	}
	if err := authz.RevokeRole(targetUserID, orgID); err != nil {
		return nil, err
	}
	if err := authz.GrantRole(targetUserID, orgID, role); err != nil {
		return nil, err
	}
	return m, nil
}

// RemoveMember kicks a user out. Removing the last owner fails.
func RemoveMember(orgID, targetUserID string) error {
	m, err := repository.GetMembership(targetUserID, orgID)
	if err != nil {
		return err
	}
	if m.Role == models.RoleOwner {
		n, err := repository.CountOwnersByRole(orgID, models.RoleOwner)
		if err != nil {
			return err
		}
		if n <= 1 {
			return ErrLastOwner
		}
	}
	if err := repository.DeleteMembership(targetUserID, orgID); err != nil {
		return err
	}
	return authz.RevokeRole(targetUserID, orgID)
}

// DeleteOrg removes the organization (memberships cascade, team contacts
// detach to personal via ON DELETE SET NULL) and its Casbin policies.
func DeleteOrg(orgID string) error {
	if _, err := repository.GetOrganizationByID(orgID); err != nil {
		return err
	}
	if err := repository.DeleteOrganization(orgID); err != nil {
		return err
	}
	return authz.RemoveOrgPolicies(orgID)
}
