package main

import (
	"advancely/internal/application"
	"advancely/internal/routes"
)

func main() {
	app := application.NewApp()
	app.Build()

	router := routes.NewRouter(app)
	router.Logger.Fatal(router.Start(app.Config.Host))
}
