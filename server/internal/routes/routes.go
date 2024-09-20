package routes

import (
	"net/http"

	"advancely/internal/application"
	"advancely/internal/auth"
	mw "advancely/internal/middleware"
	"advancely/internal/model/security"
	"advancely/internal/store"
	"advancely/internal/validation"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RouteMaker is the interface used to specify handlers must implement the MakeRoutes method.
type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

type Router struct {
	*echo.Echo
	RoleFetcher store.RoleFetcher
}

func NewRouter(app *application.App) *Router {
	r := &Router{
		Echo:        echo.New(),
		RoleFetcher: app.Store.PermissionsStore,
	}

	r.Validator = validation.NewCustomValidator()
	r.configureMiddleware(app)

	baseGroup := r.Group("/api/v1")
	for _, h := range r.getRouteHandlers(app) {
		h.MakeRoutes(baseGroup)
	}

	return r
}

type EnsurePermissionFn = func(c echo.Context, permission security.Permission) *echo.HTTPError

func (r *Router) getRouteHandlers(app *application.App) []RouteMaker {
	ensurePermissionFn := func(c echo.Context, permission security.Permission) *echo.HTTPError {
		session := auth.CurrentUser(c)
		roles, err := r.RoleFetcher.UserRoles(session.User.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if !roles.HasPermission(permission) {
			return echo.NewHTTPError(http.StatusForbidden)
		}
		return nil
	}

	return []RouteMaker{
		NewAuthHandler(app.Supabase, app.Store, app.Config, app.Logger),
		NewPermissionsHandler(app.Store, app.Config, app.Logger, ensurePermissionFn),
		NewCompaniesHandler(app.Store, app.Logger, ensurePermissionFn),
	}
}

func (r *Router) configureMiddleware(app *application.App) {
	// In production, the client is served from the api.
	// In development, the client must be served separately.
	if app.Config.Environment.IsProduction() {
		r.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Skipper:    nil,
			Root:       "../client/dist",
			Index:      "index.html",
			HTML5:      true,
			Browse:     false,
			IgnoreBase: false,
			Filesystem: nil,
		}))
	}

	r.Use(middleware.Recover())

	allowOrigins := []string{app.Config.ClientBaseURL}
	if app.Config.Environment.IsDevelopment() {
		allowOrigins = append(allowOrigins, "http://localhost:*")
	}
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	}))

	userMw := mw.NewUserMiddleware(app.Config, app.Supabase, app.Store.UserStore, app.Logger)
	r.Use(userMw.WithUserInContext)
}
