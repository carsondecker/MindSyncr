package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/carsondecker/MindSyncr/internal/api"
	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Started MindSyncr Backend.")

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	utils.JWTInit()

	queries := sqlc.New(db)

	config := utils.NewConfig(db, queries)

	config.Router = api.GetRouter(config)

	srv := &http.Server{
		Handler:      config.Router,
		Addr:         ":3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
