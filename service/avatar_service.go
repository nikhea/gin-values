package services

import (
	"errors"

	"gin-learn/models"
	"gin-learn/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SetUserAvatar stores avatarURL on the user's profile, creating the
// profile when the user has none. Other profile fields are untouched.
func SetUserAvatar(userID, avatarURL string) (*models.Profile, error) {
	if _, err := repository.GetUserByID(userID); err != nil {
		return nil, err
	}

	profile, err := repository.GetProfileByUserID(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		profile = &models.Profile{ID: uuid.New().String(), UserID: userID, AvatarURL: avatarURL}
		if err := repository.CreateProfile(profile); err != nil {
			return nil, err
		}
		return profile, nil
	}

	profile.AvatarURL = avatarURL
	if err := repository.UpdateProfile(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// SetContactAvatar stores avatarURL on the contact.
func SetContactAvatar(contactID, avatarURL string) (*models.Contact, error) {
	contact, err := repository.GetContactByID(contactID)
	if err != nil {
		return nil, err
	}

	contact.AvatarURL = avatarURL
	if err := repository.UpdateContact(contact); err != nil {
		return nil, err
	}
	return contact, nil
}
