package routes

import (
	"github.com/gin-gonic/gin"
)

// Setup registers all API routes and returns the configured router.
func Setup() *gin.Engine {
	router := gin.Default()
    base := "/api"
	api := router.Group(base + "/users")
	{
		RegisterUserRoutes(api)
		RegisterProfileRoutes(api)
	}

	return router
}
