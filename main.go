package main

import (
	"gin-learn/config"
	"gin-learn/routes"
)

// @title Gin Learn API
// @version 1.0
// @description Users, profiles and contacts API.
// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	config.LoadEnv()
	config.RunMigrations()
	config.ConnectDatabase()

	router := routes.Setup()

	router.Run(":8080")
}
