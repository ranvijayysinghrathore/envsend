-- Remove columns from secrets table
ALTER TABLE secrets
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS team_id,
    DROP COLUMN IF EXISTS organization_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_secrets_created_by;
DROP INDEX IF EXISTS idx_secrets_team_id;
DROP INDEX IF EXISTS idx_secrets_org_id;
DROP INDEX IF EXISTS idx_team_members_user_id;
DROP INDEX IF EXISTS idx_team_members_team_id;
DROP INDEX IF EXISTS idx_teams_org_id;
DROP INDEX IF EXISTS idx_org_members_user_id;
DROP INDEX IF EXISTS idx_org_members_org_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_organizations_slug;

-- Drop tables
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;
