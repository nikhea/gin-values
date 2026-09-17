package models

import "time"

type User struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`

	Name string `json:"name"`

	Email string `gorm:"uniqueIndex" json:"email"`

	Age int `json:"age"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	Profile *Profile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
}
