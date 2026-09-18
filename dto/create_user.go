package dto

type CreateUserRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50" example:"Kaige Saif"`

	Email string `json:"email" binding:"required,email" example:"kaige@example.com"`

	Age int `json:"age" binding:"required,min=18,max=100" example:"30"`
}
