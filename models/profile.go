package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id" example:"7c9e6679-7425-40de-944b-e07fc1f90ae7"`

	UserID string `gorm:"uniqueIndex;not null;type:varchar(36)" json:"user_id" example:"458622d8-daba-4252-8ce1-846277353139"`

	Bio string `json:"bio" example:"Gopher since 2021. Coffee first."`

	AvatarURL string `json:"avatar_url" example:"https://example.com/avatars/kaige.png"`

	Phone string `json:"phone" example:"+15551234567"`

	Address string `json:"address" example:"123 Main St, Lagos"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"primitive,string"`
}
