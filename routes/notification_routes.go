package routes

import (
	"grip/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterNotificationRoutes mounts the personal inbox on the protected
// group. All endpoints are self-scoped by the caller's JWT — users can
// only ever see their own rows.
func RegisterNotificationRoutes(api *gin.RouterGroup) {
	n := api.Group("/notifications")
	{
		n.GET("/", handlers.ListNotifications)
		n.GET("/unread-count", handlers.UnreadCount)
		n.PATCH("/:id/read", handlers.MarkRead)
		n.POST("/read-all", handlers.MarkAllRead)
	}
}
