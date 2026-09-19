package models

import "time"

// Canonical audit actions. Handlers/services must use these constants so
// the log stays queryable.
const (
	AuditRegister       = "users.register"
	AuditLogin          = "users.login"
	AuditVerifyEmail    = "users.verify_email"
	AuditVerifyOTP      = "users.verify_otp"
	AuditForgotPassword = "users.forgot_password"
	AuditResetPassword  = "users.reset_password"
	AuditLogout         = "users.logout"
	AuditLogoutAll      = "users.logout_all"
	AuditUserCreate     = "users.create"
	AuditUserUpdate     = "users.update"
	AuditUserDelete     = "users.delete"
	AuditProfileCreate  = "profiles.create"
	AuditProfileUpdate  = "profiles.update"
	AuditProfileDelete  = "profiles.delete"
	AuditContactCreate  = "contacts.create"
	AuditContactUpdate  = "contacts.update"
	AuditContactDelete  = "contacts.delete"
	AuditContactImport  = "contacts.import"
	AuditAvatarUpload   = "avatars.upload"
	AuditOrgCreate      = "orgs.create"
	AuditOrgDelete      = "orgs.delete"
	AuditMemberAdd      = "orgs.member_add"
	AuditMemberUpdate   = "orgs.member_update"
	AuditMemberRemove   = "orgs.member_remove"
)

// AuditLog is an immutable record of a state-changing API call.
// ActorID is nil for anonymous actions (e.g. failed logins); metadata
// carries free-form context (emails, counts, role changes).
type AuditLog struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`

	ActorID *string `gorm:"index;type:varchar(36)" json:"actor_id"`

	Action string `gorm:"index;type:varchar(64)" json:"action" example:"contacts.create"`

	ResourceType string `gorm:"type:varchar(32)" json:"resource_type" example:"contacts"`

	ResourceID *string `gorm:"type:varchar(36)" json:"resource_id"`

	OrgID *string `gorm:"index;type:varchar(36)" json:"org_id"`

	IP string `json:"ip" example:"203.0.113.7"`

	RequestID *string `gorm:"type:varchar(36)" json:"request_id"`

	Metadata map[string]any `gorm:"type:jsonb;serializer:json" json:"metadata"`

	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
