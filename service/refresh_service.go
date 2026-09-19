package services

import (
	"errors"
	"time"

	"grip/config"
	"grip/models"
	"grip/repository"
	"grip/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrRefreshReuse        = errors.New("refresh token reuse detected, all sessions revoked")
)

// IssueTokenPair creates a fresh access JWT plus a single-use opaque
// refresh token. Only the refresh hash is stored.
func IssueTokenPair(user *models.User) (access, refresh string, err error) {
	access, err = utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}

	raw, err := utils.GenerateSecureToken(32)
	if err != nil {
		return "", "", err
	}
	row := &models.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: utils.SHA256Hex(raw),
		ExpiresAt: time.Now().Add(config.RefreshTTL()),
	}
	if err := repository.CreateRefreshToken(row); err != nil {
		return "", "", err
	}
	return access, raw, nil
}

// Rotate consumes a refresh token and returns a new pair. The presented
// token is revoked and linked to its replacement. Presenting an already
// rotated (or revoked) token signals theft: the whole family is revoked
// and every session must re-authenticate.
func Rotate(raw string) (access, refresh string, user *models.User, err error) {
	if raw == "" {
		return "", "", nil, ErrInvalidRefreshToken
	}
	row, err := repository.GetRefreshTokenByHash(utils.SHA256Hex(raw))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, ErrInvalidRefreshToken
		}
		return "", "", nil, err
	}
	if row.RevokedAt != nil {
		// A rotated token presented again signals theft — but only while
		// its replacement is still alive. Otherwise (logout, logout-all,
		// or an already-killed family) it is just invalid.
		if row.ReplacedBy != nil {
			if next, err := repository.GetRefreshTokenByID(*row.ReplacedBy); err == nil &&
				next.RevokedAt == nil && time.Now().Before(next.ExpiresAt) {
				_ = repository.RevokeTokenFamily(row)
				return "", "", nil, ErrRefreshReuse
			}
		}
		return "", "", nil, ErrInvalidRefreshToken
	}
	if time.Now().After(row.ExpiresAt) {
		return "", "", nil, ErrInvalidRefreshToken
	}

	user, err = repository.GetUserByID(row.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, ErrInvalidRefreshToken
		}
		return "", "", nil, err
	}

	access, refresh, err = IssueTokenPair(user)
	if err != nil {
		return "", "", nil, err
	}
	newHash := utils.SHA256Hex(refresh)
	newRow, err := repository.GetRefreshTokenByHash(newHash)
	if err != nil {
		return "", "", nil, err
	}
	now := time.Now()
	row.RevokedAt = &now
	row.ReplacedBy = &newRow.ID
	if err := repository.UpdateRefreshToken(row); err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

// Revoke invalidates one refresh token (logout). Unknown tokens succeed
// silently so logout never leaks token validity.
func Revoke(raw string) error {
	if raw == "" {
		return nil
	}
	row, err := repository.GetRefreshTokenByHash(utils.SHA256Hex(raw))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if row.RevokedAt != nil {
		return nil
	}
	now := time.Now()
	row.RevokedAt = &now
	return repository.UpdateRefreshToken(row)
}

// RevokeAll invalidates every active refresh token of a user (logout-all).
func RevokeAll(userID string) error {
	return repository.RevokeUserTokens(userID)
}
