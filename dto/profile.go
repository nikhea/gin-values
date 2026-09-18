package dto

type CreateProfileRequest struct {
	Bio string `json:"bio" binding:"max=1000" example:"Gopher since 2021. Coffee first."`

	AvatarURL string `json:"avatar_url" binding:"omitempty,url,max=2048" example:"https://example.com/avatars/kaige.png"`

	Phone string `json:"phone" binding:"omitempty,max=32" example:"+15551234567"`

	Address string `json:"address" binding:"omitempty,max=500" example:"123 Main St, Lagos"`
}

type UpdateProfileRequest struct {
	Bio string `json:"bio" binding:"max=1000" example:"Gopher since 2021. Coffee first."`

	AvatarURL string `json:"avatar_url" binding:"omitempty,url,max=2048" example:"https://example.com/avatars/kaige.png"`

	Phone string `json:"phone" binding:"omitempty,max=32" example:"+15551234567"`

	Address string `json:"address" binding:"omitempty,max=500" example:"123 Main St, Lagos"`
}
