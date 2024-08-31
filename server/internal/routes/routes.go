package routes

import (
	"advancely/internal/application"
	"advancely/internal/validation"
	"github.com/labstack/echo/v4"
)

// RouteMaker is the interface used to specify handlers must implement the MakeRoutes method.
type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

func NewRouter(app *application.App) *echo.Echo {
	e := echo.New()
	e.Validator = validation.NewCustomValidator()

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
