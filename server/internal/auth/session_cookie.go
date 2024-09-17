package auth

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"time"

	"advancely/internal/application"
	"advancely/internal/model"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/supabase-community/gotrue-go/types"
)

const (
	SessionCookieStoreName       = "session"
	SessionCookieSessionValueKey = "session"
)

var ErrorCookieNotFound = errors.New("cookie not found")

func init() {
	// Register the session cookie with gob so it can be used in the session store.
	gob.Register(&SessionCookie{})
}

func NewSessionCookie(session types.Session) *SessionCookie {
	return &SessionCookie{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		TokenType:    session.TokenType,
		ExpiresIn:    session.ExpiresIn,
		ExpiresAt:    session.ExpiresAt,
		User: &SessionCookieUser{
			ID:    session.User.ID,
			Email: session.User.Email,
		},
	}
}

type SessionCookieUser struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
}

type SessionCookieCompany struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SessionCookie struct {
	AccessToken  string                `json:"access_token"`
	RefreshToken string                `json:"refresh_token"`
	TokenType    string                `json:"token_type"`
	ExpiresIn    int                   `json:"expires_in"`
	ExpiresAt    int64                 `json:"expires_at"`
	Company      *SessionCookieCompany `json:"company"`
	User         *SessionCookieUser    `json:"user"`
}

func (s *SessionCookie) SetCookie(c echo.Context, secret string, env application.Environment) error {
	sameSiteMode := http.SameSiteLaxMode
	if env.IsProduction() {
		sameSiteMode = http.SameSiteStrictMode
	}

	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   env.IsProduction(),
		SameSite: sameSiteMode,
		MaxAge:   3600 * 24,
	}

	storeSession, err := store.Get(c.Request(), SessionCookieStoreName)
	if err != nil {
		return fmt.Errorf("error getting session cookie: %w", err)
	}

	storeSession.Values[SessionCookieSessionValueKey] = s
	return storeSession.Save(c.Request(), c.Response())
}

// SaveInContext saves the session on both the echo context and request context.
func (s *SessionCookie) SaveInContext(c echo.Context) {
	saveInContext(c, UserSessionContextKey, *s)
}

// Expired returns true if the session ExpiresAt is in the past.
func (s *SessionCookie) Expired() bool {
	expirationTime := time.Unix(s.ExpiresAt, 0)
	currentTime := time.Now()
	return currentTime.After(expirationTime)
}

func (s *SessionCookie) SetUser(u model.UserProfile) {
	s.User = &SessionCookieUser{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (s *SessionCookie) SetCompany(c model.Company) {
	s.Company = &SessionCookieCompany{
		ID:   c.ID,
		Name: c.Name,
	}
}

// GetSessionFromCookie returns the session form the cookie on the provided echo context.
// Returns an error if no cookie is present or data cannot be decoded.
func GetSessionFromCookie(c echo.Context, secret string) (*SessionCookie, error) {
	store := sessions.NewCookieStore([]byte(secret))
	storeSession, err := store.Get(c.Request(), SessionCookieSessionValueKey)
	if err != nil {
		return nil, ErrorCookieNotFound
	}

	data, ok := storeSession.Values[SessionCookieSessionValueKey]
	if !ok {
		return nil, ErrorCookieNotFound
	}

	sessionCookie, ok := data.(*SessionCookie)
	if !ok {
		return nil, errors.New("invalid session cookie data")
	}
	return sessionCookie, nil
}

func DeleteSessionCookie(c echo.Context, secret string) error {
	store := sessions.NewCookieStore([]byte(secret))
	storeSession, err := store.Get(c.Request(), SessionCookieSessionValueKey)
	if err != nil {
		return ErrorCookieNotFound
	}
	storeSession.Values[SessionCookieSessionValueKey] = nil
	storeSession.Options.MaxAge = -1
	return storeSession.Save(c.Request(), c.Response())
}
