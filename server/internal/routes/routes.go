package routes

import (
	"net/http"

	"advancely/internal/application"
	mw "advancely/internal/middleware"
	"advancely/internal/validation"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RouteMaker is the interface used to specify handlers must implement the MakeRoutes method.
type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

func NewRouter(app *application.App) *echo.Echo {
	e := echo.New()
	e.Validator = validation.NewCustomValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{app.Config.ClientBaseURL},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	}))

	userMw := mw.NewUserMiddleware(app.Config.SessionSecret, app.Supabase, app.Store.UserStore, app.Logger)
	e.Use(userMw.WithUserInContext)

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
