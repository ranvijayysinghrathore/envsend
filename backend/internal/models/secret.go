package models

import (
	"time"

	"github.com/google/uuid"
)

// Secret represents the metadata for an encrypted secret (never stores plaintext).
type Secret struct {
	ID                 uuid.UUID              `json:"id" db:"id"`
	EncryptedBlobURL   string                 `json:"encrypted_blob_url" db:"encrypted_blob_url"`
	EncryptionMetadata map[string]interface{} `json:"encryption_metadata" db:"encryption_metadata"`
	ExpiresAt          time.Time              `json:"expires_at" db:"expires_at"`
	MaxViews           int                    `json:"max_views" db:"max_views"`
	ViewCount          int                    `json:"view_count" db:"view_count"`
	IPLock             *string                `json:"ip_lock,omitempty" db:"ip_lock"`
	RecipientID        *string                `json:"recipient_id,omitempty" db:"recipient_identifier"`
	Destroyed          bool                   `json:"destroyed" db:"destroyed"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	AccessedAt         *time.Time             `json:"accessed_at,omitempty" db:"accessed_at"`
	DestroyedAt        *time.Time             `json:"destroyed_at,omitempty" db:"destroyed_at"`
}

// IsExpired checks if the secret has expired.
func (s *Secret) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// CanAccess checks if the secret can be accessed.
func (s *Secret) CanAccess() bool {
	return !s.Destroyed && !s.IsExpired() && s.ViewCount < s.MaxViews
}

// ShouldDestroy checks if the secret should be destroyed.
func (s *Secret) ShouldDestroy() bool {
	return s.Destroyed || s.IsExpired() || s.ViewCount >= s.MaxViews
}

// IncrementViewCount increments the view count and marks as destroyed if max reached.
func (s *Secret) IncrementViewCount() {
	s.ViewCount++
	now := time.Now()
	s.AccessedAt = &now

	if s.ViewCount >= s.MaxViews {
		s.Destroyed = true
		s.DestroyedAt = &now
	}
}

// CreateSecretRequest represents the API request to create a secret.
type CreateSecretRequest struct {
	EncryptedBlob      string                 `json:"encrypted_blob" binding:"required"`
	EncryptionMetadata map[string]interface{} `json:"encryption_metadata" binding:"required"`
	ExpiresIn          string                 `json:"expires_in" binding:"required"` // e.g., "10m", "1h"
	MaxViews           int                    `json:"max_views" binding:"required,min=1"`
	IPLock             string                 `json:"ip_lock,omitempty"`
	RecipientID        string                 `json:"recipient_id,omitempty"`
}

// CreateSecretResponse represents the API response after creating a secret.
type CreateSecretResponse struct {
	SecretID  string    `json:"secret_id"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	MaxViews  int       `json:"max_views"`
}

// GetSecretResponse represents the API response when retrieving a secret.
type GetSecretResponse struct {
	EncryptedBlob      string                 `json:"encrypted_blob"`
	EncryptionMetadata map[string]interface{} `json:"encryption_metadata"`
	ViewsRemaining     int                    `json:"views_remaining"`
	ExpiresAt          time.Time              `json:"expires_at"`
}
