package security

import "github.com/google/uuid"

// UserRole is used to describe a role, including permissions on the role for a single user.
type UserRole struct {
	Role        Role
	Permissions []Permission
}

// UserRoleCollection is a collection of UserRole objects for a specific user.
type UserRoleCollection struct {
	UserID uuid.UUID
	Roles  []UserRole
}

// HasPermission returns true if the permission is present on any role, otherwise false.
// The function always returns true if the user has the RoleAdmin role.
func (collection UserRoleCollection) HasPermission(name Permission) bool {
	for _, r := range collection.Roles {
		if r.Role == RoleAdmin {
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
