package routes_test

import (
	"testing"

	"advancely/internal/tests"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/supabase-community/gotrue-go/types"
)

func setUpTestAdminUserAndCompany(t *testing.T) (*sqlx.DB, types.User, uuid.UUID) {
	db := tests.SetUpTestDatabase(t)
	sb := tests.NewTestSupabaseClient(t)
	user := tests.SignUpAdminUser(t, sb, db)
	companyId := tests.CreateTestCompany(t, db, user.ID)
	return db, user, companyId
}
