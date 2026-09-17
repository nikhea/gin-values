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
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`

	UserID string `gorm:"index;not null;type:varchar(36)" json:"user_id"`

	Name string `json:"name"`

	Type string `gorm:"index;type:varchar(16)" json:"type"`

	Value string `json:"value"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
