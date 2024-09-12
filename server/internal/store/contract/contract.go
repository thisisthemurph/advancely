package contract

import (
	"advancely/internal/model"
	"github.com/google/uuid"
)

type Store interface {
	UserStore
	CompanyStore
	PermissionsStore
}

type UserStore interface {
	// User returns the user associated with the given id.
	User(id uuid.UUID) (model.UserProfile, error)
	// BaseUserByEmail returns the auth.users user associated with the given email.
	BaseUserByEmail(email string) (model.User, error)
	// Users returns a slice of all users.
	Users(companyID uuid.UUID) ([]model.UserProfile, error)
	// CreateProfile creates a record in the profiles table.
	CreateProfile(user *model.UserProfile) error
	UpdateUser(user *model.UserProfile) error
	DeleteUser(id uuid.UUID) error
}

type CompanyStore interface {
	// Company returns the company associated with the given id.
	Company(id uuid.UUID) (model.Company, error)
	// CompanyByCreator returns the company created by the given creator user ID.
	CompanyByCreator(creatorID uuid.UUID) (model.Company, error)
	// Companies returns a slice of all companies.
	Companies() ([]model.Company, error)
	CreateCompany(c *model.Company) error
	UpdateCompany(c *model.Company) error
	DeleteCompany(id uuid.UUID) error
}

type PermissionsStore interface {
	// Role returns the role associated with the given ID
	// Passing nil for the companyID will allow searching for matching system roles
	Role(id int, companyID *uuid.UUID) (model.RoleWithPermissions, error)
	// Roles returns all roles (including system) for the given companyID
	Roles(companyID uuid.UUID) ([]model.RoleWithPermissions, error)
	CreateRole(r model.CreateRole) (model.Role, error)
	UpdateRole(r *model.Role) error
	DeleteRole(id int, companyID uuid.UUID) error
	// AssignPermissionToRole associates a given permission with the given role.
	// Users cannot associate any permissions with system roles.
	AssignPermissionToRole(roleID, permissionID int, companyID uuid.UUID) error
	// RemovePermissionFromRole removes the role - permission association.
	// Users cannot remove a permission from a system role.
	RemovePermissionFromRole(roleID, permissionID int, companyID uuid.UUID) error
	// AssignRoleToUser assigns a role to a given user.
	// A success is returned if the role already exists for the user.
	AssignRoleToUser(roleID int, userID, companyID uuid.UUID) error
	// AssignSystemRoleToUser assigns the specified system role to a given user.
	// A success is returned if the role already exists for the user.
	AssignSystemRoleToUser(role model.SystemRole, userID, companyID uuid.UUID) error
	// RemoveRoleFromUser disassociates the given role from the user.
	RemoveRoleFromUser(roleID int, userID, companyID uuid.UUID) error
}
