package store

import (
	"advancely/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrCompanyNotFound = errors.New("company not found")

func NewPostgresCompanyStore(db *sqlx.DB) *PostgresCompanyStore {
	return &PostgresCompanyStore{
		DB: db,
	}
}

type PostgresCompanyStore struct {
	*sqlx.DB
}

func (s *PostgresCompanyStore) Company(id uuid.UUID) (model.Company, error) {
	var c model.Company
	if err := s.Get(&c, "SELECT * FROM companies WHERE id = $1;", id); err != nil {
		return model.Company{}, err
	}
	return c, nil
}

func (s *PostgresCompanyStore) CompanyByCreator(creatorID uuid.UUID) (model.Company, error) {
	var c model.Company
	if err := s.Get(&c, "SELECT * FROM companies WHERE creator_id = $1;", creatorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Company{}, ErrCompanyNotFound
		}
		return model.Company{}, err
	}
	return c, nil
}

func (s *PostgresCompanyStore) Companies() ([]model.Company, error) {
	var cc []model.Company
	if err := s.Get(&cc, "SELECT * FROM companies;"); err != nil {
		return []model.Company{}, err
	}
	return cc, nil
}

func (s *PostgresCompanyStore) CreateCompany(c *model.Company) error {
	query := `insert into companies (name, creator_id) values ($1, $2) returning *;`
	if err := s.Get(c, query, c.Name, c.CreatorID); err != nil {
		return err
	}
	return nil
}

func (s *PostgresCompanyStore) UpdateCompany(c *model.Company) error {
	if err := s.Get(c, "update companies set name = $1 where id = $2 returning *;", c.Name, c.ID); err != nil {
		return fmt.Errorf("error updating company with id %s: %w", c.ID, err)
	}
	return nil
}

func (s *PostgresCompanyStore) DeleteCompany(id uuid.UUID) error {
	if _, err := s.Exec("DELETE FROM companies WHERE id = $1;", id); err != nil {
		return fmt.Errorf("error deleting company with id %s: %w", id, err)
	}
	return nil
}
