package main

import (
	"advancely/internal/application"
	"advancely/internal/routes"
	"advancely/pkg/migrator"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

const DatabaseMigrationPath = "file://cmd/migrate/migrations"

func main() {
	app := application.NewApp()

	if err := migrateDatabase(app.Config.Database); err != nil {
		slog.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	app.Build()

	router := routes.NewRouter(app)
	router.Logger.Fatal(router.Start(app.Config.Host))
}

func migrateDatabase(dbConfig application.DatabaseConfig) error {
	if !dbConfig.AutoMigrateOn {
		slog.Warn("Database migration disabled")
		return nil
	}

	db, err := sql.Open("postgres", dbConfig.URI)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}
	defer db.Close()

	m := migrator.NewPostgresMigrator(db, dbConfig.Name, DatabaseMigrationPath)
	return m.Migrate("up")
}
