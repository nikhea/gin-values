package models

import "time"

// RefreshToken is an opaque, single-use session credential. Only the
// SHA-256 hash is stored — a database leak can never mint sessions.
// Rotation links tokens via ReplacedBy; reuse of a rotated token signals
// theft and revokes the whole family (see service/refresh_service.go).
type RefreshToken struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"-"`

	UserID string `gorm:"index;not null;type:varchar(36)" json:"-"`

	TokenHash string `gorm:"uniqueIndex;not null" json:"-"`

	ExpiresAt time.Time `json:"-"`

	RevokedAt *time.Time `json:"-"`

	ReplacedBy *string `gorm:"type:varchar(36)" json:"-"`

	CreatedAt time.Time `json:"-"`
}
