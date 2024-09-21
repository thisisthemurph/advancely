package tests

import (
	"advancely/internal/application"
	"advancely/pkg/migrator"
	"advancely/pkg/sbext"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/supabase-community/supabase-go"
	"testing"
)

const (
	DatabaseName      = "postgres"
	DatabaseURL       = "postgresql://postgres:postgres@127.0.0.1:54322/postgres?sslmode=disable"
	MigrationsPath    = "file://../../cmd/migrate/migrations"
	SupabaseURL       = "http://127.0.0.1:54321"
	SupabasePublicKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0"
)

func NewTestSupabaseClient(t *testing.T) *supabase.Client {
	sb, err := supabase.NewClient(SupabaseURL, SupabasePublicKey, nil)
	require.NoError(t, err)
	return sb
}

func NewTestSupabaseExtended(t *testing.T) *sbext.SupabaseExtended {
	config := application.SupabaseConfig{
		URL:       SupabaseURL,
		PublicKey: SupabasePublicKey,
	}
	return sbext.NewSupabaseExtended(NewTestSupabaseClient(t), config)
}

func SetUpTestDatabase(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("postgres", DatabaseURL)
	require.NoError(t, err)

	m := migrator.NewPostgresMigrator(db.DB, DatabaseName, MigrationsPath)
	err = m.Migrate(migrator.MigrationDirectionDown)
	require.NoError(t, err)

	err = m.Migrate(migrator.MigrationDirectionUp)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := m.Migrate(migrator.MigrationDirectionDown)
		require.NoError(t, err)
		_, err = db.Exec("delete from auth.users;")
		require.NoError(t, err)
		db.Close()
	})

	return db
}
