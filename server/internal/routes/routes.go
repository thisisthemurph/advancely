package routes

import (
	"advancely/internal/application"
	mw "advancely/internal/middleware"
	"advancely/internal/validation"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RouteMaker is the interface used to specify handlers must implement the MakeRoutes method.
type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

func NewRouter(app *application.App) *echo.Echo {
	e := echo.New()

	e.Validator = validation.NewCustomValidator()
	setUpMiddlewares(e, app)

	baseGroup := e.Group("/api/v1")
	for _, h := range buildAPIHandlers(app) {
		h.MakeRoutes(baseGroup)
	}

	return e
}

func buildAPIHandlers(app *application.App) []RouteMaker {
	return []RouteMaker{
		NewAuthHandler(app.Supabase, app.Store, app.Config, app.Logger),
		NewPermissionsHandler(app.Store, app.Config, app.Logger),
	}
}

func setUpMiddlewares(e *echo.Echo, app *application.App) {
	// In production, the client is served from the api.
	// In development, the client must be served separately.
	if app.Config.Environment.IsProduction() {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Skipper:    nil,
			Root:       "../client/dist",
			Index:      "index.html",
			HTML5:      true,
			Browse:     false,
			IgnoreBase: false,
			Filesystem: nil,
		}))
	}

	// Show additional HTTP request logging in development only.
	if app.Config.Environment.IsDevelopment() {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.Recover())

	allowOrigins := []string{app.Config.ClientBaseURL}
	if app.Config.Environment.IsDevelopment() {
		allowOrigins = append(allowOrigins, "http://localhost:*")
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	}))

	userMw := mw.NewUserMiddleware(app.Config, app.Supabase, app.Store.UserStore, app.Logger)
	e.Use(userMw.WithUserInContext)
}
