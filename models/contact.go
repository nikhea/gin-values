package models

import "time"

// Contact types supported by ValidateContactValue.
const (
	ContactTypeEmail   = "email"
	ContactTypePhone   = "phone"
	ContactTypeAddress = "address"
	ContactTypeOther   = "other"
)

type Contact struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id" example:"9c9e6679-7425-40de-944b-e07fc1f90ae7"`

	UserID string `gorm:"index;not null;type:varchar(36)" json:"user_id" example:"458622d8-daba-4252-8ce1-846277353139"`

	Name string `json:"name" example:"Work email"`

	Type string `gorm:"index;type:varchar(16)" json:"type" example:"email"`

	Value string `json:"value" example:"kaige@work.com"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
