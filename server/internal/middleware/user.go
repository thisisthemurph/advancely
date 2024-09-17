package middleware

import (
	"advancely/internal/application"
	"advancely/internal/auth"
	"advancely/internal/store"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"log/slog"
)

var (
	ErrorNoAccessToken  = errors.New("no access token found")
	ErrorNoRefreshToken = errors.New("no refresh token found")
)

type UserMiddleware struct {
	UserStore store.UserStore
	Config    application.AppConfig
	Logger    *slog.Logger
	Supabase  *supabase.Client
}

func NewUserMiddleware(
	config application.AppConfig,
	supabaseClient *supabase.Client,
	userStore store.UserStore,
	logger *slog.Logger,
) *UserMiddleware {
	return &UserMiddleware{
		Supabase:  supabaseClient,
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
			tokenResp, err := m.refreshSupabaseUser(ctx, session)
			if err != nil {
				logger.Debug("failed to refresh supabase user", "error", err)
				return next(c)
			}

			session = auth.NewSessionCookie(tokenResp.Session)
			if err := session.SetCookie(c, m.Config.SessionSecret, m.Config.Environment); err != nil {
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
func (m *UserMiddleware) refreshSupabaseUser(ctx context.Context, session *auth.SessionCookie) (*types.TokenResponse, error) {
	if len(session.AccessToken) == 0 {
		return nil, ErrorNoAccessToken
	}
	if len(session.RefreshToken) == 0 {
		return nil, ErrorNoRefreshToken
	}

	resp, err := m.Supabase.Auth.WithToken(session.AccessToken).RefreshToken(session.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh user: %w", err)
	}
	return resp, nil
}
