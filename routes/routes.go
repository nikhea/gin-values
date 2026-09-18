package routes

import (
	"gin-learn/handlers"
	"gin-learn/middleware"

	_ "gin-learn/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup registers all API routes and returns the configured router.
//
// /api/auth/* is public (register, login, verify, password reset).
// Everything else requires a JWT via the Authorization: Bearer header.
func Setup() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.RequestID())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Liveness probe (public, no auth): used by Docker healthchecks.
	router.GET("/health", handlers.Health)

	// Publicly served uploaded files (avatars): /uploads/avatars/...
	router.Static("/uploads", "./uploads")

	RegisterAuthRoutes(router)

	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		users := protected.Group("/users")
		{
			RegisterUserRoutes(users)
			RegisterProfileRoutes(users)
		}

		RegisterContactRoutes(protected)
	}

	return router
}
