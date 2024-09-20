package security

type Permission string

// Permissions group permissions

const (
	PermissionCreateUser Permission = "create-user"
	PermissionEditUser   Permission = "edit-user"
	PermissionDeleteUser Permission = "delete-user"
)

// User management group permissions

const (
	PermissionCreateRole     Permission = "create-role"
	PermissionEditRole       Permission = "edit-role"
	PermissionDeleteRole     Permission = "delete-role"
	PermissionAssignUserRole Permission = "assign-user-role"
	PermissionRemoveUserRole Permission = "remove-user-role"
)

// Organization management group settings

const (
	PermissionEditOrganizationSettings Permission = "edit-organization-settings"
)

func (p Permission) String() string {
	return string(p)
}

type Role string

const (
	RoleAdmin Role = "Admin"
)

func (r Role) String() string {
	return string(r)
}
