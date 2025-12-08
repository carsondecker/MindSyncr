package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/carsondecker/MindSyncr/internal/app"
	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(db)

	app := app.NewApp(db, queries)

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
