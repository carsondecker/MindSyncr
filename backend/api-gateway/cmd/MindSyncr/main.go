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
	log.Println("Started MindSyncr API gateway.")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Build failed: Error loading .env file: %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) == 0 {
		log.Fatal("failed to get JWT_SECRET from .env")
	}
	utils.JWTInit(jwtSecret)

	redisAddr := os.Getenv("REDIS_URL")
	if len(redisAddr) == 0 {
		log.Fatal("failed to get REDIS_URL from .env")
	}
	redisClient, err := utils.NewRedisClient(redisAddr)
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(db)

	config := utils.NewConfig(db, queries, redisClient)

	config.Router = api.GetRouter(config)

	srv := &http.Server{
		Handler:      config.Router,
		Addr:         ":3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Build succeeded, starting MindSyncr server...")

	log.Fatal(srv.ListenAndServe())
}
