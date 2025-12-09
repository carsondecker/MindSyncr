package app

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/auth"
	"github.com/carsondecker/MindSyncr/internal/config"
)

func GetRouter(cfg *config.Config) *http.ServeMux {
	baseRouter := http.NewServeMux()
	router := http.NewServeMux()

	baseRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	authRouter := http.NewServeMux()
	authHandler := auth.NewAuthHandler(cfg)
	authRouter.HandleFunc("POST /register", authHandler.HandleRegister)

	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	return baseRouter
}
