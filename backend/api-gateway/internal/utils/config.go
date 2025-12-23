package utils

import (
	"database/sql"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Router      *http.ServeMux
	DB          *sql.DB
	Queries     *sqlc.Queries
	Validator   *validator.Validate
	RedisClient *RedisClient
}

func NewConfig(db *sql.DB, queries *sqlc.Queries, redisClient *RedisClient) *Config {
	router := http.NewServeMux()

	validate := validator.New(validator.WithRequiredStructEnabled())
	RegisterCustomValidations(validate)

	app := &Config{
		Router:      router,
		DB:          db,
		Queries:     queries,
		Validator:   validate,
		RedisClient: redisClient,
	}

	return app
}
