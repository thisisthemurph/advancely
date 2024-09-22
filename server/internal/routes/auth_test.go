package routes_test

import (
	"advancely/internal/application"
	"advancely/internal/auth"
	"advancely/internal/routes"
	"advancely/internal/store"
	"advancely/internal/tests"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func newTestAuthHandler(t *testing.T, db *sqlx.DB) routes.AuthHandler {
	return routes.AuthHandler{
		Supabase:         tests.NewTestSupabaseExtended(t),
		UserStore:        store.NewPostgresUserStore(db),
		CompanyStore:     store.NewPostgresCompanyStore(db),
		PermissionsStore: store.NewPostgresPermissionsStore(db),
		Config: application.AppConfig{
			SessionSecret: "session.secret",
		},
		Logger: tests.NewDefaultLogger(),
	}
}

func TestAuthLogin(t *testing.T) {
	db := tests.SetUpTestDatabase(t)
	sb := tests.NewTestSupabaseClient(t)
	user := tests.SignUpAdminUser(t, sb, db)

	payload := map[string]string{
		"email":    user.Email,
		"password": tests.DefaultUserPassword,
	}
	c, rec := tests.NewRequestRecorder(t, http.MethodPost, "/login", payload)
	handler := newTestAuthHandler(t, db)
	err := handler.HandleLogin()(c)
	require.NoError(t, err)
	require.NotNil(t, rec)

	b, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	var bodyResponse auth.SessionCookie
	err = json.Unmarshal(b, &bodyResponse)
	require.NoError(t, err)
	require.Equal(t, user.Email, bodyResponse.User.Email)
}

func TestAuthLoginWithInvalidRequest(t *testing.T) {
	testCases := []struct {
		name    string
		payload interface{}
	}{
		{
			name:    "nothing passed",
			payload: map[string]string{},
		},
		{
			name: "no password provided",
			payload: map[string]string{
				"email": "user@email.com",
			},
		},
		{
			name: "no email provided",
			payload: map[string]string{
				"password": "this is a secret",
			},
		},
	}

	db := tests.SetUpTestDatabase(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := tests.NewRequestRecorder(t, http.MethodPost, "/login", tc.payload)

			handler := newTestAuthHandler(t, db)
			err := handler.HandleLogin()(c)
			if err != nil {
				c.Error(err)
			}

			assertHTTPError(t, err, http.StatusBadRequest, "")
		})
	}
}

func TestAuthLoginWhenUserDoesNotExist(t *testing.T) {
	db := tests.SetUpTestDatabase(t)
	payload := map[string]string{
		"email":    "unknown-user@domain.com",
		"password": "password",
	}
	c, _ := tests.NewRequestRecorder(t, http.MethodPost, "/login", payload)
	handler := newTestAuthHandler(t, db)
	err := handler.HandleLogin()(c)

	assertHTTPError(t, err, http.StatusBadRequest, "Invalid login credentials")
}

func TestAuthLoginWhenCompanyAndProfileDoesNotExist(t *testing.T) {
	db := tests.SetUpTestDatabase(t)
	sb := tests.NewTestSupabaseClient(t)
	user := tests.CreateAdminUser(t, sb, db)

	payload := map[string]string{
		"email":    user.Email,
		"password": tests.DefaultUserPassword,
	}
	c, _ := tests.NewRequestRecorder(t, http.MethodPost, "/login", payload)
	handler := newTestAuthHandler(t, db)
	err := handler.HandleLogin()(c)

	assertHTTPError(t, err, http.StatusBadRequest, store.ErrUserNotFound)
}

func TestTriggerPasswordReset(t *testing.T) {
	db := tests.SetUpTestDatabase(t)
	sb := tests.NewTestSupabaseClient(t)
	user := tests.SignUpAdminUser(t, sb, db)

	testCases := []struct {
		name                  string
		email                 string
		expectedStatus        int
		expectDatabaseUpdated bool
	}{
		{
			"success",
			user.Email,
			http.StatusNoContent,
			true,
		},
		{
			"user not existing",
			"fake@email.com",
			http.StatusNoContent,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]string{"email": tc.email}
			c, rec := tests.NewRequestRecorder(t, http.MethodPost, "/auth/password-reset", payload)

			handler := newTestAuthHandler(t, db)
			err := handler.HandleTriggerPasswordReset()(c)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectDatabaseUpdated {
				var recoveryToken string
				stmt := "select recovery_token from auth.users where id = $1;"
				err = db.Get(&recoveryToken, stmt, user.ID)
				require.NoError(t, err)
				require.NotEmpty(t, recoveryToken)
			}
		})
	}
}
