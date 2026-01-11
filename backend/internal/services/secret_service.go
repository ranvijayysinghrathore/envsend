package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/envsend/backend/internal/models"
	"github.com/yourusername/envsend/backend/internal/storage"
)

// SecretService handles business logic for secrets.
type SecretService struct {
	postgres *storage.PostgresRepository
	s3       *storage.S3Storage
	redis    *storage.RedisClient
}

// NewSecretService creates a new secret service.
func NewSecretService(postgres *storage.PostgresRepository, s3 *storage.S3Storage, redis *storage.RedisClient) *SecretService {
	return &SecretService{
		postgres: postgres,
		s3:       s3,
		redis:    redis,
	}
}

// CreateSecret creates a new encrypted secret.
func (s *SecretService) CreateSecret(ctx context.Context, req models.CreateSecretRequest, ipAddress string) (*models.CreateSecretResponse, error) {
	// Parse expiry duration
	expiresIn, err := time.ParseDuration(req.ExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("invalid expiry duration: %w", err)
	}

	// Decode encrypted blob from base64
	encryptedData, err := base64.StdEncoding.DecodeString(req.EncryptedBlob)
	if err != nil {
		return nil, fmt.Errorf("invalid encrypted blob encoding: %w", err)
	}

	// Upload encrypted blob to S3
	blobURL, err := s.s3.UploadEncryptedBlob(ctx, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to upload blob: %w", err)
	}

	// Create secret metadata
	secret := &models.Secret{
		ID:                 uuid.New(),
		EncryptedBlobURL:   blobURL,
		EncryptionMetadata: req.EncryptionMetadata,
		ExpiresAt:          time.Now().Add(expiresIn),
		MaxViews:           req.MaxViews,
		ViewCount:          0,
		Destroyed:          false,
		CreatedAt:          time.Now(),
	}

	if req.IPLock != "" {
		secret.IPLock = &req.IPLock
	}
	if req.RecipientID != "" {
		secret.RecipientID = &req.RecipientID
	}

	// Save to database
	if err := s.postgres.CreateSecret(ctx, secret); err != nil {
		// Cleanup: delete uploaded blob
		_ = s.s3.DeleteEncryptedBlob(ctx, blobURL)
		return nil, fmt.Errorf("failed to save secret: %w", err)
	}

	// Create audit log
	auditLog := models.NewAuditLog(secret.ID, models.ActionCreated, ipAddress, "")
	_ = s.postgres.CreateAuditLog(ctx, auditLog)

	// Build response
	response := &models.CreateSecretResponse{
		SecretID:  secret.ID.String(),
		URL:       fmt.Sprintf("/s/%s", secret.ID.String()),
		ExpiresAt: secret.ExpiresAt,
		MaxViews:  secret.MaxViews,
	}

	return response, nil
}

// GetSecret retrieves and decrypts a secret (server never sees plaintext).
func (s *SecretService) GetSecret(ctx context.Context, secretID string, ipAddress, userAgent string) (*models.GetSecretResponse, error) {
	// Parse UUID
	id, err := uuid.Parse(secretID)
	if err != nil {
		return nil, fmt.Errorf("invalid secret ID: %w", err)
	}

	// Get secret from database
	secret, err := s.postgres.GetSecretByID(ctx, id)
	if err != nil {
		// Log failed access attempt
		auditLog := models.NewAuditLog(id, models.ActionFailedAccess, ipAddress, userAgent)
		_ = s.postgres.CreateAuditLog(ctx, auditLog)
		return nil, fmt.Errorf("secret not found: %w", err)
	}

	// Check if secret can be accessed
	if !secret.CanAccess() {
		action := models.ActionFailedAccess
		if secret.Destroyed {
			action = models.ActionDestroyed
		} else if secret.IsExpired() {
			action = models.ActionExpired
		}

		auditLog := models.NewAuditLog(id, action, ipAddress, userAgent)
		_ = s.postgres.CreateAuditLog(ctx, auditLog)

		return nil, fmt.Errorf("secret cannot be accessed (destroyed or expired)")
	}

	// Check IP lock
	if secret.IPLock != nil && *secret.IPLock != ipAddress {
		auditLog := models.NewAuditLog(id, models.ActionFailedAccess, ipAddress, userAgent)
		auditLog.Metadata = map[string]interface{}{
			"reason":       "ip_mismatch",
			"expected_ip":  *secret.IPLock,
			"actual_ip":    ipAddress,
		}
		_ = s.postgres.CreateAuditLog(ctx, auditLog)

		return nil, fmt.Errorf("IP address mismatch")
	}

	// Download encrypted blob
	encryptedData, err := s.s3.DownloadEncryptedBlob(ctx, secret.EncryptedBlobURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob: %w", err)
	}

	// Increment view count
	if err := s.postgres.IncrementViewCount(ctx, id); err != nil {
		return nil, fmt.Errorf("failed to increment view count: %w", err)
	}

	// Create audit log
	auditLog := models.NewAuditLog(id, models.ActionViewed, ipAddress, userAgent)
	_ = s.postgres.CreateAuditLog(ctx, auditLog)

	// Check if should be destroyed after this view
	secret.IncrementViewCount()
	if secret.ShouldDestroy() {
		// Schedule async deletion
		go func() {
			deleteCtx := context.Background()
			_ = s.DeleteSecret(deleteCtx, secretID)
		}()
	}

	// Build response
	response := &models.GetSecretResponse{
		EncryptedBlob:      base64.StdEncoding.EncodeToString(encryptedData),
		EncryptionMetadata: secret.EncryptionMetadata,
		ViewsRemaining:     secret.MaxViews - secret.ViewCount,
		ExpiresAt:          secret.ExpiresAt,
	}

	return response, nil
}

// DeleteSecret deletes a secret and its encrypted blob.
func (s *SecretService) DeleteSecret(ctx context.Context, secretID string) error {
	id, err := uuid.Parse(secretID)
	if err != nil {
		return fmt.Errorf("invalid secret ID: %w", err)
	}

	// Get secret to find blob URL
	secret, err := s.postgres.GetSecretByID(ctx, id)
	if err != nil {
		return fmt.Errorf("secret not found: %w", err)
	}

	// Delete encrypted blob from S3
	if err := s.s3.DeleteEncryptedBlob(ctx, secret.EncryptedBlobURL); err != nil {
		// Log error but continue with database deletion
		fmt.Printf("Warning: failed to delete blob: %v\n", err)
	}

	// Delete from database
	if err := s.postgres.DeleteSecret(ctx, id); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	// Create audit log
	auditLog := models.NewAuditLog(id, models.ActionDestroyed, "", "")
	_ = s.postgres.CreateAuditLog(ctx, auditLog)

	return nil
}

// CleanupExpiredSecrets removes expired secrets (called by worker).
func (s *SecretService) CleanupExpiredSecrets(ctx context.Context, batchSize int) (int, error) {
	secrets, err := s.postgres.GetExpiredSecrets(ctx, batchSize)
	if err != nil {
		return 0, fmt.Errorf("failed to get expired secrets: %w", err)
	}

	count := 0
	for _, secret := range secrets {
		if err := s.DeleteSecret(ctx, secret.ID.String()); err != nil {
			fmt.Printf("Warning: failed to delete secret %s: %v\n", secret.ID, err)
			continue
		}
		count++
	}

	return count, nil
}
