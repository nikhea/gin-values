package dto

import "grip/models"

// NotificationQuery carries inbox list params.
type NotificationQuery struct {
	UnreadOnly bool `form:"unread_only"`
	Page       int  `form:"page"`
	PageSize   int  `form:"page_size"`
}

type NotificationEnvelope struct {
	Message      string              `json:"message"`
	Notification models.Notification `json:"notification"`
}

type NotificationsEnvelope struct {
	Message       string                `json:"message"`
	Notifications []models.Notification `json:"notifications"`
	Meta          PageMeta              `json:"meta"`
}

type UnreadEnvelope struct {
	Message string `json:"message"`
	Unread  int64  `json:"unread" example:"3"`
}
