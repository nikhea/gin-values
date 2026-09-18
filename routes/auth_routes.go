package routes

import (
	"gin-learn/handlers"
	"gin-learn/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes mounts public auth endpoints plus the
// authenticated /me endpoint on the /api/auth group.
func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/verify", handlers.VerifyEmail)
		auth.POST("/resend-verification", handlers.ResendVerification)
		auth.POST("/forgot-password", handlers.ForgotPassword)
		auth.POST("/reset-password", handlers.ResetPassword)
		auth.GET("/me", middleware.AuthRequired(), handlers.Me)
	}
}
