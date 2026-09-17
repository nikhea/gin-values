package main

import (
	"gin-learn/config"
	"gin-learn/handlers"

	"github.com/gin-gonic/gin"
)


func main(){

	config.LoadEnv()
	config.RunMigrations()
	config.ConnectDatabase()

	router := gin.Default()

	api := router.Group("/api/users")
	{

		api.POST("/", handlers.CreateUser)

		api.GET("/", handlers.GetUsers)

		api.GET("/:id", handlers.GetUser)

		api.PUT("/:id", handlers.UpdateUser)

		api.DELETE("/:id", handlers.DeleteUser)

	}



	router.Run(":8080")
}