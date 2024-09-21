package tests

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateTestCompany(t *testing.T, db *sqlx.DB, userId uuid.UUID) uuid.UUID {
	stmt := "insert into companies (name, creator_id) values ('test-company', $1) returning id;"
	var companyId uuid.UUID
	err := db.Get(&companyId, stmt, userId)
	require.NoError(t, err)
	return companyId
}
