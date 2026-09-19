package repository

import (
	"time"

	"grip/config"
	"grip/models"
)

// CreateNotification inserts one inbox row.
func CreateNotification(n *models.Notification) error {
	return config.DB.Create(n).Error
}

// NotificationFilter scopes inbox listing. UnreadOnly shows read_at IS NULL.
type NotificationFilter struct {
	UserID     string
	UnreadOnly bool
	Page       int
	PageSize   int
}

// Normalize applies pagination defaults (page 1, 10 per page, max 100).
func (f *NotificationFilter) Normalize() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 10
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
}

// Offset returns the SQL offset for the current page.
func (f NotificationFilter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// ListNotifications returns newest-first rows plus total.
func ListNotifications(filter NotificationFilter) ([]models.Notification, int64, error) {
	filter.Normalize()

	var rows []models.Notification
	var total int64

	q := config.DB.Model(&models.Notification{}).Where("user_id = ?", filter.UserID)
	if filter.UnreadOnly {
		q = q.Where("read_at IS NULL")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Limit(filter.PageSize).Offset(filter.Offset()).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// UnreadCount counts unread inbox rows for a user.
func UnreadCount(userID string) (int64, error) {
	var n int64
	err := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&n).Error
	return n, err
}

// GetNotificationByID loads one row scoped to its owner.
func GetNotificationByID(userID, id string) (*models.Notification, error) {
	var n models.Notification
	err := config.DB.First(&n, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// MarkNotificationRead stamps read_at (idempotent).
func MarkNotificationRead(userID, id string) (*models.Notification, error) {
	n, err := GetNotificationByID(userID, id)
	if err != nil {
		return nil, err
	}
	if n.ReadAt == nil {
		now := time.Now()
		n.ReadAt = &now
		if err := config.DB.Save(n).Error; err != nil {
			return nil, err
		}
	}
	return n, nil
}

// MarkAllNotificationsRead stamps every unread row; returns the count marked.
func MarkAllNotificationsRead(userID string) (int64, error) {
	res := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", time.Now())
	return res.RowsAffected, res.Error
}
