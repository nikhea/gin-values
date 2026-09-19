package handlers

import (
	"errors"
	"net/http"
	"time"

	"grip/dto"
	"grip/middleware"
	"grip/models"
	"grip/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	_ = dto.AuditEnvelope{}
	_ = models.AuditLog{}
)

func parseAuditTime(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, raw)
}

// ListAuditLogs godoc
// @Summary List organization audit logs (admin+)
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param actor_id query string false "Filter by actor user ID"
// @Param action query string false "Filter by action" example(contacts.create)
// @Param resource_type query string false "Filter by resource" example(contacts)
// @Param from query string false "RFC3339 start" example(2026-01-01T00:00:00Z)
// @Param to query string false "RFC3339 end" example(2026-12-31T23:59:59Z)
// @Param page query int false "Page number" default(1) example(1)
// @Param page_size query int false "Page size" default(10) example(10)
// @Success 200 {object} dto.AuditEnvelope
// @Failure 403 {object} dto.ErrorEnvelope
// @Router /audit-logs/ [get]
func ListAuditLogs(c *gin.Context) {
	orgID, ok := middleware.GetOrgID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Org-ID header (or org_id query param) is required"})
		return
	}

	var q dto.AuditQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from, err := parseAuditTime(q.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from time (RFC3339)"})
		return
	}
	to, err := parseAuditTime(q.To)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to time (RFC3339)"})
		return
	}

	filter := repository.AuditFilter{
		ActorID: q.ActorID, Action: q.Action, ResourceType: q.ResourceType,
		OrgID: orgID, From: from, To: to, Page: q.Page, PageSize: q.PageSize,
	}
	filter.Normalize()

	rows, total, err := repository.ListAuditLogs(filter)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Audit logs found",
		"logs":    rows,
		"meta": dto.PageMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
		},
	})
}

// GetAuditLog godoc
// @Summary Get one organization audit log (admin+)
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audit log ID"
// @Success 200 {object} models.AuditLog
// @Failure 403 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /audit-logs/{id} [get]
func GetAuditLog(c *gin.Context) {
	orgID, ok := middleware.GetOrgID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Org-ID header (or org_id query param) is required"})
		return
	}

	row, err := repository.GetAuditLogByID(c.Param("id"), orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Audit log not found"})
			return
		}
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, row)
}
