package services

import (
	"errors"

	"grip/dto"
	"grip/models"
	"grip/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateContact(req dto.CreateContactRequest) (*models.Contact, error) {
	if err := dto.ValidateContactValue(req.Type, req.Value); err != nil {
		return nil, err
	}

	// Ensure the user exists so the FK holds and callers get 404 otherwise.
	if _, err := repository.GetUserByID(req.UserID); err != nil {
		return nil, err
	}

	contact := &models.Contact{
		ID:     uuid.New().String(),
		UserID: req.UserID,
		Name:   req.Name,
		Type:   req.Type,
		Value:  req.Value,
	}

	if err := repository.CreateContact(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

func ListContacts(filter dto.ContactFilter) ([]models.Contact, int64, error) {
	filter.Normalize()
	return repository.ListContacts(filter)
}

func GetContactByID(id string) (*models.Contact, error) {
	return repository.GetContactByID(id)
}

func UpdateContact(id string, req dto.UpdateContactRequest) (*models.Contact, error) {
	if err := dto.ValidateContactValue(req.Type, req.Value); err != nil {
		return nil, err
	}

	contact, err := repository.GetContactByID(id)
	if err != nil {
		return nil, err
	}

	contact.Name = req.Name
	contact.Type = req.Type
	contact.Value = req.Value

	if err := repository.UpdateContact(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

func DeleteContact(id string) error {
	if _, err := repository.GetContactByID(id); err != nil {
		return err
	}

	return repository.DeleteContact(id)
}

// CreateContactForUser is the HTTP-facing create: the contact is always
// owned by the caller (req.UserID is overwritten), and team contacts
// require a writer role (owner/admin/member) in the org.
func CreateContactForUser(actorID string, req dto.CreateContactRequest) (*models.Contact, error) {
	if err := dto.ValidateContactValue(req.Type, req.Value); err != nil {
		return nil, err
	}

	if _, err := repository.GetUserByID(actorID); err != nil {
		return nil, err
	}

	var orgID *string
	if req.OrgID != nil && *req.OrgID != "" {
		if err := requireOrgWriter(actorID, *req.OrgID); err != nil {
			return nil, err
		}
		orgID = req.OrgID
	}

	contact := &models.Contact{
		ID:     uuid.New().String(),
		UserID: actorID,
		OrgID:  orgID,
		Name:   req.Name,
		Type:   req.Type,
		Value:  req.Value,
	}

	if err := repository.CreateContact(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// ListContactsForUser lists what the caller may see: with an explicit org
// scope (member-only, any role) or otherwise their personal contacts plus
// contacts of orgs they belong to.
func ListContactsForUser(actorID string, filter dto.ContactFilter) ([]models.Contact, int64, error) {
	filter.Normalize()

	if filter.OrgID != "" {
		if _, err := repository.GetMembership(actorID, filter.OrgID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 0, ErrForbidden
			}
			return nil, 0, err
		}
		return repository.ListContacts(filter)
	}

	orgIDs, err := repository.ListUserOrgIDs(actorID)
	if err != nil {
		return nil, 0, err
	}
	return repository.ListVisibleContacts(actorID, orgIDs, filter)
}

// GetContactForUser loads a contact the caller may read: own contacts,
// plus team contacts for org members (any role). Anything else is 404.
func GetContactForUser(actorID, id string) (*models.Contact, error) {
	contact, err := repository.GetContactByID(id)
	if err != nil {
		return nil, err
	}
	return contact, authorizeContactRead(actorID, contact)
}

// UpdateContactForUser overwrites a contact the caller may write: own
// contacts, plus team contacts for owner/admin/member. Viewers get 403,
// strangers get 404.
func UpdateContactForUser(actorID, id string, req dto.UpdateContactRequest) (*models.Contact, error) {
	if err := dto.ValidateContactValue(req.Type, req.Value); err != nil {
		return nil, err
	}

	contact, err := repository.GetContactByID(id)
	if err != nil {
		return nil, err
	}
	if err := authorizeContactWrite(actorID, contact); err != nil {
		return nil, err
	}

	contact.Name = req.Name
	contact.Type = req.Type
	contact.Value = req.Value

	if err := repository.UpdateContact(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// DeleteContactForUser soft-deletes a contact the caller may write.
func DeleteContactForUser(actorID, id string) error {
	contact, err := repository.GetContactByID(id)
	if err != nil {
		return err
	}
	if err := authorizeContactWrite(actorID, contact); err != nil {
		return err
	}

	return repository.DeleteContact(id)
}

// authorizeContactRead allows owners and org members (any role).
func authorizeContactRead(actorID string, contact *models.Contact) error {
	if contact.UserID == actorID {
		return nil
	}
	if contact.OrgID == nil {
		return gorm.ErrRecordNotFound
	}
	if _, err := repository.GetMembership(actorID, *contact.OrgID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// authorizeContactWrite allows owners and writer roles; viewers get 403.
func authorizeContactWrite(actorID string, contact *models.Contact) error {
	if contact.UserID == actorID {
		return nil
	}
	if contact.OrgID == nil {
		return gorm.ErrRecordNotFound
	}
	return requireOrgWriter(actorID, *contact.OrgID)
}

// requireOrgWriter enforces owner/admin/member membership. Non-members
// surface as 404 (no existence leak); viewers as 403.
func requireOrgWriter(actorID, orgID string) error {
	if _, err := repository.GetOrganizationByID(orgID); err != nil {
		return err
	}
	m, err := repository.GetMembership(actorID, orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	switch m.Role {
	case models.RoleOwner, models.RoleAdmin, models.RoleMember:
		return nil
	default:
		return ErrForbidden
	}
}
