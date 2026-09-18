package routes

import (
	"gin-learn/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterContactRoutes mounts contact endpoints on the protected
// group: /api/contacts (requires JWT via AuthRequired).
func RegisterContactRoutes(router *gin.RouterGroup) {
	contacts := router.Group("/contacts")
	{
		contacts.POST("/", handlers.CreateContact)
		contacts.GET("/", handlers.GetContacts)
		contacts.GET("/:id", handlers.GetContact)
		contacts.PUT("/:id", handlers.UpdateContact)
		contacts.DELETE("/:id", handlers.DeleteContact)
	}
}
