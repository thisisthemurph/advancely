package routes

import (
	"advancely/internal/application"
	"advancely/internal/validation"
	"github.com/labstack/echo/v4"
)

func NewRouter(app *application.App) *echo.Echo {
	e := echo.New()
	e.Validator = validation.NewCustomValidator()

	authHandler := authHandler{
		Supabase:     app.Supabase,
		UserStore:    app.Store.UserStore,
		CompanyStore: app.Store.CompanyStore,
	}
	newAuthRouter(e, authHandler)

	return e
}
