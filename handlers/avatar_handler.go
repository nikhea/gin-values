package handlers

import (
	"errors"
	"net/http"
	"os"

	"grip/audit"
	"grip/dto"
	"grip/models"
	services "grip/service"
	"grip/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Compile-time reference so Swagger can resolve dto.AvatarResponse.
var _ = dto.AvatarResponse{}

// uploadAvatar validates and stores a multipart image, then links it via link.
// formKey is the multipart field name ("avatar"), subdir is the avatar dir
// (utils.AvatarUsersDir / utils.AvatarContactsDir).
func uploadAvatar(c *gin.Context, formKey, subdir string, link func(string) (any, error)) {
	file, err := c.FormFile(formKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar image is required (multipart form field \"avatar\")"})
		return
	}
	if err := utils.ValidateUploadSize(file.Size); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.ValidateAvatarExt(file.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlPath, diskPath, err := utils.AvatarTarget(subdir, file.Filename)
	if err != nil {
		internalError(c, err)
		return
	}
	if err := c.SaveUploadedFile(file, diskPath); err != nil {
		internalError(c, err)
		return
	}

	if _, err := link(urlPath); err != nil {
		_ = os.Remove(diskPath) // don't leave orphan files on failure
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Avatar uploaded",
		"avatar_url": urlPath,
	})
	audit.Log(c, models.AuditAvatarUpload, "avatars", "", map[string]any{"avatar_url": urlPath})
}

// UploadUserAvatar godoc
// @Summary Upload a user's avatar image
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" example(458622d8-daba-4252-8ce1-846277353139)
// @Param avatar formData file true "Avatar image (.jpg, .jpeg, .png, .gif, .webp; max 5 MB)"
// @Success 200 {object} dto.AvatarResponse
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /users/{id}/avatar [post]
func UploadUserAvatar(c *gin.Context) {
	id := c.Param("id")
	uploadAvatar(c, "avatar", utils.AvatarUsersDir, func(urlPath string) (any, error) {
		return services.SetUserAvatar(id, urlPath)
	})
}

// UploadContactAvatar godoc
// @Summary Upload a contact's avatar image
// @Tags contacts
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact ID" example(9c9e6679-7425-40de-944b-e07fc1f90ae7)
// @Param avatar formData file true "Avatar image (.jpg, .jpeg, .png, .gif, .webp; max 5 MB)"
// @Success 200 {object} dto.AvatarResponse
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /contacts/{id}/avatar [post]
func UploadContactAvatar(c *gin.Context) {
	id := c.Param("id")
	uploadAvatar(c, "avatar", utils.AvatarContactsDir, func(urlPath string) (any, error) {
		return services.SetContactAvatar(id, urlPath)
	})
}
