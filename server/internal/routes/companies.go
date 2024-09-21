package routes

import (
	"advancely/internal/auth"
	"advancely/internal/model/security"
	"advancely/internal/store"
	"advancely/internal/validation"
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func NewCompaniesHandler(
	s *store.PostgresStore,
	logger *slog.Logger,
	ensurePermissionFn EnsurePermissionFn) CompaniesHandler {
	return CompaniesHandler{
		CompanySettingsStore: s.CompanySettingsStore,
		Logger:               logger,
		EnsurePermission:     ensurePermissionFn,
	}
}

type CompaniesHandler struct {
	CompanySettingsStore store.CompanySettingsStore
	Logger               *slog.Logger
	EnsurePermission     EnsurePermissionFn
}

func (h CompaniesHandler) MakeRoutes(e *echo.Group) {
	group := e.Group("/company/settings")
	group.POST("/domain", h.HandleAddAllowedDomain())
}

type AddAllowedDomainRequest struct {
	Domain              string `json:"domain"`
	AllowUnknownDomains bool   `json:"allowUnknownDomains"`
}

func (h CompaniesHandler) HandleAddAllowedDomain() echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := h.EnsurePermission(c, security.PermissionEditOrganizationSettings); err != nil {
			return err
		}

		ctx := c.Request().Context()
		user := auth.CurrentUser(c)

		var req AddAllowedDomainRequest
		if err := validation.BindAndValidate(c, &req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validation.ValidateDomain(req.Domain); err != nil {
			if !(errors.Is(err, validation.ErrUnknownDomain) && req.AllowUnknownDomains) {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}

		if err := h.CompanySettingsStore.AddAllowedEmailDomain(ctx, user.Company.ID, req.Domain); err != nil {
			if errors.Is(err, store.ErrDomainAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusCreated)
	}
}
