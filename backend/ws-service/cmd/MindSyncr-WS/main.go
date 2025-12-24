package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/carsondecker/MindSyncr-WS/internal/app"
	"github.com/carsondecker/MindSyncr-WS/internal/utils"
)

func main() {
	redisAddr := os.Getenv("REDIS_URL")
	if len(redisAddr) == 0 {
		log.Fatal("failed to get REDIS_URL from .env")
	}
	redisClient, err := utils.NewRedisClient(redisAddr)
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(redisClient)

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         ":3001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
