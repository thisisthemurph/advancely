package store

import (
	"fmt"

	"advancely/internal/store/contract"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewStore(connectionString string) (*Store, error) {
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Store{
		UserStore:        NewUserStore(db),
		CompanyStore:     NewCompanyStore(db),
		PermissionsStore: NewPermissionsStore(db),
	}, nil
}

type Store struct {
	contract.UserStore
	contract.CompanyStore
	contract.PermissionsStore
}
