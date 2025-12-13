package config

import (
	"database/sql"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Router    *http.ServeMux
	DB        *sql.DB
	Queries   *sqlc.Queries
	Validator *validator.Validate
}

func NewConfig(db *sql.DB, queries *sqlc.Queries) *Config {
	router := http.NewServeMux()

	validate := validator.New(validator.WithRequiredStructEnabled())
	utils.RegisterCustomValidations(validate)

	app := &Config{
		Router:    router,
		DB:        db,
		Queries:   queries,
		Validator: validate,
	}

	return app
}
