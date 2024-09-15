package routes

import (
	"net/http"

	"advancely/internal/auth"
	"advancely/internal/model/security"
	"advancely/internal/store"
	"github.com/labstack/echo/v4"
)

// EnsurePermission returns an echo.HTTPError if the user session within the given echo.Context
// does not have the specified permission.
func EnsurePermission(c echo.Context, fetcher store.RoleFetcher, permission security.Permission) *echo.HTTPError {
	session := auth.CurrentUser(c)
	roles, err := fetcher.UserRoles(session.User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !roles.HasPermission(permission) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	return nil
}
