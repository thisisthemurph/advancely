package auth

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthenticatedSessionUser struct {
	ID    uuid.UUID
	Email string
}

type AuthenticatedSession struct {
	LoggedIn bool

	SessionCookie
	AuthenticatedSessionUser
}

// Session returns the current user session stored in the context.
// If no session is present, a default session is returned.
func Session(c echo.Context) AuthenticatedSession {
	session, ok := c.Get(string(UserSessionContextKey)).(SessionCookie)
	if !ok {
		return AuthenticatedSession{}
	}

	user := AuthenticatedSessionUser{
		ID:    session.Sub,
		Email: session.Email,
	}

	return AuthenticatedSession{
		LoggedIn:                 true,
		SessionCookie:            session,
		AuthenticatedSessionUser: user,
	}
}
