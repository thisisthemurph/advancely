package main

import (
	"advancely/internal/application"
	"advancely/pkg/migrator"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func run(direction migrator.MigrationDirection, config application.AppConfig) error {
	db, err := sql.Open("postgres", config.Database.URI)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}
	defer db.Close()

	m := migrator.NewPostgresMigrator(db, config.Database.Name, migrator.DefaultMigrationPath)
	return m.Migrate(direction)
}

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintln(os.Stdout, "usage: migrate [direction: up|down]")
		os.Exit(1)
	}

	directionArg := os.Args[len(os.Args)-1]
	direction := migrator.NewMigrationDirection(directionArg)
	config := application.NewAppConfig(os.Getenv)

	if err := run(direction, config); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
