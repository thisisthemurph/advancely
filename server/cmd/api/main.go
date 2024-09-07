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

func main() {
	app := application.NewApp()

	if err := migrateDatabase(app); err != nil {
		slog.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	app.Build()

	router := routes.NewRouter(app)
	router.Logger.Fatal(router.Start(app.Config.Host))
}

func migrateDatabase(app *application.App) error {
	dbConfig := app.Config.Database
	environment := app.Config.Environment
	shouldMigrateDatabase := environment.IsProduction() || dbConfig.AutoMigrateOn

	if !shouldMigrateDatabase {
		app.Logger.Warn("skipping database migration", "environment", environment, "AutoMigrateOn", dbConfig.AutoMigrateOn)
		return nil
	}

	if environment.IsProduction() {
		app.Logger.Info("migrating database as in production")
	} else if environment.IsDevelopment() && dbConfig.AutoMigrateOn {
		app.Logger.Info("migrating database in development as AutoMigrateOn is set to true")
	}

	db, err := sql.Open("postgres", dbConfig.URI)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}
	defer db.Close()

	m := migrator.NewPostgresMigrator(db, dbConfig.Name, migrator.DefaultMigrationPath).WithLogger(app.Logger)
	return m.Migrate(migrator.MigrationDirectionUp)
}
