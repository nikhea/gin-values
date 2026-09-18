package repository

import (
	"grip/config"
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
	return config.DB.Delete(&models.User{}, "id = ?", id).Error
}
