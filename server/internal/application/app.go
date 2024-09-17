package application

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"advancely/internal/store"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

// App represents the configuration of the server application.
type App struct {
	Config   AppConfig
	Store    *store.PostgresStore
	Supabase *supabase.Client
	Logger   *slog.Logger
}

// NewApp creates a new basic App instance.
// Before this can be used Build needs to be called to set up the application.
func NewApp() *App {
	if err := loadPossibleEnv(".env", "../.env", "/etc/secrets/.env"); err != nil {
		log.Fatal(fmt.Errorf("error loading .env file: %w", err))
	}

	config := NewAppConfig(os.Getenv)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: nil,
		Level:       config.LogLevel,
	}))

	return &App{
		Config: config,
		Logger: logger,
	}
}

// Build connects the primary resources and persists them on the App struct.
func (app *App) Build() {
	app.Logger.Info("building app", "mode", app.Config.Environment.String())

	app.createSupabaseClient()

	if err := app.configureStores(); err != nil {
		app.Logger.Error("failed to configure stores", "error", err)
		os.Exit(1)
	}
}

// createSupabaseClient creates a new supabase client using the environment variables set on App.
func (app *App) createSupabaseClient() {
	app.Logger.Info("initializing Supabase client")

	client, err := supabase.NewClient(app.Config.Supabase.URL, app.Config.Supabase.PublicKey, nil)
	if err != nil {
		app.Logger.Error("failed to create Supabase client", "error", err)
		panic(err)
	}
	app.Supabase = client
}

// configureStores sets up the stores using the configured database URI.
func (app *App) configureStores() error {
	app.Logger.Info("connecting stores to the database")
	s, err := store.NewPostgresStore(app.Config.Database.URI)
	if err != nil {
		return err
	}
	app.Store = s
	return nil
}

// loadPossibleEnv takes any number of env paths and returns on the first success.
// Returns an error if none of the possible paths could be located.
func loadPossibleEnv(paths ...string) error {
	var err error
	for _, path := range paths {
		if err = godotenv.Load(path); err == nil {
			return nil
		}
	}
	return err
}
