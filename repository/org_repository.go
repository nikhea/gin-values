package repository

import (
	"grip/config"
	"grip/models"
)

func CreateOrganization(org *models.Organization) error {
	return config.DB.Create(org).Error
}

func GetOrganizationByID(id string) (*models.Organization, error) {
	var org models.Organization
	err := config.DB.First(&org, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func DeleteOrganization(id string) error {
	return config.DB.Delete(&models.Organization{}, "id = ?", id).Error
}

// ListOrganizationsForUser returns orgs the user belongs to, newest first.
func ListOrganizationsForUser(userID string) ([]models.Organization, error) {
	var orgs []models.Organization
	err := config.DB.Joins("JOIN memberships ON memberships.org_id = organizations.id").
		Where("memberships.user_id = ?", userID).
		Order("organizations.created_at DESC").
		Find(&orgs).Error
	return orgs, err
}

func CreateMembership(m *models.Membership) error {
	return config.DB.Create(m).Error
}

func GetMembership(userID, orgID string) (*models.Membership, error) {
	var m models.Membership
	err := config.DB.First(&m, "user_id = ? AND org_id = ?", userID, orgID).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// UpsertMembership inserts or updates the role. Callers validate first.
func UpsertMembership(m *models.Membership) error {
	return config.DB.Save(m).Error
}

func DeleteMembership(userID, orgID string) error {
	return config.DB.Delete(&models.Membership{}, "user_id = ? AND org_id = ?", userID, orgID).Error
}

// CountOwnersByRole counts memberships holding a role (used to protect
// the last owner from demotion/removal).
func CountOwnersByRole(orgID, role string) (int64, error) {
	var n int64
	err := config.DB.Model(&models.Membership{}).
		Where("org_id = ? AND role = ?", orgID, role).
		Count(&n).Error
	return n, err
}

// ListUserOrgIDs returns every org ID the user belongs to.
func ListUserOrgIDs(userID string) ([]string, error) {
	var ids []string
	err := config.DB.Model(&models.Membership{}).
		Where("user_id = ?", userID).
		Pluck("org_id", &ids).Error
	return ids, err
}
