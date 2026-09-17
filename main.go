package main

import (
	"gin-learn/handlers"

	"github.com/gin-gonic/gin"
)


func main(){

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