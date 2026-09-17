package routes

import (
	"gin-learn/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterContactRoutes mounts contact endpoints on their own
// top-level group: /api/contacts.
func RegisterContactRoutes(router *gin.Engine) {
	contacts := router.Group("/api/contacts")
	{
		contacts.POST("/", handlers.CreateContact)
		contacts.GET("/", handlers.GetContacts)
		contacts.GET("/:id", handlers.GetContact)
		contacts.PUT("/:id", handlers.UpdateContact)
		contacts.DELETE("/:id", handlers.DeleteContact)
	}
}
