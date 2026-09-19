package handlers

import (
	"errors"
	"grip/audit"
	"grip/dto"
	"grip/middleware"
	"grip/models"
	services "grip/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Compile-time reference so Swagger can resolve models.Contact in annotations.
var _ = models.Contact{}

// contactOrgMeta attaches the team scope to audit rows when present.
func contactOrgMeta(contact *models.Contact) map[string]any {
	if contact.OrgID != nil && *contact.OrgID != "" {
		return map[string]any{"org_id": *contact.OrgID}
	}
	return nil
}

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
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	var req dto.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	contact, err := services.CreateContactForUser(actorID, req)
	if err != nil {
		writeContactError(c, err, "User not found")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Contact created",
		"contact": contact,
	})
	audit.Log(c, models.AuditContactCreate, "contacts", contact.ID, contactOrgMeta(contact))
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
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	var filter dto.ContactFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Org scope also accepted via header; explicit query wins.
	if filter.OrgID == "" {
		filter.OrgID = middleware.OrgIDFromRequest(c)
	}
	filter.Normalize()

	contacts, total, err := services.ListContactsForUser(actorID, filter)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
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
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	contact, err := services.GetContactForUser(actorID, c.Param("id"))
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
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	var req dto.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	contact, err := services.UpdateContactForUser(actorID, c.Param("id"), req)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		writeContactError(c, err, "Contact not found")
		return
	}

	audit.Log(c, models.AuditContactUpdate, "contacts", contact.ID, contactOrgMeta(contact))

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
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	if err := services.DeleteContactForUser(actorID, c.Param("id")); err != nil {
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		writeContactError(c, err, "Contact not found")
		return
	}

	audit.Log(c, models.AuditContactDelete, "contacts", c.Param("id"), nil)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Contact deleted",
		"soft_deleted": true,
	})
}
