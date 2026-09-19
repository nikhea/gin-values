package handlers

import (
	"errors"
	"net/http"

	"grip/dto"
	"grip/middleware"
	"grip/models"
	services "grip/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	_ = dto.NotificationEnvelope{}
	_ = dto.NotificationsEnvelope{}
	_ = dto.UnreadEnvelope{}
	_ = models.Notification{}
)

// ListNotifications godoc
// @Summary List my notifications (newest first)
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Param unread_only query bool false "Unread only"
// @Param page query int false "Page number" default(1) example(1)
// @Param page_size query int false "Page size" default(10) example(10)
// @Success 200 {object} dto.NotificationsEnvelope
// @Router /notifications/ [get]
func ListNotifications(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	var q dto.NotificationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, total, err := services.ListNotifications(userID, q.UnreadOnly, q.Page, q.PageSize)
	if err != nil {
		internalError(c, err)
		return
	}

	page, pageSize := q.Page, q.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "Notifications found",
		"notifications": rows,
		"meta": dto.PageMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		},
	})
}

// UnreadCount godoc
// @Summary Count my unread notifications
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UnreadEnvelope
// @Router /notifications/unread-count [get]
func UnreadCount(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	n, err := services.UnreadNotificationCount(userID)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unread count", "unread": n})
}

// MarkRead godoc
// @Summary Mark one notification read
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} dto.NotificationEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /notifications/{id}/read [patch]
func MarkRead(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	n, err := services.MarkNotificationRead(userID, c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Notification not found"})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Marked read", "notification": n})
}

// MarkAllRead godoc
// @Summary Mark all my notifications read
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MessageEnvelope
// @Router /notifications/read-all [post]
func MarkAllRead(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	n, err := services.MarkAllNotificationsRead(userID)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications marked read",
		"marked":  n,
	})
}
