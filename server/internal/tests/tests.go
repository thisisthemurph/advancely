package tests

import (
	"advancely/internal/auth"
	"advancely/internal/validation"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
)

// NewDefaultLogger creates a basic logger that logs to os.Stdout.
func NewDefaultLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// NewEchoInstance creates a basic echo instance with a validator configured.
func NewEchoInstance() *echo.Echo {
	e := echo.New()
	e.Validator = validation.NewCustomValidator()
	return e
}

// SaveSessionInContext saves the most basic data in the echo context for use with authorization.
func SaveSessionInContext(c echo.Context, userID, companyID uuid.UUID) {
	session := auth.SessionCookie{
		Company: &auth.SessionCookieCompany{
			ID: companyID,
		},
		User: &auth.SessionCookieUser{
			ID: userID,
		},
	}

	session.SaveInContext(c)
}

// NewRequestRecorder creates a test HTTP request and recorder, returning the associated echo context and response.
// The body should be a struct representing the request for the handler.
func NewRequestRecorder(t *testing.T, method, url string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	reqJSON, err := json.Marshal(body)
	require.NoError(t, err)

	e := NewEchoInstance()
	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
