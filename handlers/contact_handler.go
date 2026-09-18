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

// Compile-time reference so Swagger can resolve models.Contact in annotations.
var _ = models.Contact{}

func writeContactError(c *gin.Context, err error, notFoundMsg string) {
	var validationErr *dto.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validationErr.Error(),
		})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": notFoundMsg,
		})
		return
	}
	internalError(c, err)
}

// CreateContact godoc
// @Summary Create contact
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateContactRequest true "Contact payload"
// @Success 201 {object} dto.ContactEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /contacts/ [post]
func CreateContact(c *gin.Context) {
	var req dto.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	contact, err := services.CreateContact(req)
	if err != nil {
		writeContactError(c, err, "User not found")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Contact created",
		"contact": contact,
	})
}

// GetContacts godoc
// @Summary List contacts with pagination and filters
// @Tags contacts
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "Filter by user ID" example(458622d8-daba-4252-8ce1-846277353139)
// @Param type query string false "Filter by type" Enums(email, phone, address, other) example(email)
// @Param search query string false "Search name and value" example(kaige)
// @Param page query int false "Page number" default(1) example(1)
// @Param page_size query int false "Page size" default(10) example(10)
// @Success 200 {object} dto.ContactsEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /contacts/ [get]
func GetContacts(c *gin.Context) {
	var filter dto.ContactFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	filter.Normalize()

	contacts, total, err := services.ListContacts(filter)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Contacts found",
		"contacts": contacts,
		"meta":     dto.NewPageMeta(filter.Pagination, total),
	})
}

// GetContact godoc
// @Summary Get contact by ID
// @Tags contacts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact ID" example(9c9e6679-7425-40de-944b-e07fc1f90ae7)
// @Success 200 {object} models.Contact
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /contacts/{id} [get]
func GetContact(c *gin.Context) {
	contact, err := services.GetContactByID(c.Param("id"))
	if err != nil {
		writeContactError(c, err, "Contact not found")
		return
	}

	c.JSON(http.StatusOK, contact)
}

// UpdateContact godoc
// @Summary Update contact
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact ID" example(9c9e6679-7425-40de-944b-e07fc1f90ae7)
// @Param request body dto.UpdateContactRequest true "Contact payload"
// @Success 200 {object} models.Contact
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /contacts/{id} [put]
func UpdateContact(c *gin.Context) {
	var req dto.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	contact, err := services.UpdateContact(c.Param("id"), req)
	if err != nil {
		writeContactError(c, err, "Contact not found")
		return
	}

	c.JSON(http.StatusOK, contact)
}

// DeleteContact godoc
// @Summary Delete contact
// @Tags contacts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact ID" example(9c9e6679-7425-40de-944b-e07fc1f90ae7)
// @Success 200 {object} dto.MessageEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /contacts/{id} [delete]
func DeleteContact(c *gin.Context) {
	if err := services.DeleteContact(c.Param("id")); err != nil {
		writeContactError(c, err, "Contact not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contact deleted",
	})
}
