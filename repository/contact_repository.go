package repository

import (
	"grip/config"
	"grip/dto"
	"grip/models"
)

func CreateContact(contact *models.Contact) error {
	return config.DB.Create(contact).Error
}

func GetContactByID(id string) (*models.Contact, error) {
	var contact models.Contact
	err := config.DB.First(&contact, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// ListContacts returns a filtered page of contacts plus the total
// matching row count for pagination metadata.
func ListContacts(filter dto.ContactFilter) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	q := config.DB.Model(&models.Contact{})

	if filter.UserID != "" {
		q = q.Where("user_id = ?", filter.UserID)
	}
	if filter.Type != "" {
		q = q.Where("type = ?", filter.Type)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("name ILIKE ? OR value ILIKE ?", like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order("created_at DESC").
		Limit(filter.PageSize).
		Offset(filter.Offset()).
		Find(&contacts).Error
	if err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func UpdateContact(contact *models.Contact) error {
	return config.DB.Save(contact).Error
}

func DeleteContact(id string) error {
	return config.DB.Delete(&models.Contact{}, "id = ?", id).Error
}
