package middleware

import (
	"advancely/internal/auth"
	"advancely/internal/store/contract"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
)

var (
	ErrorNoAccessToken  = errors.New("no access token found")
	ErrorNoRefreshToken = errors.New("no refresh token found")
)

type UserMiddleware struct {
	*supabase.Client
	UserStore     contract.UserStore
	SessionSecret string
}

func NewUserMiddleware(sessionSecret string, client *supabase.Client, userStore contract.UserStore) *UserMiddleware {
	return &UserMiddleware{
		Client:        client,
		SessionSecret: sessionSecret,
		UserStore:     userStore,
	}
}

func (m *UserMiddleware) WithUserInContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		session, err := auth.GetSessionFromCookie(c, m.SessionSecret)
		if err != nil {
			return next(c)
		}

		if session.Expired() {
			refreshedAuthDetails, err := m.refreshSupabaseUser(ctx, session)
			if err != nil {
				return next(c)
			}

			session = auth.NewSessionCookie(refreshedAuthDetails)
			if err := session.SetCookie(c, m.SessionSecret); err != nil {
				return next(c)
			}
		}

		session.SaveInContext(c)
		return next(c)
	}
}

// refreshSupabaseUser refreshes the user wit the refresh token, returning the new supabase.AuthenticatedDetails.
// Returns an error if tokens are missing in the session or if there is an error refreshing via Supabase.
func (m *UserMiddleware) refreshSupabaseUser(ctx context.Context, session *auth.SessionCookie) (*supabase.AuthenticatedDetails, error) {
	if len(session.AccessToken) == 0 {
		return nil, ErrorNoAccessToken
	}
	if len(session.RefreshToken) == 0 {
		return nil, ErrorNoRefreshToken
	}

	refreshedAuth, err := m.Auth.RefreshUser(ctx, session.AccessToken, session.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh user: %w", err)
	}
	return refreshedAuth, nil
}
