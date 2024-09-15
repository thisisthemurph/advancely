package model

import (
	"github.com/google/uuid"
)

type SystemRole string

const SystemRoleAdmin SystemRole = "Admin"

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

// UserRole is used to describe a role, including permissions on the role for a single user.
type UserRole struct {
	Name        string
	Permissions []string
}

// UserRoleCollection is a collection of UserRole objects for a specific user.
type UserRoleCollection struct {
	UserID uuid.UUID
	Roles  []UserRole
}

// HasPermission returns true if the permission is present on any role, otherwise false.
// The function always returns true if the user has the Admin role.
func (collection UserRoleCollection) HasPermission(name string) bool {
	for _, r := range collection.Roles {
		if r.Name == "Admin" {
			return true
		}
		for _, p := range r.Permissions {
			if p == name {
				return true
			}
		}
	}
	return false
}
