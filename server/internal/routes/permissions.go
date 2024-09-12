package routes

import (
	"advancely/internal/application"
	"advancely/internal/auth"
	"advancely/internal/model"
	"advancely/internal/store"
	"advancely/internal/store/contract"
	"advancely/internal/validation"
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

func NewPermissionsHandler(
	s *store.Store,
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
	UserStore        contract.UserStore
	PermissionsStore contract.PermissionsStore
	Config           application.AppConfig
	Logger           *slog.Logger
}

func (h PermissionsHandler) MakeRoutes(e *echo.Group) {
	group := e.Group("/auth/permissions")
	group.GET("/role/:roleId", h.handleGetRoleWithPermissions())
	group.GET("/role", h.handleListRolesWithPermissions())
	group.POST("/role", h.handleCreateRole())
	group.PUT("/role/:roleId", h.handleUpdateRole())
	group.DELETE("/role/:roleId", h.handleDeleteRole())
}

func (h PermissionsHandler) handleGetRoleWithPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.Session(c)
		user, err := h.UserStore.User(session.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		roleId, err := strconv.Atoi(c.Param("roleId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		// Get the role
		role, err := h.PermissionsStore.Role(roleId, &user.CompanyID)
		if err != nil {
			if errors.Is(err, store.ErrRoleNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return c.JSON(http.StatusOK, role)
	}
}

func (h PermissionsHandler) handleListRolesWithPermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.Session(c)
		user, err := h.UserStore.User(session.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		roles, err := h.PermissionsStore.Roles(user.CompanyID)
		if err != nil {
			if errors.Is(err, store.ErrRoleNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
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
		session := auth.Session(c)
		user, err := h.UserStore.User(session.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// Get the role to create

		var request CreateRoleRequest
		if err := validation.BindAndValidate(c, &request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		role := model.CreateRole{
			CompanyID:   user.CompanyID,
			Name:        request.Name,
			Description: request.Description,
		}

		// Create the role

		createdRole, err := h.PermissionsStore.CreateRole(role)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		s := auth.Session(c)
		user, err := h.UserStore.User(s.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		roleId, err := strconv.Atoi(c.Param("roleId"))
		if err != nil {
			h.Logger.Error("error parsing roleId param", "error", err)
			return echo.NewHTTPError(400, "could not determine the role ID")
		}

		role, err := h.PermissionsStore.Role(roleId, &user.CompanyID)
		if err != nil {
			if errors.Is(err, store.ErrRoleNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(204)
	}
}

func (h PermissionsHandler) handleDeleteRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.Session(c)
		user, err := h.UserStore.User(session.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		roleIdParam := c.Param("roleId")
		if roleIdParam == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "role id is required")
		}

		roleId, err := strconv.Atoi(roleIdParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "role id is invalid")
		}

		if err := h.PermissionsStore.DeleteRole(roleId, user.CompanyID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	}
}
