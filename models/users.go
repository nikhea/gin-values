package models

import "time"

type User struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id" example:"458622d8-daba-4252-8ce1-846277353139"`

	Name string `json:"name" example:"Kaige Saif"`

	Email string `gorm:"uniqueIndex" json:"email" example:"kaige@example.com"`

	Age int `json:"age" example:"30"`

	// PasswordHash never leaves the server: json:"-" keeps it out of API responses.
	PasswordHash string `gorm:"column:password_hash" json:"-"`

	EmailVerified bool `gorm:"column:email_verified" json:"email_verified" example:"true"`

	// Email verification / password-reset tokens. Never serialized.
	VerificationToken string    `gorm:"column:verification_token" json:"-"`
	ResetToken        string    `gorm:"column:reset_token" json:"-"`
	ResetExpiresAt    time.Time `gorm:"column:reset_expires_at" json:"-"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	Profile *Profile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
}
