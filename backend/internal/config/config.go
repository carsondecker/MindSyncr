package config

import (
	"database/sql"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
)

type Config struct {
	Router  *http.ServeMux
	DB      *sql.DB
	Queries *sqlc.Queries
}

func NewConfig(db *sql.DB, queries *sqlc.Queries) *Config {
	router := http.NewServeMux()

	app := &Config{
		Router:  router,
		DB:      db,
		Queries: queries,
	}

	return app
}
