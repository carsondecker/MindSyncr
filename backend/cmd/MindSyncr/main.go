package main

import (
	"log"
	"net/http"
	"time"

	"github.com/carsondecker/MindSyncr/internal/app"
	"github.com/carsondecker/MindSyncr/internal/realtime"
)

func main() {
	app := app.NewApp()

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go testEvents(app)

	log.Fatal(srv.ListenAndServe())
}

func testEvents(app *app.App) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop() // cleanup

	for range ticker.C {
		log.Println("Running every 2 seconds")

		app.Hub.Broadcast <- realtime.Event{
			RoomID: "1",
			Msg:    "Hi! This message will repeat every 2 seconds.",
		}
	}
}
