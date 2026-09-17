package routes

import (
	"gin-learn/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes mounts CRUD endpoints on the /api/users group.
func RegisterUserRoutes(api *gin.RouterGroup) {
	api.POST("/", handlers.CreateUser)
	api.GET("/", handlers.GetUsers)
	api.GET("/:id", handlers.GetUser)
	api.PUT("/:id", handlers.UpdateUser)
	api.DELETE("/:id", handlers.DeleteUser)
}
