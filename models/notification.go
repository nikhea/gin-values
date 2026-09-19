package models

import "time"

// Notification types emitted by the app (see jobs NotifyWorker).
const (
	NotifyWelcome         = "welcome"
	NotifyPasswordChanged = "password_changed"
	NotifyImportFinished  = "import_finished"
	NotifyAddedToOrg      = "added_to_org"
	NotifyRoleChanged     = "role_changed"
)

// Notification is one inbox row for a user. Rows are never updated except
// for the read marker; Data carries event context (counts, org/role).
type Notification struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`

	UserID string `gorm:"index;not null;type:varchar(36)" json:"user_id"`

	Type string `gorm:"type:varchar(32)" json:"type" example:"welcome"`

	Title string `json:"title" example:"Welcome to Grip"`

	Body string `json:"body" example:"Your email is verified. Glad to have you!"`

	Data map[string]any `gorm:"type:jsonb;serializer:json" json:"data"`

	ReadAt *time.Time `gorm:"index" json:"read_at"`

	CreatedAt time.Time `json:"created_at"`
}
