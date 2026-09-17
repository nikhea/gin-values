package routes

import (
	"gin-learn/middleware"

	_ "gin-learn/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup registers all API routes and returns the configured router.
func Setup() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.RequestID())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	base := "/api"
	api := router.Group(base + "/users")
	{
		RegisterUserRoutes(api)
		RegisterProfileRoutes(api)
	}

	RegisterContactRoutes(router)

	return router
}
