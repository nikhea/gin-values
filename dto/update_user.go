package dto


type UpdateUserRequest struct {

	Name string `json:"name" binding:"required,min=3,max=50"`

	Email string `json:"email" binding:"required,email"`

	Age int `json:"age" binding:"required,min=18,max=100"`
}