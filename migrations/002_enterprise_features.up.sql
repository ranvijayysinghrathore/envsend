-- Create organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(255) NOT NULL UNIQUE,
    plan VARCHAR(50) NOT NULL DEFAULT 'free', -- free, team, enterprise
    max_members INTEGER NOT NULL DEFAULT 5,
    max_secrets_per_month INTEGER NOT NULL DEFAULT 1000,
    custom_retention_days INTEGER,
    sso_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    sso_provider VARCHAR(50),
    sso_config JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_plan CHECK (plan IN ('free', 'team', 'enterprise'))
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT, -- nullable for SSO users
    ssh_public_key TEXT,
    github_username VARCHAR(255),
    gitlab_username VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Create organization_members table
CREATE TABLE IF NOT EXISTS organization_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member', -- owner, admin, member, viewer
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_role CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    CONSTRAINT unique_org_user UNIQUE (organization_id, user_id)
);

-- Create teams table (within organizations)
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unique_org_team_name UNIQUE (organization_id, name)
);

-- Create team_members table
CREATE TABLE IF NOT EXISTS team_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member', -- lead, member
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_team_role CHECK (role IN ('lead', 'member')),
    CONSTRAINT unique_team_user UNIQUE (team_id, user_id)
);

-- Add organization and team references to secrets table
ALTER TABLE secrets
    ADD COLUMN organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    ADD COLUMN team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    ADD COLUMN created_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Indexes for enterprise features
CREATE INDEX idx_organizations_slug ON organizations(slug);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_org_members_org_id ON organization_members(organization_id);
CREATE INDEX idx_org_members_user_id ON organization_members(user_id);
CREATE INDEX idx_teams_org_id ON teams(organization_id);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_secrets_org_id ON secrets(organization_id);
CREATE INDEX idx_secrets_team_id ON secrets(team_id);
CREATE INDEX idx_secrets_created_by ON secrets(created_by);

-- Comments
COMMENT ON TABLE organizations IS 'Multi-tenant organizations for enterprise features';
COMMENT ON TABLE users IS 'User accounts with optional SSO integration';
COMMENT ON TABLE organization_members IS 'Organization membership with RBAC';
COMMENT ON TABLE teams IS 'Teams within organizations for access control';
COMMENT ON TABLE team_members IS 'Team membership';
