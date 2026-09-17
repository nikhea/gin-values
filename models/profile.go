package models

import "time"

type Profile struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`

	UserID string `gorm:"uniqueIndex;not null;type:varchar(36)" json:"user_id"`

	Bio string `json:"bio"`

	AvatarURL string `json:"avatar_url"`

	Phone string `json:"phone"`

	Address string `json:"address"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
