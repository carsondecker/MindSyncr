package app

import (
	"github.com/carsondecker/MindSyncr-WS/internal/realtime"
	"github.com/gorilla/mux"
)

type App struct {
	Config *Config
	Router *mux.Router
	Hub    *realtime.Hub
}

func NewApp() *App {
	cfg := NewConfig()
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
