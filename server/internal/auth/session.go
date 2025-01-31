package auth

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthenticatedSessionUser struct {
	ID    uuid.UUID
	Email string
}

type AuthenticatedSessionCompany struct {
	ID uuid.UUID
}

type AuthenticatedSession struct {
	LoggedIn bool

	SessionCookie
	AuthenticatedSessionUser
	AuthenticatedSessionCompany
}

// CurrentUser returns the current user session stored in the echo.Context.
// If no session is present or the session could not be parsed, a default session is returned.
func CurrentUser(c echo.Context) AuthenticatedSession {
	session, ok := c.Get(string(UserSessionContextKey)).(SessionCookie)
	if !ok {
		return AuthenticatedSession{}
	}

	user := AuthenticatedSessionUser{
		ID:    session.User.ID,
		Email: session.User.Email,
	}

	company := AuthenticatedSessionCompany{
		ID: session.Company.ID,
	}

	return AuthenticatedSession{
		LoggedIn:                    true,
		SessionCookie:               session,
		AuthenticatedSessionUser:    user,
		AuthenticatedSessionCompany: company,
	}
}
