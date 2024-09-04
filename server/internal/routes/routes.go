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
		NewAuthHandler(app.Supabase, app.Store, app.Config.SessionSecret, app.Logger),
	}
}

func setUpMiddlewares(e *echo.Echo, app *application.App) {
	if app.IsDevelopment() {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{app.Config.ClientBaseURL},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	}))

	userMw := mw.NewUserMiddleware(app.Config.SessionSecret, app.Supabase, app.Store.UserStore, app.Logger)
	e.Use(userMw.WithUserInContext)
}
