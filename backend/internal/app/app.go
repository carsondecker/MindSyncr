package app

import (
	"database/sql"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
)

type App struct {
	Router  *http.ServeMux
	DB      *sql.DB
	Queries *sqlc.Queries
}

func NewApp(db *sql.DB, queries *sqlc.Queries) *App {
	router := http.NewServeMux()

	app := &App{
		Router:  router,
		DB:      db,
		Queries: queries,
	}

	app.registerRoutes()

	return app
}
