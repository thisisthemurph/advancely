package application

import (
	"advancely/internal/store"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
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
