package application

import (
	"advancely/internal/store"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"

	_ "github.com/lib/pq"
)

type App struct {
	Config   AppConfig
	Store    *store.Store
	Supabase *supabase.Client
}

func NewApp() *App {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &App{
		Config: NewAppConfig(os.Getenv),
	}
}

func (app *App) Start() {
	app.Supabase = supabase.CreateClient(app.Config.Supabase.URL, app.Config.Supabase.PublicKey)

	s, err := store.NewStore(app.Config.Database.URI)
	if err != nil {
		log.Fatal(err)
	}
	app.Store = s
}
