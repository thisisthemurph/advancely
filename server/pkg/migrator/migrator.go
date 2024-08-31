package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
)

var ErrUnknownDirection = errors.New("unknown direction expected up or down")

type Migrator interface {
	Migrate(direction string) error
}

func NewPostgresMigrator(db *sql.DB, dbName, migrationPath string) Migrator {
	return &PostgresMigrator{
		DB:            db,
		DBName:        dbName,
		MigrationPath: migrationPath,
	}
}

type PostgresMigrator struct {
	DB            *sql.DB
	DBName        string
	MigrationPath string
}

func (m *PostgresMigrator) Migrate(direction string) error {
	driver, err := postgres.WithInstance(m.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	mig, err := migrate.NewWithDatabaseInstance(m.MigrationPath, m.DBName, driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	slog.Info("starting database migration", "direction", direction)

	switch direction {
	case "up":
		if err := mig.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("nothing to migrate")
				return nil
			}
			return fmt.Errorf("could not run migration: %w", err)
		}
	case "down":
		if err := mig.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("nothing to migrate")
				return nil
			}
			return fmt.Errorf("could not run migration: %w", err)
		}
	default:
		return fmt.Errorf("%w got %s", ErrUnknownDirection, direction)
	}

	slog.Info("database migration complete")
	return nil
}
