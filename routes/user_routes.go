package routes

import (
	"grip/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes mounts CRUD endpoints plus avatar upload on the /api/users group.
func RegisterUserRoutes(api *gin.RouterGroup) {
	api.POST("/", handlers.CreateUser)
	api.GET("/", handlers.GetUsers)
	api.GET("/:id", handlers.GetUser)
	api.PUT("/:id", handlers.UpdateUser)
	api.DELETE("/:id", handlers.DeleteUser)
	api.POST("/:id/avatar", handlers.UploadUserAvatar)
}
