package store

import (
	"advancely/pkg/errs"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrDomainAlreadyExists = errors.New("domain already exists")
)

func NewPostgresCompanySettingsStore(db *sqlx.DB) *PostgresCompanySettingsStore {
	return &PostgresCompanySettingsStore{
		DB: db,
	}
}

type PostgresCompanySettingsStore struct {
	*sqlx.DB
}

func (s *PostgresCompanySettingsStore) AddAllowedEmailDomain(
	ctx context.Context,
	companyID uuid.UUID,
	domain string,
) error {
	stmt := "insert into allowed_email_domains (company_id, domain) values ($1, $2);"
	if _, err := s.ExecContext(ctx, stmt, companyID, domain); err != nil {
		if pgErr := errs.CheckPgErr(err); errors.Is(pgErr, errs.PgErrCodeUniqueViolation) {
			return ErrDomainAlreadyExists
		}
		return err
	}
	return nil
}
