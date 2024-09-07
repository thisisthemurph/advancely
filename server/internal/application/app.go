package application

import (
	"advancely/internal/store"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"

	_ "github.com/lib/pq"
)

type App struct {
	Config   AppConfig
	Store    *store.Store
	Supabase *supabase.Client
	Logger   *slog.Logger
}

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

func (app *App) Build() {
	app.Logger.Info("Initializing Supabase client")
	app.Supabase = supabase.CreateClient(app.Config.Supabase.URL, app.Config.Supabase.PublicKey)

	app.Logger.Info("Connecting stores to the database")
	s, err := store.NewStore(app.Config.Database.URI)
	if err != nil {
		app.Logger.Error("Failed to create database stores", "error", err)
		os.Exit(1)
	}
	app.Store = s
}

func (app *App) IsDevelopment() bool {
	return app.Config.IsDevelopment
}

func loadPossibleEnv(paths ...string) error {
	var err error
	for _, path := range paths {
		if err = godotenv.Load(path); err == nil {
			return nil
		}
	}
	return err
}
