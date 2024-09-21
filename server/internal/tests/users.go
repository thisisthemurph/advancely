package tests

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"testing"
)

func SignUpAdminUser(t *testing.T, sb *supabase.Client, db *sqlx.DB) types.User {
	resp, err := sb.Auth.Signup(types.SignupRequest{
		Email:    "user@company-email.com",
		Password: "password",
	})
	require.NoError(t, err)

	var adminRoleId int
	err = db.Get(&adminRoleId, "select id from security.roles where name = 'Admin' limit 1;")
	require.NoError(t, err)
	_, err = db.Exec(
		"insert into security.user_roles (user_id, role_id) values ($1, $2);",
		resp.User.ID, adminRoleId)
	require.NoError(t, err)

	return resp.User
}
