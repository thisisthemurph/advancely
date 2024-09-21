package routes_test

import (
	"net/http"
	"testing"

	"advancely/internal/application"
	"advancely/internal/routes"
	"advancely/internal/store"
	"advancely/internal/tests"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func newPermissionsHandler(db *sqlx.DB, roleFetcher store.RoleFetcher) routes.PermissionsHandler {
	permissionsStore := store.NewPostgresPermissionsStore(db)
	var rf = roleFetcher
	if rf == nil {
		rf = permissionsStore
	}
	return routes.PermissionsHandler{
		UserStore:        store.NewPostgresUserStore(db),
		PermissionsStore: permissionsStore,
		EnsurePermission: routes.EnsurePermissionsFnFactory(rf),
		Config:           application.AppConfig{},
		Logger:           tests.NewDefaultLogger(),
	}
}

func TestHandleCreateRole(t *testing.T) {
	db, user, companyId := setUpTestAdminUserAndCompany(t)
	payload := map[string]string{
		"name":        "new-role",
		"description": "new role description",
	}

	c, rec := tests.NewRequestRecorder(t, http.MethodPost, "/admin/companies/add-domain", payload)
	tests.SaveSessionInContext(c, user.ID, companyId)

	handler := newPermissionsHandler(db, nil)
	err := handler.HandleCreateRole()(c)
	if err != nil {
		c.Error(err)
	}

	require.NoError(t, err)
	require.Equal(t, 201, rec.Code)

	var roleName, roleDescription string
	err = db.QueryRow("select name, description from security.roles where company_id = $1 limit 1;", companyId).Scan(&roleName, &roleDescription)
	require.NoError(t, err)
	require.Equal(t, "new-role", roleName)
	require.Equal(t, "new role description", roleDescription)
}

func TestHandleCreateRoleWithIncorrectPermission(t *testing.T) {
	c, rec := tests.NewRequestRecorder(t, http.MethodPost, "/admin/companies/add-domain", map[string]string{})
	tests.SaveSessionInContext(c, uuid.New(), uuid.New())

	handler := newPermissionsHandler(nil, tests.NewFakeRoleFetcher())
	err := handler.HandleCreateRole()(c)
	if err != nil {
		c.Error(err)
	}

	require.Error(t, err)
	require.Equal(t, http.StatusForbidden, rec.Code)
}
