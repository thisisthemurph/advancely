package model

import "github.com/google/uuid"

// Role represents the security.roles table.
type Role struct {
	ID           int        `db:"id" json:"id"`
	CompanyID    *uuid.UUID `db:"company_id" json:"companyId,omitempty"`
	Name         string     `db:"name" json:"name"`
	Description  string     `db:"description" json:"description"`
	IsSystemRole bool       `db:"is_system_role" json:"system"`
}

// CreateRole is a model used to create roles in the store.
type CreateRole struct {
	CompanyID   uuid.UUID `json:"companyId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// PermissionGroup represents the security.permission_groups table.
type PermissionGroup struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

// Permission represents the security.permissions table.
type Permission struct {
	ID          int             `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description"`
	Group       PermissionGroup `json:"group"`
}

// RoleWithPermissions represents the join between the security.roles and security.permissions table.
type RoleWithPermissions struct {
	Role
	Permissions []Permission `json:"permissions"`
}
