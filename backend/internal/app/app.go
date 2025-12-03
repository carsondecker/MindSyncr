package app

import (
	"github.com/carsondecker/MindSyncr/internal/config"
	"github.com/carsondecker/MindSyncr/internal/realtime"
	"github.com/gorilla/mux"
)

type App struct {
	Config *config.Config
	Router *mux.Router
	Hub    *realtime.Hub
}

func NewApp() *App {
	cfg := config.NewConfig("3000")
	router := mux.NewRouter()
	hub := realtime.NewHub()
	go hub.Run()

	app := &App{
		Config: cfg,
		Router: router,
		Hub:    hub,
	}

	app.registerRoutes()

	return app
}
