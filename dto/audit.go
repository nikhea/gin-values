package dto

import "grip/models"

// AuditQuery carries the list endpoint's filter params. Org scope comes
// from X-Org-ID / org_id (required); other filters are optional.
type AuditQuery struct {
	ActorID      string `form:"actor_id"`
	Action       string `form:"action" example:"contacts.create"`
	ResourceType string `form:"resource_type" example:"contacts"`
	From         string `form:"from" example:"2026-01-01T00:00:00Z"`
	To           string `form:"to" example:"2026-12-31T23:59:59Z"`
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
}

type AuditEnvelope struct {
	Message string            `json:"message"`
	Logs    []models.AuditLog `json:"logs"`
	Meta    PageMeta          `json:"meta"`
}
