package main

import (
	"advancely/internal/application"
	"advancely/internal/routes"
)

func main() {
	app := application.NewApp()
	app.Start()

	e := routes.NewRouter(app)
	e.Logger.Fatal(e.Start(app.Config.Host))
}
