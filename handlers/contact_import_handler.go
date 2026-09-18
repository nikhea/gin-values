package handlers

import (
	"errors"
	"io"
	"net/http"

	"gin-learn/dto"
	services "gin-learn/service"
	"gin-learn/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Compile-time reference so Swagger can resolve dto.ImportSummary.
var _ = dto.ImportSummary{}

func importContacts(c *gin.Context, ext string, parse func(string, io.Reader) (*dto.ImportSummary, error)) {
	userID := c.PostForm("user_id")
	if userID == "" {
		userID = c.Query("user_id")
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required (form field or query param)"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required (multipart form field \"file\")"})
		return
	}
	if err := utils.ValidateUploadSize(file.Size); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.ValidateImportExt(file.Filename, ext); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read uploaded file"})
		return
	}
	defer src.Close()

	summary, err := parse(userID, src)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ImportContactsCSV godoc
// @Summary Import contacts from a CSV file
// @Tags contacts
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param user_id formData string true "Owner user ID"
// @Param file formData file true "CSV file with header name,type,value"
// @Success 200 {object} dto.ImportSummary
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /contacts/import/csv [post]
func ImportContactsCSV(c *gin.Context) {
	importContacts(c, ".csv", func(userID string, r io.Reader) (*dto.ImportSummary, error) {
		return services.ImportContactsFromCSV(userID, r)
	})
}

// ImportContactsJSON godoc
// @Summary Import contacts from a JSON file
// @Tags contacts
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param user_id formData string true "Owner user ID"
// @Param file formData file true "JSON file with an array of {name,type,value}"
// @Success 200 {object} dto.ImportSummary
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /contacts/import/json [post]
func ImportContactsJSON(c *gin.Context) {
	importContacts(c, ".json", func(userID string, r io.Reader) (*dto.ImportSummary, error) {
		return services.ImportContactsFromJSON(userID, r)
	})
}
