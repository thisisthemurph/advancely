package store

import (
	"advancely/internal/model"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRoleNotFound           = errors.New("role not found")
	ErrPermissionNotFount     = errors.New("permission not found")
	ErrCannotDeleteSystemRole = errors.New("cannot delete system role")
	ErrCannotUpdateSystemRole = errors.New("cannot update system role")
)

func NewPermissionsStore(db *sqlx.DB) *PermissionsStore {
	return &PermissionsStore{
		DB: db,
	}
}

type PermissionsStore struct {
	*sqlx.DB
}

type rolePermission struct {
	RoleID         int            `db:"id"`
	CompanyID      *uuid.UUID     `db:"company_id"`
	RoleName       string         `db:"name"`
	RoleDesc       string         `db:"description"`
	IsSystemRole   bool           `db:"is_system_role"`
	PermissionID   sql.NullInt64  `db:"permission_id"`
	PermissionName sql.NullString `db:"permission_name"`
	PermissionDesc sql.NullString `db:"permission_description"`
}

func (s *PermissionsStore) Role(id int, companyID *uuid.UUID) (model.RoleWithPermissions, error) {
	stmt := `
		select
		  r.id, r.company_id, r.name, r.description, r.is_system_role,
		  p.id as permission_id, p.name as permission_name, p.description as permission_description
		from security.roles r
		  left join security.role_permissions rp on r.id = rp.role_id
		  left join security.permissions p on rp.permission_id = p.id
		where r.id = $1
		  and (r.company_id = $2 or r.is_system_role = true);`

	var rpList []rolePermission
	if err := s.Select(&rpList, stmt, id, companyID); err != nil {
		return model.RoleWithPermissions{}, err
	}

	if rpList == nil {
		return model.RoleWithPermissions{}, ErrRoleNotFound
	}

	if companyID != nil && *companyID == uuid.Nil {
		// testing
		companyID = nil
	}

	role := model.RoleWithPermissions{
		Role: model.Role{
			ID:           rpList[0].RoleID,
			CompanyID:    rpList[0].CompanyID,
			Name:         rpList[0].RoleName,
			Description:  rpList[0].RoleDesc,
			IsSystemRole: rpList[0].IsSystemRole,
		},
		Permissions: []model.Permission{},
	}

	for _, rp := range rpList {
		if !rp.PermissionID.Valid {
			continue
		}
		permission := model.Permission{
			ID:          int(rp.PermissionID.Int64),
			Name:        rp.PermissionName.String,
			Description: rp.PermissionDesc.String,
		}
		role.Permissions = append(role.Permissions, permission)
	}

	return role, nil
}

func (s *PermissionsStore) Roles(companyID uuid.UUID) ([]model.RoleWithPermissions, error) {
	stmt := `
		select
		  r.id, r.company_id, r.name, r.description, r.is_system_role,
		  p.id as permission_id, p.name as permission_name, p.description as permission_description
		from security.roles r
		  left join security.role_permissions rp on r.id = rp.role_id
		  left join security.permissions p on rp.permission_id = p.id
		where company_id = $1 or is_system_role = true
		order by r.id, p.id;`

	var rpList []rolePermission
	if err := s.Select(&rpList, stmt, companyID); err != nil {
		return nil, fmt.Errorf("failed to list roles for company ID %v: %w", companyID, err)
	}

	if rpList == nil {
		return []model.RoleWithPermissions{}, nil
	}

	var roles []model.RoleWithPermissions
	roleMap := make(map[int]*model.RoleWithPermissions)

	for _, rp := range rpList {
		role, exists := roleMap[rp.RoleID]
		if !exists {
			role = &model.RoleWithPermissions{
				Role: model.Role{
					ID:           rp.RoleID,
					CompanyID:    rp.CompanyID,
					Name:         rp.RoleName,
					Description:  rp.RoleDesc,
					IsSystemRole: rp.IsSystemRole,
				},
				Permissions: []model.Permission{},
			}
			roleMap[rp.RoleID] = role
			roles = append(roles, *role)
		}

		if !rp.PermissionID.Valid {
			continue
		}

		currentIndex := len(roles) - 1
		roles[currentIndex].Permissions = append(roles[currentIndex].Permissions, model.Permission{
			ID:          int(rp.PermissionID.Int64),
			Name:        rp.PermissionName.String,
			Description: rp.PermissionDesc.String,
		})
	}
	return roles, nil
}

func (s *PermissionsStore) CreateRole(r model.CreateRole) (model.Role, error) {
	stmt := `
		insert into security.roles (company_id, name, description)
		values ($1, $2, $3)
		returning id, company_id, name, description, is_system_role;`

	var createdRole model.Role
	if err := s.Get(&createdRole, stmt, r.CompanyID, r.Name, r.Description); err != nil {
		return model.Role{}, fmt.Errorf("failed to create role: %w", err)
	}

	return createdRole, nil
}

