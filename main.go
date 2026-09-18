package main

import (
	"gin-learn/config"
	"gin-learn/routes"
)

// @title Gin Learn API
// @version 1.0
// @description Users, profiles and contacts API with JWT auth and email verification.
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT. Example: "Bearer eyJhbGciOiJIUzI1NiJ9..."
func main() {
	config.LoadEnv()
	config.RunMigrations()
	config.ConnectDatabase()

	router := routes.Setup()

	router.Run(":8080")
}
