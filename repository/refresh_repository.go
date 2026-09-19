package repository

import (
	"time"

	"grip/config"
	"grip/models"
)

// CreateRefreshToken stores a new refresh token row.
func CreateRefreshToken(token *models.RefreshToken) error {
	return config.DB.Create(token).Error
}

// GetRefreshTokenByHash loads a token by its SHA-256 hex digest.
func GetRefreshTokenByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := config.DB.First(&token, "token_hash = ?", hash).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetRefreshTokenByID loads a token row by primary key.
func GetRefreshTokenByID(id string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := config.DB.First(&token, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// UpdateRefreshToken persists rotation/revocation changes.
func UpdateRefreshToken(token *models.RefreshToken) error {
	return config.DB.Save(token).Error
}

// RevokeUserTokens marks every active token of a user revoked (logout-all).
func RevokeUserTokens(userID string) error {
	now := time.Now()
	return config.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// RevokeTokenFamily marks every token in the rotation chain revoked.
// The chain is walked via ReplacedBy links starting from the given row.
func RevokeTokenFamily(start *models.RefreshToken) error {
	seen := map[string]bool{}
	current := start
	for current != nil && !seen[current.ID] {
		seen[current.ID] = true
		now := time.Now()
		current.RevokedAt = &now
		if err := config.DB.Save(current).Error; err != nil {
			return err
		}
		if current.ReplacedBy == nil {
			break
		}
		next := *current.ReplacedBy
		var child models.RefreshToken
		if err := config.DB.First(&child, "id = ?", next).Error; err != nil {
			break
		}
		current = &child
	}
	return nil
}
