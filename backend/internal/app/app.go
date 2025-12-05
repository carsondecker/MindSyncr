package app

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/config"
	"github.com/carsondecker/MindSyncr/internal/realtime"
)

type App struct {
	Config *config.Config
	Router *http.ServeMux
	Hub    *realtime.Hub
}

func NewApp() *App {
	cfg := config.NewConfig("3000")
	router := http.NewServeMux()
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
