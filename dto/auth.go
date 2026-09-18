package dto

import "grip/models"

type RegisterRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50" example:"Kaige Saif"`

	Email string `json:"email" binding:"required,email" example:"kaige@example.com"`

	Password string `json:"password" binding:"required,min=8,max=72" example:"supersecret123"`

	Age int `json:"age" binding:"required,min=18,max=100" example:"30"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email" example:"kaige@example.com"`

	Password string `json:"password" binding:"required,min=1,max=72" example:"supersecret123"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"kaige@example.com"`
}

type ResetPasswordRequest struct {
	// Token may alternatively be supplied as ?token= query param.
	Token string `json:"token" example:"9f2c4a1e6b4d4f8a9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b"`

	NewPassword string `json:"new_password" binding:"required,min=8,max=72" example:"newsecret123"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email" example:"kaige@example.com"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"kaige@example.com"`

	Code string `json:"code" binding:"required,len=6,numeric" example:"482914"`
}

type AuthResponse struct {
	Message string      `json:"message" example:"Logged in"`
	Token   string      `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3OCJ9.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlF"`
	User    models.User `json:"user"`
}
