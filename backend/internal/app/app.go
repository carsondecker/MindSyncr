package app

import (
	"github.com/carsondecker/MindSyncr/internal/config"
	"github.com/gorilla/mux"
)

type App struct {
	Config *config.Config
	Router *mux.Router
}

func NewApp() *App {
	cfg := config.NewConfig("8080")
	router := mux.NewRouter()

	app := &App{
		Config: cfg,
		Router: router,
	}

	app.registerRoutes()

	return app
}
