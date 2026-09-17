package handlers

import (
	"gin-learn/dto"
	"gin-learn/models"
	services "gin-learn/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CREATE USER
func CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); 
	
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := services.CreateUser(req)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"user": user,
	})
}



// GET ALL USERS
func GetUsers(c *gin.Context){

	user := services.GetUsers()

	if len(user) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message":"No users yet",
			"user": []models.User{},
		})
		return
	}

	c.JSON(http.StatusOK,  gin.H{
		"message":"Users found",
		"user": user,
	})
}



// GET SINGLE USER

func GetUser(c *gin.Context){
	id := c.Param("id")

	 users, found := services.GetUserByID(id)

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"message":"User not found",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}



// UPDATE USER

func UpdateUser(c *gin.Context){
		var req dto.CreateUserRequest
			id := c.Param("id")


	if err := c.ShouldBindJSON(&req); 
	
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

		 users, found := services.UpdateUser(id, req)


		 if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"message":"User not found",
		})
		return
	}


		c.JSON(http.StatusOK, users)

		

}




// DELETE USER

func DeleteUser(c *gin.Context){

	id := c.Param("id")

	 found := services.DeleteUser(id)

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"message":"User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":"User deleted",
	})
}