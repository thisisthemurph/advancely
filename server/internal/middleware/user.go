package middleware

import (
	"advancely/internal/application"
	"advancely/internal/auth"
	"advancely/internal/store/contract"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
	"log/slog"
)

var (
	ErrorNoAccessToken  = errors.New("no access token found")
	ErrorNoRefreshToken = errors.New("no refresh token found")
)

type UserMiddleware struct {
	*supabase.Client
	UserStore contract.UserStore
	Config    application.AppConfig
	Logger    *slog.Logger
}

func NewUserMiddleware(
	config application.AppConfig,
	client *supabase.Client,
	userStore contract.UserStore,
	logger *slog.Logger,
) *UserMiddleware {
	return &UserMiddleware{
		Client:    client,
		Config:    config,
		UserStore: userStore,
		Logger:    logger,
	}
}

func (m *UserMiddleware) WithUserInContext(next echo.HandlerFunc) echo.HandlerFunc {
	logger := m.Logger.With("mw", "WithUserInContext")
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		session, err := auth.GetSessionFromCookie(c, m.Config.SessionSecret)
		if err != nil {
			logger.Debug("failed to get session from cookie", "error", err)
			return next(c)
		}

		if session.Expired() {
			refreshedAuthDetails, err := m.refreshSupabaseUser(ctx, session)
			if err != nil {
				logger.Debug("failed to refresh supabase user", "error", err)
				return next(c)
			}

			session = auth.NewSessionCookie(refreshedAuthDetails)
			if err := session.SetCookie(c, m.Config.SessionSecret, m.Config.IsDevelopment); err != nil {
				logger.Error("failed to set cookie", "error", err)
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