func (s *PermissionsStore) UpdateRole(r *model.Role) error {
	role, err := s.Role(r.ID, r.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to find role with ID %d: %w", r.ID, err)
	}
	if role.IsSystemRole {
		return ErrCannotUpdateSystemRole
	}

	stmt := `
		update security.roles
		set name = $1, description = $2
		where id = $3
		  and company_id = $4
		  and is_system_role = false -- prevent updating of system roles
		returning id, company_id, name, description, is_system_role;`

	if err := s.Get(r, stmt, r.Name, r.Description, r.ID, r.CompanyID); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (s *PermissionsStore) DeleteRole(id int, companyID uuid.UUID) error {
	role, err := s.Role(id, &companyID)
	if err != nil {
		return fmt.Errorf("failed to find role with ID %d: %w", id, err)
	}
	if role.IsSystemRole {
		return ErrCannotDeleteSystemRole
	}

	stmt := "delete from security.roles where id = $1 and company_id = $2;"
	if _, err := s.Exec(stmt, id, companyID); err != nil {
		return err
	}
	return nil
}

func (s *PermissionsStore) Permission(id int) (model.Permission, error) {
	stmt := `
		select p.id, p.name, p.description,
		       g.id as group_id,
		       g.name as group_name,
		       g.description as group_description
		from security.permissions p
		join security.permission_groups g
			on p.group_id = g.id
		where p.id = $1;`

	var result struct {
		ID          int    `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		GroupID     int    `db:"group_id"`
		GroupName   string `db:"group_name"`
		GroupDesc   string `db:"group_description"`
	}

	if err := s.Get(&result, stmt, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Permission{}, ErrPermissionNotFount
		}
		return model.Permission{}, err
	}

	permission := model.Permission{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Group: model.PermissionGroup{
			ID:          result.GroupID,
			Name:        result.GroupName,
			Description: result.GroupDesc,
		},
	}
	return permission, nil
}

func (s *PermissionsStore) AssignPermissionToRole(roleID, permissionID int, companyID uuid.UUID) error {
	role, err := s.Role(roleID, &companyID)
	if err != nil {
		return err
	}
	if role.IsSystemRole {
		return ErrCannotUpdateSystemRole
	}

	if _, err = s.Permission(permissionID); err != nil {
		return err
	}

	stmt := "insert into security.role_permissions (role_id, permission_id) values ($1, $2)"
	if _, err := s.Exec(stmt, roleID, permissionID); err != nil {
		// Check for unique_violation error, the relationship already exists.
		if pge := checkPgErr(err); errors.Is(pge, PgErrCodeUniqueViolation) {
			return nil
		}
		return fmt.Errorf("failed to insert role permission: %w", err)
	}

	return nil
}

func (s *PermissionsStore) RemovePermissionFromRole(roleID, permissionID int, companyID uuid.UUID) error {
	role, err := s.Role(roleID, &companyID)
	if err != nil {
		return err
	}
	if role.IsSystemRole {
		return ErrCannotUpdateSystemRole
	}

	stmt := "delete from security.role_permissions where role_id = $1 and permission_id = $2"
	if _, err := s.Exec(stmt, roleID, permissionID); err != nil {
		return fmt.Errorf("failed to delete role permission: %w", err)
	}
	return nil
}

func (s *PermissionsStore) AssignRoleToUser(roleID int, userID, companyID uuid.UUID) error {
	_, err := s.Role(roleID, &companyID)
	if err != nil {
		return err
	}

	stmt := "insert into security.user_roles (user_id, role_id) values ($1, $2);"
	if _, err := s.Exec(stmt, userID, roleID); err != nil {
		// Check for postgres unique_violation, relationship already exists
		if pgErr := checkPgErr(err); errors.Is(pgErr, PgErrCodeUniqueViolation) {
			return nil
		}
		return fmt.Errorf("failed to insert user role: %w", err)
	}
	return nil
}

func (s *PermissionsStore) AssignSystemRoleToUser(role model.SystemRole, userID, companyID uuid.UUID) error {
	var roleId int
	stmt := "select id from security.roles where name = $1 and is_system_role = true;"
	if err := s.Get(&roleId, stmt, role); err != nil {
		return err
	}
	return s.AssignRoleToUser(roleId, userID, companyID)
}

func (s *PermissionsStore) RemoveRoleFromUser(roleID int, userID, companyID uuid.UUID) error {
	stmt := "delete from security.user_roles where user_id = $1 and role_id = $2;"
	if _, err := s.Exec(stmt, userID, roleID); err != nil {
		return fmt.Errorf("failed to delete user role: %w", err)
	}
	return nil
}
