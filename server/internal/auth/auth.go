package auth

import (
	"context"
	"github.com/labstack/echo/v4"
)

type UserContextKey string

const UserSessionContextKey UserContextKey = "session"

// saveInContext persists the given value within both the Echo context and the Echo request context.
func saveInContext(c echo.Context, key UserContextKey, value any) {
	// Persist session in the echo context
	c.Set(string(key), value)
	// Persist the session in the request context
	c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), key, value)))
}
