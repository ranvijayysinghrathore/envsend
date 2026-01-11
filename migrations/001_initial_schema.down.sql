-- Drop triggers
DROP TRIGGER IF EXISTS trigger_prevent_destroyed_access ON secrets;
DROP TRIGGER IF EXISTS trigger_check_secret_destruction ON secrets;

-- Drop functions
DROP FUNCTION IF EXISTS prevent_destroyed_secret_access();
DROP FUNCTION IF EXISTS check_secret_destruction();

-- Drop indexes
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_timestamp;
DROP INDEX IF EXISTS idx_audit_logs_secret_id;
DROP INDEX IF EXISTS idx_secrets_created_at;
DROP INDEX IF EXISTS idx_secrets_destroyed;
DROP INDEX IF EXISTS idx_secrets_expires_at;

-- Drop tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS secrets;
