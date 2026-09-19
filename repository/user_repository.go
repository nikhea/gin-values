package repository

import (
	"grip/config"
	"grip/dto"
	"grip/models"
)

func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Preload("Profile").Find(&users).Error
	return users, err
}

// ListUsers returns a filtered page of users plus the total matching row
// count for pagination metadata.
func ListUsers(filter dto.UserFilter) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	q := config.DB.Model(&models.User{})

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("name ILIKE ? OR email ILIKE ?", like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Preload("Profile").Order("created_at DESC").
		Limit(filter.PageSize).
		Offset(filter.Offset()).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := config.DB.Preload("Profile").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Preload("Profile").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByVerificationToken(token string) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, "verification_token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByResetToken(token string) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, "reset_token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

func DeleteUser(id string) error {
	// Soft delete: sets deleted_at (models.User.DeletedAt).
	return config.DB.Delete(&models.User{}, "id = ?", id).Error
}

// HardDeleteUser permanently removes the row, bypassing soft delete.
// Reserved for admin purges; normal flows use DeleteUser.
func HardDeleteUser(id string) error {
	return config.DB.Unscoped().Delete(&models.User{}, "id = ?", id).Error
}
