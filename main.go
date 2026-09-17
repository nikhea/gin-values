package main

import (
	"gin-learn/config"
	"gin-learn/routes"
)

func main() {
	config.LoadEnv()
	config.RunMigrations()
	config.ConnectDatabase()

	router := routes.Setup()

	router.Run(":8080")
}
