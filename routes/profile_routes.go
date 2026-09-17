package routes

import (
	"gin-learn/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterProfileRoutes mounts one-to-one profile endpoints nested
// under the /api/users group: /api/users/:id/profile.
func RegisterProfileRoutes(api *gin.RouterGroup) {
	api.POST("/:id/profile", handlers.CreateProfile)
	api.GET("/:id/profile", handlers.GetProfile)
	api.PUT("/:id/profile", handlers.UpdateProfile)
	api.DELETE("/:id/profile", handlers.DeleteProfile)
}
