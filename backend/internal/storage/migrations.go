package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	ctx := context.Background()

	// Check if migrations table exists
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Check if migration 001 already applied
	var exists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = 1)").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if exists {
		log.Println("✓ Database migrations already applied")
		return nil
	}

	// Run migration 001
	log.Println("Running database migration 001...")
	
	migration := `
-- Create secrets table
CREATE TABLE IF NOT EXISTS secrets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encrypted_blob_url TEXT NOT NULL,
    encryption_metadata JSONB NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    max_views INTEGER NOT NULL DEFAULT 1,
    view_count INTEGER NOT NULL DEFAULT 0,
    ip_lock INET,
    recipient_identifier TEXT,
    destroyed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    accessed_at TIMESTAMP WITH TIME ZONE,
    destroyed_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT valid_max_views CHECK (max_views > 0),
    CONSTRAINT valid_view_count CHECK (view_count >= 0),
    CONSTRAINT view_count_not_exceed_max CHECK (view_count <= max_views)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    secret_id UUID NOT NULL REFERENCES secrets(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_action CHECK (action IN ('created', 'viewed', 'destroyed', 'expired', 'failed_access'))
);

CREATE INDEX IF NOT EXISTS idx_secrets_expires_at ON secrets(expires_at) WHERE NOT destroyed;
CREATE INDEX IF NOT EXISTS idx_secrets_destroyed ON secrets(destroyed);
CREATE INDEX IF NOT EXISTS idx_secrets_created_at ON secrets(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_secret_id ON audit_logs(secret_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

CREATE OR REPLACE FUNCTION check_secret_destruction()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.view_count >= NEW.max_views THEN
        NEW.destroyed := TRUE;
        NEW.destroyed_at := NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_secret_destruction
    BEFORE UPDATE OF view_count ON secrets
    FOR EACH ROW
    EXECUTE FUNCTION check_secret_destruction();

CREATE OR REPLACE FUNCTION prevent_destroyed_secret_access()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.destroyed = TRUE THEN
        RAISE EXCEPTION 'Cannot access destroyed secret';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_prevent_destroyed_access
    BEFORE UPDATE OF view_count ON secrets
    FOR EACH ROW
    EXECUTE FUNCTION prevent_destroyed_secret_access();
`

	_, err = db.ExecContext(ctx, migration)
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	// Mark migration as applied
	_, err = db.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (1)")
	if err != nil {
		return fmt.Errorf("failed to mark migration as applied: %w", err)
	}

	log.Println("✓ Database migration 001 completed successfully")
	return nil
}
