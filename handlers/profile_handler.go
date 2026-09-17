package handlers

import (
	"errors"
	"gin-learn/dto"
	services "gin-learn/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CREATE PROFILE — POST /api/users/:id/profile
func CreateProfile(c *gin.Context) {
	id := c.Param("id")

	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	profile, err := services.CreateProfile(id, req)
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile created",
		"profile": profile,
	})
}

// GET PROFILE — GET /api/users/:id/profile
func GetProfile(c *gin.Context) {
	id := c.Param("id")

	profile, err := services.GetProfileByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UPDATE PROFILE — PUT /api/users/:id/profile
func UpdateProfile(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	profile, err := services.UpdateProfile(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// DELETE PROFILE — DELETE /api/users/:id/profile
func DeleteProfile(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteProfile(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted",
	})
}
