package handlers

import (
	"errors"
	"grip/dto"
	"grip/models"
	services "grip/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Compile-time reference so Swagger can resolve models.Profile in annotations.
var _ = models.Profile{}

// CreateProfile godoc
// @Summary Create profile for user
// @Tags profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example(458622d8-daba-4252-8ce1-846277353139)
// @Param request body dto.CreateProfileRequest true "Profile payload"
// @Success 201 {object} dto.ProfileEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/{id}/profile [post]
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
		internalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile created",
		"profile": profile,
	})
}

// GetProfile godoc
// @Summary Get profile by user ID
// @Tags profiles
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example(458622d8-daba-4252-8ce1-846277353139)
// @Success 200 {object} models.Profile
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/{id}/profile [get]
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
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Update profile by user ID
// @Tags profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example(458622d8-daba-4252-8ce1-846277353139)
// @Param request body dto.UpdateProfileRequest true "Profile payload"
// @Success 200 {object} models.Profile
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/{id}/profile [put]
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
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// DeleteProfile godoc
// @Summary Delete profile by user ID
// @Tags profiles
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example(458622d8-daba-4252-8ce1-846277353139)
// @Success 200 {object} dto.MessageEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /users/{id}/profile [delete]
func DeleteProfile(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteProfile(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Profile not found",
			})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted",
	})
}
