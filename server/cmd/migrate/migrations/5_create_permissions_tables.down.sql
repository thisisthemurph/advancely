-- Drop triggers for the tables first
DROP TRIGGER IF EXISTS trg_set_updated_at_roles ON security.roles;
DROP TRIGGER IF EXISTS trg_set_updated_at_user_roles ON security.user_roles;
DROP TRIGGER IF EXISTS trg_set_updated_at_role_permissions ON security.role_permissions;

-- Drop the tables
DROP TABLE IF EXISTS security.role_permissions;
DROP TABLE IF EXISTS security.permissions;
DROP TABLE IF EXISTS security.permission_groups;
DROP TABLE IF EXISTS security.user_roles;
DROP TABLE IF EXISTS security.roles;

-- Drop the schema
DROP SCHEMA IF EXISTS security;
