package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"grip/audit"
	"grip/dto"
	"grip/jobs"
	"grip/middleware"
	"grip/models"
	services "grip/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Compile-time references so Swagger can resolve org envelopes.
var (
	_ = dto.OrgEnvelope{}
	_ = dto.OrgsEnvelope{}
	_ = dto.MemberEnvelope{}
)

func writeOrgError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidRole):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrAlreadyMember):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrLastOwner):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	default:
		internalError(c, err)
	}
}

func currentUserID(c *gin.Context) (string, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return "", false
	}
	return userID, true
}

// CreateOrg godoc
// @Summary Create an organization (caller becomes owner)
// @Tags orgs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateOrgRequest true "Organization payload"
// @Success 201 {object} dto.OrgEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Router /orgs/ [post]
func CreateOrg(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := services.CreateOrg(userID, req)
	if err != nil {
		writeOrgError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Organization created",
		"org":     org,
	})
	audit.LogFor(c, userID, models.AuditOrgCreate, "orgs", org.ID, "", map[string]any{"name": org.Name})
}

// ListOrgs godoc
// @Summary List my organizations
// @Tags orgs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.OrgsEnvelope
// @Router /orgs/ [get]
func ListOrgs(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	orgs, err := services.ListOrgs(userID)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organizations found",
		"orgs":    orgs,
	})
}

// AddMember godoc
// @Summary Add a member to the organization
// @Tags orgs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Param request body dto.AddMemberRequest true "Member payload"
// @Success 201 {object} dto.MemberEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 403 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Failure 409 {object} dto.ErrorEnvelope
// @Router /orgs/{id}/members [post]
func AddMember(c *gin.Context) {
	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := services.AddMember(c.Param("id"), req.UserID, req.Role)
	if err != nil {
		writeOrgError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Member added",
		"member":  m,
	})
	audit.Log(c, models.AuditMemberAdd, "members", m.UserID, map[string]any{"org_id": m.OrgID, "role": m.Role})
	if err := jobs.EnqueueNotify(context.Background(), m.UserID, models.NotifyAddedToOrg,
		"Added to organization",
		"You were added to an organization as "+m.Role+".",
		map[string]any{"org_id": m.OrgID, "role": m.Role}); err != nil {
		slog.Warn("added-to-org notification failed", "error", err, "user", m.UserID)
	}
}

// UpdateMember godoc
// @Summary Change a member's role
// @Tags orgs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Param uid path string true "Member user ID"
// @Param request body dto.UpdateMemberRequest true "Role payload"
// @Success 200 {object} dto.MemberEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 403 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /orgs/{id}/members/{uid} [put]
func UpdateMember(c *gin.Context) {
	var req dto.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := services.UpdateMemberRole(c.Param("id"), c.Param("uid"), req.Role)
	if err != nil {
		writeOrgError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Member updated",
		"member":  m,
	})
	audit.Log(c, models.AuditMemberUpdate, "members", m.UserID, map[string]any{"org_id": m.OrgID, "role": m.Role})
	if err := jobs.EnqueueNotify(context.Background(), m.UserID, models.NotifyRoleChanged,
		"Role changed",
		"Your role is now "+m.Role+".",
		map[string]any{"org_id": m.OrgID, "role": m.Role}); err != nil {
		slog.Warn("role-changed notification failed", "error", err, "user", m.UserID)
	}
}

// RemoveMember godoc
// @Summary Remove a member from the organization
// @Tags orgs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Param uid path string true "Member user ID"
// @Success 200 {object} dto.MessageEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 403 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /orgs/{id}/members/{uid} [delete]
func RemoveMember(c *gin.Context) {
	if err := services.RemoveMember(c.Param("id"), c.Param("uid")); err != nil {
		writeOrgError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
	audit.Log(c, models.AuditMemberRemove, "members", c.Param("uid"), nil)
}

// DeleteOrg godoc
// @Summary Delete the organization
// @Tags orgs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} dto.MessageEnvelope
// @Failure 403 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /orgs/{id} [delete]
func DeleteOrg(c *gin.Context) {
	if err := services.DeleteOrg(c.Param("id")); err != nil {
		writeOrgError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted"})
	audit.Log(c, models.AuditOrgDelete, "orgs", c.Param("id"), nil)
}
