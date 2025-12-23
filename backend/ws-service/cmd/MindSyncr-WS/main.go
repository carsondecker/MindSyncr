package main

import (
	"log"
	"net/http"
	"time"

	"github.com/carsondecker/MindSyncr-WS/internal/app"
)

func main() {
	app := app.NewApp()

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
