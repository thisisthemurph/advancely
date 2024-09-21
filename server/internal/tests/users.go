package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

const DefaultUserPassword = "password"

// CreateAdminUser creates a user and associates them with the Admin role.
// This method does not to a complete signup. The profiles and companies table is not affected.
func CreateAdminUser(t *testing.T, sb *supabase.Client, db *sqlx.DB) types.User {
	resp, err := sb.Auth.Signup(types.SignupRequest{
		Email:    "user@company-email.com",
		Password: DefaultUserPassword,
	})
	require.NoError(t, err)

	// Add the admin role
	var adminRoleId int
	err = db.Get(&adminRoleId, "select id from security.roles where name = 'Admin' limit 1;")
	require.NoError(t, err)
	_, err = db.Exec(
		"insert into security.user_roles (user_id, role_id) values ($1, $2);",
		resp.User.ID, adminRoleId)
	require.NoError(t, err)

	return resp.User
}

// SignUpAdminUser simulates a completely signed up Admin user.
// The user is created with the Admin role and associate profile and company records.
func SignUpAdminUser(t *testing.T, sb *supabase.Client, db *sqlx.DB) types.User {
	user := CreateAdminUser(t, sb, db)

	// Create company
	var companyID uuid.UUID
	stmt := "insert into companies (name, creator_id) values ($1, $2) returning id;"
	err := db.Get(&companyID, stmt, "test-company", user.ID)
	require.NoError(t, err)

	// Create profile
	stmt = "insert into profiles (id, company_id, first_name, last_name) values ($1, $2, $3, $4);"
	_, err = db.Exec(stmt, user.ID, companyID, "Joe", "Blogs")
	require.NoError(t, err)

	return user
}
