package routes

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"advancely/internal/application"
	"advancely/internal/auth"
	"advancely/internal/model"
	"advancely/internal/model/security"
	"advancely/internal/store"
	"advancely/internal/validation"
	"advancely/pkg/errs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewPermissionsHandler(
	s *store.PostgresStore,
	config application.AppConfig,
	logger *slog.Logger,
) PermissionsHandler {
	return PermissionsHandler{
		UserStore:        s.UserStore,
		PermissionsStore: s.PermissionsStore,
		Config:           config,
		Logger:           logger,
	}
}

type PermissionsHandler struct {
	UserStore        store.UserStore
	PermissionsStore store.PermissionsStore
	Config           application.AppConfig
	Logger           *slog.Logger
}

// EnsurePermission returns an echo.HTTPError if the permission does not exist for the
// current authenticated session user.
func (h PermissionsHandler) EnsurePermission(c echo.Context, permission security.Permission) *echo.HTTPError {
	return EnsurePermission(c, h.PermissionsStore, permission)
}

func (h PermissionsHandler) MakeRoutes(e *echo.Group) {
	group := e.Group("/auth/permissions")

	roleGroup := group.Group("/role")
	roleGroup.GET("/:roleId", h.handleGetRoleWithPermissions())
	roleGroup.GET("", h.handleListRolesWithPermissions())
	roleGroup.POST("", h.handleCreateRole())
	roleGroup.PUT("/:roleId", h.handleUpdateRole())
	roleGroup.DELETE("/:roleId", h.handleDeleteRole())

	permissionGroup := group.Group("/role/:roleId/permission")
	permissionGroup.POST("/:permissionId", h.handleAssignPermissionToRole())
	permissionGroup.DELETE("/:permissionId", h.handleRemovePermissionFromRole())

	userGroup := group.Group("/role/:roleId/user")
	userGroup.POST("/:userId", h.handleAssignRoleToUser())
	userGroup.DELETE("/:userId", h.handleRemoveRoleFromUser())
}

func (h PermissionsHandler) handleGetRoleWithPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)

		roleId, err := strconv.Atoi(c.Param("roleId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Role ID could not be parsed")
		}

		role, err := h.PermissionsStore.Role(roleId, &session.Company.ID)
		if err != nil {
			if errors.Is(err, store.ErrRoleNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, role)
	}
}

func (h PermissionsHandler) handleListRolesWithPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		roles, err := h.PermissionsStore.Roles(session.Company.ID)
		if err != nil {
			h.Logger.Error("failed to list roles", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, roles)
	}
}

func (h PermissionsHandler) handleCreateRole() echo.HandlerFunc {
	type CreateRoleRequest struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		if err := h.EnsurePermission(c, security.PermissionCreateRole); err != nil {
			return err
		}

		// Get the role to create

		var request CreateRoleRequest
		if err := validation.BindAndValidate(c, &request); err != nil {
			return err
		}

		role := model.CreateRole{
			CompanyID:   session.Company.ID,
			Name:        request.Name,
			Description: request.Description,
		}

		// Create the role

		createdRole, err := h.PermissionsStore.CreateRole(role)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusCreated, createdRole)
	}
}

func (h PermissionsHandler) handleUpdateRole() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		if err := h.EnsurePermission(c, security.PermissionEditRole); err != nil {
			return err
		}

		roleId, err := strconv.Atoi(c.Param("roleId"))
		if err != nil {
			h.Logger.Error("error parsing roleId param", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "could not determine the role ID")
		}

		role, err := h.PermissionsStore.Role(roleId, &session.Company.ID)
		if err != nil {
			if errors.Is(err, store.ErrRoleNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if role.IsSystemRole {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "cannot update system role")
		}

		var request Request
		if err := validation.BindAndValidate(c, &request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		update := model.Role{
			ID:           role.ID,
			CompanyID:    role.CompanyID,
			Name:         request.Name,
			Description:  request.Description,
			IsSystemRole: role.IsSystemRole,
		}

		if err := h.PermissionsStore.UpdateRole(&update); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (h PermissionsHandler) handleDeleteRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		if err := h.EnsurePermission(c, security.PermissionDeleteRole); err != nil {
			return err
		}

		roleIdParam := c.Param("roleId")
		if roleIdParam == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "role id is required")
		}

		roleId, err := strconv.Atoi(roleIdParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "role id is invalid")
		}

		if err := h.PermissionsStore.DeleteRole(roleId, session.Company.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (h PermissionsHandler) handleAssignPermissionToRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		if err := h.EnsurePermission(c, security.PermissionEditRole); err != nil {
			return err
		}

		roleID, ridErr := strconv.Atoi(c.Param("roleId"))
		permissionId, pidErr := strconv.Atoi(c.Param("permissionId"))
		if ridErr != nil || pidErr != nil {
			return echo.NewHTTPError(400, "role or permission ID not valid")
		}

		if err := h.PermissionsStore.AssignPermissionToRole(roleID, permissionId, session.Company.ID); err != nil {
			if errors.Is(err, store.ErrCannotUpdateSystemRole) {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
			}
			if errs.IsOne(err, store.ErrRoleNotFound, store.ErrPermissionNotFount) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			h.Logger.Error("error assigning permission to role", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusCreated)
	}
}

func (h PermissionsHandler) handleRemovePermissionFromRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		if err := h.EnsurePermission(c, security.PermissionEditRole); err != nil {
			return err
		}

		roleID, ridErr := strconv.Atoi(c.Param("roleId"))
		permissionId, pidErr := strconv.Atoi(c.Param("permissionId"))
		if ridErr != nil || pidErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "role or permission ID not valid")
		}

		if err := h.PermissionsStore.RemovePermissionFromRole(roleID, permissionId, session.Company.ID); err != nil {
			if errors.Is(err, store.ErrRoleNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			if errs.IsOne(err, store.ErrCannotUpdateSystemRole, store.ErrCannotUpdateSystemRole) {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
			}
			h.Logger.Error("error removing permission from role", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (h PermissionsHandler) handleAssignRoleToUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		if err := h.EnsurePermission(c, security.PermissionAssignUserRole); err != nil {
			return err
		}

		roleID, err := strconv.Atoi(c.Param("roleId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "permission ID not valid")
		}

		userID, err := uuid.Parse(c.Param("userId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "user ID not valid")
		}

		if err := h.PermissionsStore.AssignRoleToUser(roleID, userID, session.Company.ID); err != nil {
			if errors.Is(err, store.ErrRoleNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			h.Logger.Error("error assigning role to user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusCreated)
	}
}

func (h PermissionsHandler) handleRemoveRoleFromUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		if err := h.EnsurePermission(c, security.PermissionRemoveUserRole); err != nil {
			return err
		}

		roleID, err := strconv.Atoi(c.Param("roleId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "permission ID not valid")
		}

		userID, err := uuid.Parse(c.Param("userId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "user ID not valid")
		}

		if err := h.PermissionsStore.RemoveRoleFromUser(roleID, userID, session.Company.ID); err != nil {
			h.Logger.Error("error removing role from user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
