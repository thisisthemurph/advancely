package auth

import (
	"advancely/internal/model"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
	"net/http"
	"time"
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

type SessionCookieCompany struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SessionCookieUser struct {
	ID        uuid.UUID             `json:"id"`
	FirstName string                `json:"firstName"`
	LastName  string                `json:"lastName"`
	Email     string                `json:"email"`
	Company   *SessionCookieCompany `json:"company,omitempty"`
}

type SessionCookie struct {
	Sub          uuid.UUID          `json:"sub"`
	Aud          string             `json:"aud"`
	Role         string             `json:"role"`
	Email        string             `json:"email"`
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	ExpiresAt    time.Time          `json:"expiresAt"`
	User         *SessionCookieUser `json:"user,omitempty"`
}

func NewSessionCookie(details *supabase.AuthenticatedDetails) *SessionCookie {
	userID, _ := uuid.Parse(details.User.ID)
	return &SessionCookie{
		Sub:          userID,
		Aud:          details.User.Aud,
		Role:         details.User.Role,
		Email:        details.User.Email,
		AccessToken:  details.AccessToken,
		RefreshToken: details.RefreshToken,
		ExpiresAt:    time.Now().UTC().Add(time.Duration(details.ExpiresIn) * time.Second),
	}
}

// SetCookie saves the SessionCookie using the provided echo context.
// Returns an error if there are any issues encoding the session data.
func (s *SessionCookie) SetCookie(c echo.Context, secret string, isDevelopment bool) error {
	sameSite := http.SameSiteLaxMode
	if isDevelopment {
		sameSite = http.SameSiteNoneMode
	}

	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   !isDevelopment,
		SameSite: sameSite,
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
	return time.Now().After(s.ExpiresAt)
}

func (s *SessionCookie) SetUser(u model.UserProfile) {
	s.User = &SessionCookieUser{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

func (s *SessionCookie) SetUserCompany(c model.Company) {
	if s.User == nil {
		s.User = &SessionCookieUser{}
	}
	s.User.Company = &SessionCookieCompany{
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
