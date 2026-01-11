package models

import (
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of action performed on a secret.
type AuditAction string

const (
	ActionCreated      AuditAction = "created"
	ActionViewed       AuditAction = "viewed"
	ActionDestroyed    AuditAction = "destroyed"
	ActionExpired      AuditAction = "expired"
	ActionFailedAccess AuditAction = "failed_access"
)

// AuditLog represents an immutable audit trail entry.
type AuditLog struct {
	ID        int64                  `json:"id" db:"id"`
	SecretID  uuid.UUID              `json:"secret_id" db:"secret_id"`
	Action    AuditAction            `json:"action" db:"action"`
	IPAddress *string                `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string                `json:"user_agent,omitempty" db:"user_agent"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	Timestamp time.Time              `json:"timestamp" db:"timestamp"`
}

// NewAuditLog creates a new audit log entry.
func NewAuditLog(secretID uuid.UUID, action AuditAction, ipAddress, userAgent string) *AuditLog {
	var ip, ua *string
	if ipAddress != "" {
		ip = &ipAddress
	}
	if userAgent != "" {
		ua = &userAgent
	}

	return &AuditLog{
		SecretID:  secretID,
		Action:    action,
		IPAddress: ip,
		UserAgent: ua,
		Timestamp: time.Now(),
	}
}
