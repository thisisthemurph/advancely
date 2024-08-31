package routes

import (
	"advancely/internal/application"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{app.Config.ClientBaseURL},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	baseGroup := e.Group("/api/v1")
	for _, h := range getAPIHandlers(app) {
		h.MakeRoutes(baseGroup)
	}

	return e
}

func getAPIHandlers(app *application.App) []RouteMaker {
	return []RouteMaker{
		NewAuthHandler(app.Supabase, app.Store),
	}
}
