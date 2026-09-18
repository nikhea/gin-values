package services

import (
	"grip/dto"
	"grip/models"
	"grip/repository"

	"github.com/google/uuid"
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
