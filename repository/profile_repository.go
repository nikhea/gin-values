package repository

import (
	"gin-learn/config"
	"gin-learn/models"
)

func CreateProfile(profile *models.Profile) error {
	return config.DB.Create(profile).Error
}

func GetProfileByUserID(userID string) (*models.Profile, error) {
	var profile models.Profile
	err := config.DB.First(&profile, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func UpdateProfile(profile *models.Profile) error {
	return config.DB.Save(profile).Error
}

func DeleteProfileByUserID(userID string) error {
	return config.DB.Delete(&models.Profile{}, "user_id = ?", userID).Error
}
