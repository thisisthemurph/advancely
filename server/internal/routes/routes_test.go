package routes_test

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"testing"

	"advancely/internal/tests"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/supabase-community/gotrue-go/types"
)

func setUpTestAdminUserAndCompany(t *testing.T) (*sqlx.DB, types.User, uuid.UUID) {
	db := tests.SetUpTestDatabase(t)
	sb := tests.NewTestSupabaseClient(t)
	user := tests.CreateAdminUser(t, sb, db)
	companyId := tests.CreateTestCompany(t, db, user.ID)
	return db, user, companyId
}

func assertHTTPError(t *testing.T, err error, expectedCode int, expectedMessage interface{}) {
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok, "expected echo.HTTPError")
	require.Equal(t, expectedCode, httpErr.Code)
	if expectedMessage != "" {
		require.Equal(t, expectedMessage, httpErr.Message)
	}
}
