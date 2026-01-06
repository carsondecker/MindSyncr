package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

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

	wsSecret := os.Getenv("WS_SECRET")
	if len(wsSecret) == 0 {
		log.Fatal("failed to get WS_SECRET from .env")
	}

	utils.JWTInit(jwtSecret, wsSecret)

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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      c.Handler(config.Router),
		Addr:         ":3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Build succeeded, starting MindSyncr server...")

	log.Fatal(srv.ListenAndServe())
}
