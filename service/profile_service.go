package services

import (
	"gin-learn/dto"
	"gin-learn/models"
	"gin-learn/repository"

	"github.com/google/uuid"
)

func CreateProfile(userID string, req dto.CreateProfileRequest) (*models.Profile, error) {
	// Ensure the user exists so the FK holds and callers get 404 otherwise.
	if _, err := repository.GetUserByID(userID); err != nil {
		return nil, err
	}

	profile := &models.Profile{
		ID:        uuid.New().String(),
		UserID:    userID,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		Phone:     req.Phone,
		Address:   req.Address,
	}

	if err := repository.CreateProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func GetProfileByUserID(userID string) (*models.Profile, error) {
	return repository.GetProfileByUserID(userID)
}

func UpdateProfile(userID string, req dto.UpdateProfileRequest) (*models.Profile, error) {
	profile, err := repository.GetProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	profile.Bio = req.Bio
	profile.AvatarURL = req.AvatarURL
	profile.Phone = req.Phone
	profile.Address = req.Address

	if err := repository.UpdateProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func DeleteProfile(userID string) error {
	if _, err := repository.GetProfileByUserID(userID); err != nil {
		return err
	}

	return repository.DeleteProfileByUserID(userID)
}
