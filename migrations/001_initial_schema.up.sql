-- Create secrets table
CREATE TABLE IF NOT EXISTS secrets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encrypted_blob_url TEXT NOT NULL,
    encryption_metadata JSONB NOT NULL, -- stores algorithm, key derivation params, etc.
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    max_views INTEGER NOT NULL DEFAULT 1,
    view_count INTEGER NOT NULL DEFAULT 0,
    ip_lock INET,
    recipient_identifier TEXT, -- github:username, email, etc.
    destroyed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    accessed_at TIMESTAMP WITH TIME ZONE,
    destroyed_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT valid_max_views CHECK (max_views > 0),
    CONSTRAINT valid_view_count CHECK (view_count >= 0),
    CONSTRAINT view_count_not_exceed_max CHECK (view_count <= max_views)
);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    secret_id UUID NOT NULL REFERENCES secrets(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL, -- created, viewed, destroyed, expired
    ip_address INET,
    user_agent TEXT,
    metadata JSONB, -- additional context
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT valid_action CHECK (action IN ('created', 'viewed', 'destroyed', 'expired', 'failed_access'))
);

-- Indexes for performance
CREATE INDEX idx_secrets_expires_at ON secrets(expires_at) WHERE NOT destroyed;
CREATE INDEX idx_secrets_destroyed ON secrets(destroyed);
CREATE INDEX idx_secrets_created_at ON secrets(created_at);
CREATE INDEX idx_audit_logs_secret_id ON audit_logs(secret_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- Function to automatically mark secrets as destroyed when max views reached
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

-- Trigger to auto-destroy on max views
CREATE TRIGGER trigger_check_secret_destruction
    BEFORE UPDATE OF view_count ON secrets
    FOR EACH ROW
    EXECUTE FUNCTION check_secret_destruction();

-- Function to prevent updates to destroyed secrets
CREATE OR REPLACE FUNCTION prevent_destroyed_secret_access()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.destroyed = TRUE THEN
        RAISE EXCEPTION 'Cannot access destroyed secret';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to prevent access to destroyed secrets
CREATE TRIGGER trigger_prevent_destroyed_access
    BEFORE UPDATE OF view_count ON secrets
    FOR EACH ROW
    EXECUTE FUNCTION prevent_destroyed_secret_access();

-- Comments for documentation
COMMENT ON TABLE secrets IS 'Stores metadata for encrypted secrets (never plaintext)';
COMMENT ON COLUMN secrets.encrypted_blob_url IS 'S3/MinIO URL to encrypted blob';
COMMENT ON COLUMN secrets.encryption_metadata IS 'Client-side encryption parameters (algorithm, IV, etc.)';
COMMENT ON COLUMN secrets.ip_lock IS 'Optional IP address restriction for access';
COMMENT ON COLUMN secrets.recipient_identifier IS 'Optional recipient identifier for SSH key exchange';
COMMENT ON TABLE audit_logs IS 'Immutable audit trail for all secret operations';
