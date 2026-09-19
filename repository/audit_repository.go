package repository

import (
	"time"

	"grip/config"
	"grip/models"
)

// CreateAuditLog appends one immutable row. There is intentionally no
// update or delete function for audit logs.
func CreateAuditLog(entry *models.AuditLog) error {
	return config.DB.Create(entry).Error
}

// AuditFilter scopes audit listing. Zero values mean "no constraint".
type AuditFilter struct {
	ActorID      string
	Action       string
	ResourceType string
	OrgID        string
	From         time.Time
	To           time.Time
	Page         int
	PageSize     int
}

// Normalize applies pagination defaults (page 1, 10 per page, max 100).
func (f *AuditFilter) Normalize() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 10
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
}

// Offset returns the SQL offset for the current page.
func (f AuditFilter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// ListAuditLogs returns newest-first rows matching the filter plus total.
func ListAuditLogs(filter AuditFilter) ([]models.AuditLog, int64, error) {
	filter.Normalize()

	var rows []models.AuditLog
	var total int64

	q := config.DB.Model(&models.AuditLog{})
	if filter.ActorID != "" {
		q = q.Where("actor_id = ?", filter.ActorID)
	}
	if filter.Action != "" {
		q = q.Where("action = ?", filter.Action)
	}
	if filter.ResourceType != "" {
		q = q.Where("resource_type = ?", filter.ResourceType)
	}
	if filter.OrgID != "" {
		q = q.Where("org_id = ?", filter.OrgID)
	}
	if !filter.From.IsZero() {
		q = q.Where("created_at >= ?", filter.From)
	}
	if !filter.To.IsZero() {
		q = q.Where("created_at <= ?", filter.To)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Limit(filter.PageSize).Offset(filter.Offset()).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// GetAuditLogByID loads one row, scoped to an org when orgID != "".
func GetAuditLogByID(id, orgID string) (*models.AuditLog, error) {
	var row models.AuditLog
	q := config.DB.Where("id = ?", id)
	if orgID != "" {
		q = q.Where("org_id = ?", orgID)
	}
	if err := q.First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
