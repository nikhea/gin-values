package routes

import (
	"grip/handlers"
	"grip/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes mounts public auth endpoints plus the
// authenticated /me endpoint on the /api/auth group.
func RegisterAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	// Strict per-IP budget: these endpoints are anonymous and brute-forceable.
	auth.Use(middleware.RateLimit("auth", middleware.AuthRateLimit))
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/verify", handlers.VerifyEmail)
		auth.POST("/verify-otp", handlers.VerifyOTP)
		auth.POST("/resend-verification", handlers.ResendVerification)
		auth.POST("/forgot-password", handlers.ForgotPassword)
		auth.POST("/reset-password", handlers.ResetPassword)
		auth.GET("/me", middleware.AuthRequired(), handlers.Me)
	}
}
