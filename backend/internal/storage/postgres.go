package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ranvijayysinghrathore/envsend/backend/internal/models"
)

// PostgresRepository handles database operations for secrets and audit logs.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(databaseURL string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Production Tuning: Connection Pool Settings
	// Limit max open connections to prevent crashing the database (e.g., max connections limit on RDS/Render)
	db.SetMaxOpenConns(25)
	// Keep a pool of idle connections ready for fast access
	db.SetMaxIdleConns(25)
	// Recycle connections every 5 minutes to prevent stale connection issues
	db.SetConnMaxLifetime(5 * time.Minute)

	repo := &PostgresRepository{db: db}

	// Run migrations automatically
	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return repo, nil
}

// Close closes the database connection.
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// CreateSecret creates a new secret in the database.
func (r *PostgresRepository) CreateSecret(ctx context.Context, secret *models.Secret) error {
	metadataJSON, err := json.Marshal(secret.EncryptionMetadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO secrets (
			id, encrypted_blob_url, encryption_metadata, expires_at, 
			max_views, view_count, ip_lock, recipient_identifier, 
			destroyed, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = r.db.ExecContext(ctx, query,
		secret.ID,
		secret.EncryptedBlobURL,
		metadataJSON,
		secret.ExpiresAt,
		secret.MaxViews,
		secret.ViewCount,
		secret.IPLock,
		secret.RecipientID,
		secret.Destroyed,
		secret.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

// GetSecretByID retrieves a secret by its ID.
func (r *PostgresRepository) GetSecretByID(ctx context.Context, id uuid.UUID) (*models.Secret, error) {
	query := `
		SELECT id, encrypted_blob_url, encryption_metadata, expires_at, 
		       max_views, view_count, ip_lock, recipient_identifier, 
		       destroyed, created_at, accessed_at, destroyed_at
		FROM secrets
		WHERE id = $1
	`

	var secret models.Secret
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&secret.ID,
		&secret.EncryptedBlobURL,
		&metadataJSON,
		&secret.ExpiresAt,
		&secret.MaxViews,
		&secret.ViewCount,
		&secret.IPLock,
		&secret.RecipientID,
		&secret.Destroyed,
		&secret.CreatedAt,
		&secret.AccessedAt,
		&secret.DestroyedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("secret not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	// Unmarshal metadata
	if err := json.Unmarshal(metadataJSON, &secret.EncryptionMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &secret, nil
}

// IncrementViewCount increments the view count for a secret.
func (r *PostgresRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE secrets
		SET view_count = view_count + 1,
		    accessed_at = $1
		WHERE id = $2 AND NOT destroyed
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("secret not found or already destroyed")
	}

	return nil
}

// DeleteSecret deletes a secret from the database.
func (r *PostgresRepository) DeleteSecret(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM secrets WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

// GetExpiredSecrets retrieves all expired secrets.
func (r *PostgresRepository) GetExpiredSecrets(ctx context.Context, limit int) ([]*models.Secret, error) {
	query := `
		SELECT id, encrypted_blob_url, encryption_metadata, expires_at, 
		       max_views, view_count, ip_lock, recipient_identifier, 
		       destroyed, created_at, accessed_at, destroyed_at
		FROM secrets
		WHERE (expires_at < $1 OR destroyed = true)
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired secrets: %w", err)
	}
	defer rows.Close()

	var secrets []*models.Secret
	for rows.Next() {
		var secret models.Secret
		var metadataJSON []byte

		err := rows.Scan(
			&secret.ID,
			&secret.EncryptedBlobURL,
			&metadataJSON,
			&secret.ExpiresAt,
			&secret.MaxViews,
			&secret.ViewCount,
			&secret.IPLock,
			&secret.RecipientID,
			&secret.Destroyed,
			&secret.CreatedAt,
			&secret.AccessedAt,
			&secret.DestroyedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}

		if err := json.Unmarshal(metadataJSON, &secret.EncryptionMetadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		secrets = append(secrets, &secret)
	}

	return secrets, nil
}

// CreateAuditLog creates a new audit log entry.
func (r *PostgresRepository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	var metadataJSON []byte
	var err error

	if log.Metadata != nil {
		metadataJSON, err = json.Marshal(log.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	query := `
		INSERT INTO audit_logs (secret_id, action, ip_address, user_agent, metadata, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.db.ExecContext(ctx, query,
		log.SecretID,
		log.Action,
		log.IPAddress,
		log.UserAgent,
		metadataJSON,
		log.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetAuditLogsBySecretID retrieves all audit logs for a secret.
func (r *PostgresRepository) GetAuditLogsBySecretID(ctx context.Context, secretID uuid.UUID) ([]*models.AuditLog, error) {
	query := `
		SELECT id, secret_id, action, ip_address, user_agent, metadata, timestamp
		FROM audit_logs
		WHERE secret_id = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, secretID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.SecretID,
			&log.Action,
			&log.IPAddress,
			&log.UserAgent,
			&metadataJSON,
			&log.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &log.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		logs = append(logs, &log)
	}

	return logs, nil
}
