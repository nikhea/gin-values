package handlers

import (
	"errors"
	"gin-learn/dto"
	"gin-learn/models"
	services "gin-learn/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser godoc
// @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User payload"
// @Success 201 {object} dto.UserEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/ [post]
func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"user":    user,
	})
}

// GetUsers godoc
// @Summary List users
// @Tags users
// @Produce json
// @Success 200 {object} dto.UsersEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/ [get]
func GetUsers(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No users yet",
			"user":    []models.User{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users found",
		"user":    users,
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := services.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.CreateUserRequest true "User payload"
// @Success 200 {object} models.User
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.UpdateUser(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.MessageEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
	})
}
