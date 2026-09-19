package models

import (
	"time"

	"gorm.io/gorm"
)

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

	// OrgID is nil for personal contacts (owner-only). Set for team contacts.
	OrgID *string `gorm:"index;type:varchar(36)" json:"org_id,omitempty" example:"b3c4d5e6-f7a8-49b0-c1d2-e3f4a5b6c7d8"`

	Name string `json:"name" example:"Work email"`

	Type string `gorm:"index;type:varchar(16)" json:"type" example:"email"`

	Value string `json:"value" example:"kaige@work.com"`

	AvatarURL string `gorm:"column:avatar_url" json:"avatar_url,omitempty" example:"/uploads/avatars/contacts/9c9e6679-7425-40de-944b-e07fc1f90ae7.png"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"primitive,string"`
}
