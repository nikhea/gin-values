package repository

import (
	"gin-learn/config"
	"gin-learn/models"
)

func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Find(&users).Error
	return users, err
}

func GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, "id = ?", id).Error
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
