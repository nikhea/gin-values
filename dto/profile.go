package dto

type CreateProfileRequest struct {
	Bio string `json:"bio" binding:"max=1000"`

	AvatarURL string `json:"avatar_url" binding:"omitempty,url,max=2048"`

	Phone string `json:"phone" binding:"omitempty,max=32"`

	Address string `json:"address" binding:"omitempty,max=500"`
}

type UpdateProfileRequest struct {
	Bio string `json:"bio" binding:"max=1000"`

	AvatarURL string `json:"avatar_url" binding:"omitempty,url,max=2048"`

	Phone string `json:"phone" binding:"omitempty,max=32"`

	Address string `json:"address" binding:"omitempty,max=500"`
}
