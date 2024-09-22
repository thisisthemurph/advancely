package routes_test

import (
	"advancely/internal/model/security"
	"advancely/internal/routes"
	"advancely/internal/store"
	"advancely/internal/tests"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/supabase-community/supabase-go"
	"net/http"
	"testing"
)

func newTestUsersHandler(db *sqlx.DB, sb *supabase.Client, roleFetcher store.RoleFetcher) routes.UsersHandler {
	permissionsStore := store.NewPostgresPermissionsStore(db)
	var rf = roleFetcher
	if rf == nil {
		rf = permissionsStore
	}
	return routes.UsersHandler{
		UserStore:        store.NewPostgresUserStore(db),
		EnsurePermission: routes.EnsurePermissionsFnFactory(rf),
		Supabase:         sb,
		Logger:           tests.NewDefaultLogger(),
	}
}

func TestHandleCreateNewUser(t *testing.T) {
	db, user, companyId := setUpTestAdminUserAndCompany(t)
	sb := tests.NewTestSupabaseClient(t)

	payload := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "johndoe@advancelyexample.com",
	}

	c, rec := tests.NewRequestRecorder(t, http.MethodPost, "/login", payload)
	tests.SaveSessionInContext(c, user.ID, companyId)
	handler := newTestUsersHandler(db, sb, tests.NewFakeRoleFetcher(security.PermissionCreateUser))
	err := handler.HandleCreateNewUser()(c)
	if err != nil {
		c.Error(err)
	}

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Ensure there is a record in the auth.users and public.profiles tables
	var exists bool
	stmt := `
		select exists(
			select 1 
			from auth.users u 
			join public.profiles p on u.id = p.id
			where email = 'johndoe@advancelyexample.com');`
	err = db.QueryRow(stmt).Scan(&exists)
	require.NoError(t, err)
	require.True(t, exists)
}
