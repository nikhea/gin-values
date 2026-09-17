package models

import "time"

type User struct {
	ID string `json:"id"`

	Name string `json:"name" binding:"required,min=3,max=50"`

	Email string `json:"email" binding:"required,email"`

	Age int `json:"age" binding:"required,min=18,max=100"`

    CreatedAt time.Time `json:"created_at"`

    UpdatedAt time.Time `json:"updated_at"`
}