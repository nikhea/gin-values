// Package audit records an append-only trail of state-changing API calls.
// Handlers call Log after a successful mutation; failures inside Log are
// swallowed (logged) so auditing can never break a request.
package audit

import (
	"log/slog"
	"time"

	"grip/middleware"
	"grip/models"
	"grip/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Log writes one audit row. actorID/orgID come from the request context
// (empty when anonymous); meta carries free-form context and may be nil.
func Log(c *gin.Context, action, resourceType, resourceID string, meta map[string]any) {
	var actorID, orgID, requestID *string
	if id, ok := middleware.GetUserID(c); ok {
		actorID = &id
	}
	if id, ok := middleware.GetOrgID(c); ok {
		orgID = &id
	}
	if rid, exists := c.Get("requestID"); exists {
		if s, ok := rid.(string); ok && s != "" {
			requestID = &s
		}
	}
	var rid *string
	if resourceID != "" {
		rid = &resourceID
	}
	if meta == nil {
		meta = map[string]any{}
	}

	entry := &models.AuditLog{
		ID:           uuid.New().String(),
		ActorID:      actorID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   rid,
		OrgID:        orgID,
		IP:           c.ClientIP(),
		RequestID:    requestID,
		Metadata:     meta,
		CreatedAt:    time.Now(),
	}
	if err := repository.CreateAuditLog(entry); err != nil {
		slog.Error("audit log failed", "error", err, "action", action)
	}
}

// LogFor records an audit row for a known actor outside AuthRequired
// (e.g. register/login on anonymous routes). Context-derived fields
// (IP, request ID) are still captured; prefer Log in handlers.
func LogFor(c *gin.Context, actorID, action, resourceType, resourceID, orgID string, meta map[string]any) {
	var actor, org, rid, requestID *string
	if actorID != "" {
		actor = &actorID
	}
	if orgID != "" {
		org = &orgID
	}
	if resourceID != "" {
		rid = &resourceID
	}
	if c != nil {
		if v, ok := c.Get("requestID"); ok {
			if s, ok := v.(string); ok && s != "" {
				requestID = &s
			}
		}
	}
	if meta == nil {
		meta = map[string]any{}
	}
	entry := &models.AuditLog{
		ID:           uuid.New().String(),
		ActorID:      actor,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   rid,
		OrgID:        org,
		RequestID:    requestID,
		Metadata:     meta,
		CreatedAt:    time.Now(),
	}
	if c != nil {
		entry.IP = c.ClientIP()
	}
	if err := repository.CreateAuditLog(entry); err != nil {
		slog.Error("audit log failed", "error", err, "action", action)
	}
}
