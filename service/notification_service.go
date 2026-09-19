package services

import (
	"grip/models"
	"grip/repository"
)

// ListNotifications returns the caller's inbox page plus total.
func ListNotifications(userID string, unreadOnly bool, page, pageSize int) ([]models.Notification, int64, error) {
	return repository.ListNotifications(repository.NotificationFilter{
		UserID: userID, UnreadOnly: unreadOnly, Page: page, PageSize: pageSize,
	})
}

// UnreadNotificationCount counts the caller's unread rows.
func UnreadNotificationCount(userID string) (int64, error) {
	return repository.UnreadCount(userID)
}

// MarkNotificationRead marks one owned row read (404 otherwise).
func MarkNotificationRead(userID, id string) (*models.Notification, error) {
	return repository.MarkNotificationRead(userID, id)
}

// MarkAllNotificationsRead marks the whole inbox read, returning the count.
func MarkAllNotificationsRead(userID string) (int64, error) {
	return repository.MarkAllNotificationsRead(userID)
}
